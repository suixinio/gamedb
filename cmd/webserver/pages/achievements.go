package pages

import (
	"net/http"
	"sync"

	"github.com/gamedb/gamedb/cmd/webserver/helpers/datatable"
	"github.com/gamedb/gamedb/cmd/webserver/helpers/session"
	"github.com/gamedb/gamedb/pkg/log"
	"github.com/gamedb/gamedb/pkg/mongo"
	"github.com/go-chi/chi"
	"go.mongodb.org/mongo-driver/bson"
)

func AchievementsRouter() http.Handler {

	r := chi.NewRouter()
	r.Get("/", achievementsHandler)
	r.Get("/achievements.json", achievementsAjaxHandler)
	return r
}

func achievementsHandler(w http.ResponseWriter, r *http.Request) {

	t := GlobalTemplate{}
	t.fill(w, r, "Achievements", "Games with the most achievements")
	t.addAssetJSON2HTML()

	returnTemplate(w, r, "achievements", t)
}

func achievementsAjaxHandler(w http.ResponseWriter, r *http.Request) {

	query := datatable.NewDataTableQuery(r, false)

	var wg sync.WaitGroup
	var filter = bson.D{{"achievements_count", bson.M{"$gt": 0}}}
	var filter2 = filter
	var countLock sync.Mutex

	var search = query.GetSearchString("search")
	if search != "" {
		filter2 = append(filter2, bson.E{Key: "$text", Value: bson.M{"$search": search}})
	}

	var apps []mongo.App
	wg.Add(1)
	go func() {

		defer wg.Done()

		var columns = map[string]string{
			"1": "achievements_count",
			"2": "achievements_average_completion",
		}

		var projection = bson.M{"_id": 1, "name": 1, "icon": 1, "achievements_5": 1, "achievements_count": 1, "achievements_average_completion": 1, "prices": 1}
		var sort = query.GetOrderMongo(columns)
		sort = append(sort, bson.E{Key: "achievements_average_completion", Value: -1})

		var err error
		apps, err = mongo.GetApps(query.GetOffset64(), 100, sort, filter2, projection, nil)
		if err != nil {
			log.Err(err, r)
		}
	}()

	var filtered int64
	wg.Add(1)
	go func() {

		defer wg.Done()

		var err error
		countLock.Lock()
		filtered, err = mongo.CountDocuments(mongo.CollectionApps, filter2, 0)
		countLock.Unlock()
		if err != nil {
			log.Err(err, r)
		}
	}()

	var count int64
	wg.Add(1)
	go func() {

		defer wg.Done()

		var err error
		countLock.Lock()
		count, err = mongo.CountDocuments(mongo.CollectionApps, filter, 60*60*24)
		countLock.Unlock()
		if err != nil {
			log.Err(err, r)
		}
	}()

	wg.Wait()

	//
	var code = session.GetProductCC(r)
	var response = datatable.NewDataTablesResponse(r, query, count, filtered)
	for _, app := range apps {

		for k, v := range app.Achievements {
			app.Achievements[k].Icon = v.GetIcon()
		}

		response.AddRow([]interface{}{
			app.ID,                            // 0
			app.GetName(),                     // 1
			app.GetIcon(),                     // 2
			app.GetPath() + "#achievements",   // 3
			app.Prices.Get(code).GetFinal(),   // 4
			app.AchievementsCount,             // 5
			app.AchievementsAverageCompletion, // 6
			app.Achievements,                  // 7
		})
	}

	returnJSON(w, r, response)
}
