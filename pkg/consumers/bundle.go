package consumers

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Jleagle/steam-go/steam"
	"github.com/cenkalti/backoff"
	"github.com/gamedb/gamedb/pkg/consumers/framework"
	"github.com/gamedb/gamedb/pkg/helpers"
	steamHelper "github.com/gamedb/gamedb/pkg/helpers/steam"
	"github.com/gamedb/gamedb/pkg/log"
	"github.com/gamedb/gamedb/pkg/mongo"
	"github.com/gamedb/gamedb/pkg/sql"
	"github.com/gamedb/gamedb/pkg/websockets"
	"github.com/gocolly/colly"
	"go.mongodb.org/mongo-driver/bson"
)

type BundleMessage struct {
	ID int `json:"id"`
}

func (msg BundleMessage) Produce() error {
	return channels[framework.Producer][queueBundles].ProduceInterface(msg)
}

func bundleHandler(messages []*framework.Message) {

	for _, message := range messages {

		payload := BundleMessage{}

		err := helpers.Unmarshal(message.Message.Body, &payload)
		if err != nil {
			log.Err(err, message.Message.Body)
			sendToFailQueue(message)
			return
		}

		// Load current bundle
		gorm, err := sql.GetMySQLClient()
		if err != nil {
			log.Err(err, message.Message.Body)
			sendToRetryQueue(message)
			return
		}

		bundle := sql.Bundle{}
		gorm = gorm.FirstOrInit(&bundle, sql.Bundle{ID: payload.ID})
		if gorm.Error != nil {
			log.Err(gorm.Error, message.Message.Body)
			sendToRetryQueue(message)
			return
		}

		oldBundle := bundle

		err = updateBundle(&bundle)
		if err != nil && err != steam.ErrAppNotFound {
			steamHelper.LogSteamError(err, message.Message.Body)
			sendToRetryQueue(message)
			return
		}

		var wg sync.WaitGroup

		// Save new data
		wg.Add(1)
		go func() {

			defer wg.Done()

			gorm = gorm.Save(&bundle)
			if gorm.Error != nil {
				log.Err(gorm.Error, message.Message.Body)
				sendToRetryQueue(message)
				return
			}
		}()

		// Save to InfluxDB
		wg.Add(1)
		go func() {

			defer wg.Done()

			var err error

			err = saveBundlePriceToMongo(bundle, oldBundle)
			if err != nil {
				log.Err(err, message.Message.Body)
				sendToRetryQueue(message)
				return
			}
		}()

		wg.Wait()

		if message.ActionTaken {
			return
		}

		// Send websocket
		wsPaload := websockets.PubSubIDPayload{}
		wsPaload.ID = payload.ID
		wsPaload.Pages = []websockets.WebsocketPage{websockets.PageBundle, websockets.PageBundles}

		_, err = helpers.Publish(helpers.PubSubTopicWebsockets, wsPaload)
		if err != nil {
			log.Err(err, payload.ID)
		}

		message.Ack()
	}
}

func updateBundle(bundle *sql.Bundle) (err error) {

	c := colly.NewCollector(
		colly.AllowedDomains("store.steampowered.com"),
		colly.AllowURLRevisit(), // This is for retrys
	)

	steamHelper.SetAgeCheckCookieJar(c)

	// Title
	c.OnHTML("h2.pageheader", func(e *colly.HTMLElement) {
		bundle.Name = e.Text
	})

	// Image
	c.OnHTML("img.package_header", func(e *colly.HTMLElement) {
		bundle.Image = e.Attr("src")
	})

	// Discount
	c.OnHTML(".game_purchase_discount .bundle_base_discount", func(e *colly.HTMLElement) {
		var discount int
		discount, err = strconv.Atoi(strings.Replace(e.Text, "%", "", 1))
		bundle.SetDiscount(discount)
	})

	// Bigger discount
	c.OnHTML(".game_purchase_discount .discount_pct", func(e *colly.HTMLElement) {
		var discount int
		discount, err = strconv.Atoi(strings.Replace(e.Text, "%", "", 1))
		bundle.SetDiscount(discount)
	})

	// Apps
	var apps []string
	c.OnHTML("[data-ds-appid]", func(e *colly.HTMLElement) {
		apps = append(apps, strings.Split(e.Attr("data-ds-appid"), ",")...)
	})

	// Packages
	var packages []string
	c.OnHTML("[data-ds-packageid]", func(e *colly.HTMLElement) {
		packages = append(packages, strings.Split(e.Attr("data-ds-packageid"), ",")...)
	})

	// Retry call
	operation := func() (err error) {
		return c.Visit("https://store.steampowered.com/bundle/" + strconv.Itoa(bundle.ID))
	}

	policy := backoff.NewExponentialBackOff()

	err = backoff.RetryNotify(operation, policy, func(err error, t time.Duration) { log.Info(err) })
	if err != nil {
		return err
	}

	//
	if len(apps) == 0 && len(packages) == 0 {
		return nil
	}

	// Apps
	b, err := json.Marshal(helpers.StringSliceToIntSlice(apps))
	if err != nil {
		return err
	}

	bundle.AppIDs = string(b)

	// Packages
	b, err = json.Marshal(helpers.StringSliceToIntSlice(packages))
	if err != nil {
		return err
	}

	bundle.PackageIDs = string(b)

	return nil
}

var bundlePriceLock sync.Mutex

func saveBundlePriceToMongo(bundle sql.Bundle, oldBundle sql.Bundle) (err error) {

	bundlePriceLock.Lock()
	defer bundlePriceLock.Unlock()

	time.Sleep(time.Second) // prices are keyed by the second

	if bundle.Discount != oldBundle.Discount {

		doc := mongo.BundlePrice{
			CreatedAt: time.Now(),
			BundleID:  bundle.ID,
			Discount:  bundle.Discount,
		}

		// Does a replace, as sometimes doing a InsertOne would error on key already existing
		_, err = mongo.ReplaceOne(mongo.CollectionBundlePrices, bson.D{{"_id", doc.GetKey()}}, doc)
	}

	return err
}
