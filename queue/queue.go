package queue

import (
	"encoding/json"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/Jleagle/go-durationfmt"
	"github.com/cenkalti/backoff"
	"github.com/gamedb/website/config"
	"github.com/gamedb/website/log"
	"github.com/streadway/amqp"
)

type queueName string

//noinspection GoUnusedConst
const (
	// C#
	queueCSApps     queueName = "GameDB_CS_Apps"
	queueCSPackages queueName = "GameDB_CS_Packages"
	queueCSProfiles queueName = "GameDB_CS_Profiles"

	// Go
	queueGoApps     queueName = "GameDB_Go_Apps"
	queueGoBundles  queueName = "GameDB_Go_Bundles"
	queueGoChanges  queueName = "GameDB_Go_Changes"
	queueGoDelays   queueName = "GameDB_Go_Delays"
	queueGoPackages queueName = "GameDB_Go_Packages"
	queueGoProfiles queueName = "GameDB_Go_Profiles"
	queueGoFailed   queueName = "GameDB_Go_Failed"

	//
	maxBytesToStore int = 1024 * 10
)

var (
	consumeLock sync.Mutex
	produceLock sync.Mutex

	consumerConnection *amqp.Connection
	producerConnection *amqp.Connection

	consumerCloseChannel = make(chan *amqp.Error)
	producerCloseChannel = make(chan *amqp.Error)

	queues = []baseQueue{
		{
			name:    queueGoApps,
			queue:   &appQueue{},
			maxTime: time.Hour * 24 * 7,
		},
		{
			name:    queueGoBundles,
			queue:   &bundleQueue{},
			maxTime: time.Hour * 24 * 7,
		},
		{
			name:    queueGoChanges,
			queue:   &changeQueue{},
			maxTime: time.Hour * 24 * 7,
		},
		{
			name:    queueGoDelays,
			queue:   &delayQueue{},
			maxTime: time.Hour * 24 * 7,
		},
		{
			name:    queueGoPackages,
			queue:   &packageQueue{},
			maxTime: time.Hour * 24 * 7,
		},
		{
			name:    queueGoProfiles,
			queue:   &playerQueue{},
			maxTime: time.Hour * 24 * 7,
		},
	}
)

type baseMessage struct {
	Message       interface{} `json:"message"`
	FirstSeen     time.Time   `json:"first_seen"`
	Attempt       int         `json:"attempt"`
	OriginalQueue queueName   `json:"original_queue"`
}

func (payload baseMessage) getNextAttempt() time.Time {

	var min float64 = 1
	var max float64 = 600

	var seconds float64
	seconds = math.Pow(1.5, float64(payload.Attempt))
	seconds = math.Min(seconds, max)
	seconds = math.Max(seconds, min)

	return payload.FirstSeen.Add(time.Second * time.Duration(int64(seconds)))
}

// Remove from queue
func (payload baseMessage) ack(msg amqp.Delivery) {

	err := msg.Ack(false)
	logError(err)
}

// Send to failed queue
func (payload baseMessage) fail(msg amqp.Delivery) {

	err := produce(payload, queueGoFailed)
	if err != nil {
		logError(err)
		return
	}

	err = msg.Ack(false)
	if err != nil {
		logError(err)
		return
	}
}

// Send to delay queue
func (payload baseMessage) ackRetry(msg amqp.Delivery) {

	payload.Attempt++

	totalStr, err := durationfmt.Format(payload.getNextAttempt().Sub(payload.FirstSeen), "%mm %ss")
	if err != nil {
		logError(err)
	}

	leftStr, err := durationfmt.Format(payload.getNextAttempt().Sub(time.Now()), "%mm %ss")
	if err != nil {
		logError(err)
	}

	logInfo("Adding to delay queue for " + leftStr + ", " + totalStr + " total, attempt " + strconv.Itoa(payload.Attempt))

	err = produce(payload, queueGoDelays)
	if err != nil {
		logError(err)
		return
	}

	err = msg.Ack(false)
	if err != nil {
		logError(err)
		return
	}
}

type queueInterface interface {
	processMessages(msgs []amqp.Delivery)
}

type baseQueue struct {
	queue       queueInterface
	name        queueName
	batchSize   int
	maxAttempts int
	maxTime     time.Duration
}

func (q baseQueue) consumeMessages() {

	var err error

	for {

		// Connect
		err = func() error {

			consumeLock.Lock()
			defer consumeLock.Unlock()

			if consumerConnection == nil {

				consumerConnection, err = makeAConnection()
				if err != nil {
					logCritical("Connecting to Rabbit: " + err.Error())
					return err
				}
				consumerConnection.NotifyClose(consumerCloseChannel)
			}

			return nil
		}()

		if err != nil {
			logError(err)
			return
		}

		//
		ch, qu, err := getQueue(consumerConnection, q.name)
		if err != nil {
			logError(err)
			return
		}

		msgs, err := ch.Consume(qu.Name, "", false, false, false, false, nil)
		if err != nil {
			logError(err)
			return
		}

		// In a anon function so can return at anytime
		func(msgs <-chan amqp.Delivery, q baseQueue) {

			var msgSlice []amqp.Delivery

			for {
				select {
				case err = <-consumerCloseChannel:
					logWarning(err)
					return
				case msg := <-msgs:
					msgSlice = append(msgSlice, msg)
				}

				if len(msgSlice) >= q.batchSize {
					q.queue.processMessages(msgSlice)
					msgSlice = []amqp.Delivery{}
				}
			}

		}(msgs, q)

		logWarning("Rabbit consumer connection has disconnected")

		err = ch.Close()
		logError(err)
	}
}

func RunConsumers() {
	for _, queue := range queues {
		go queue.consumeMessages()
	}
}

func produce(payload baseMessage, queue queueName) (err error) {

	if payload.OriginalQueue == "" {
		payload.OriginalQueue = queue
	}
	if payload.FirstSeen.IsZero() {
		payload.FirstSeen = time.Now()
	}
	if payload.Attempt == 0 {
		payload.Attempt = 1
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Connect
	err = func() error {

		produceLock.Lock()
		defer produceLock.Unlock()

		if producerConnection == nil {

			producerConnection, err = makeAConnection()
			if err != nil {
				logCritical("Connecting to Rabbit: " + err.Error())
				return err
			}
			producerConnection.NotifyClose(producerCloseChannel)
		}

		return nil
	}()

	if err != nil {
		return err
	}

	//
	ch, qu, err := getQueue(producerConnection, queue)
	if err != nil {
		return err
	}

	err = ch.Publish("", qu.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         b,
	})
	if err != nil {
		return err
	}

	return ch.Close()
}

func makeAConnection() (conn *amqp.Connection, err error) {

	operation := func() (err error) {

		logInfo("Connecting to Rabbit")

		conn, err = amqp.Dial(config.Config.RabbitDSN())
		logError(err) // Logging here as no max elasped time
		return err
	}

	policy := backoff.NewExponentialBackOff()
	policy.MaxElapsedTime = 0
	policy.InitialInterval = 5 * time.Second

	err = backoff.RetryNotify(operation, policy, func(err error, t time.Duration) { logInfo(err) })

	return conn, err
}

func getQueue(conn *amqp.Connection, queue queueName) (ch *amqp.Channel, qu amqp.Queue, err error) {

	ch, err = conn.Channel()
	if err != nil {
		return
	}

	err = ch.Qos(10, 0, false)
	if err != nil {
		return
	}

	qu, err = ch.QueueDeclare(string(queue), true, false, false, false, nil)

	return ch, qu, err
}

//
type steamKitJob struct {
	SequentialCount int    `json:"SequentialCount"`
	StartTime       string `json:"StartTime"`
	ProcessID       int    `json:"ProcessID"`
	BoxID           int    `json:"BoxID"`
	Value           int64  `json:"Value"`
}

func logInfo(interfaces ...interface{}) {
	log.Info(append(interfaces, log.LogNameConsumers)...)
}

func logError(interfaces ...interface{}) {
	log.Err(append(interfaces, log.LogNameConsumers)...)
}

func logWarning(interfaces ...interface{}) {
	log.Warning(append(interfaces, log.LogNameConsumers)...)
}

func logCritical(interfaces ...interface{}) {
	log.Critical(append(interfaces, log.LogNameConsumers)...)
}

func ProduceBundle(ID int, appID int) (err error) {

	return produce(baseMessage{
		Message: bundleMessage{
			ID:    ID,
			AppID: appID,
		},
	}, queueGoBundles)
}

func ProduceApp(ID int) (err error) {

	return produce(baseMessage{
		Message: appMessage{
			ID: ID,
		},
	}, queueCSApps)
}

func ProducePackage(ID int) (err error) {

	return produce(baseMessage{
		Message: packageMessage{
			ID: ID,
		},
	}, queueCSPackages)
}

func ProducePlayer(ID int64) (err error) {

	return produce(baseMessage{
		Message: playerMessage{
			ID: ID,
		},
	}, queueCSProfiles)
}
