package queue

import (
	"sort"
	"time"

	"github.com/Jleagle/rabbit-go"
	"github.com/gamedb/gamedb/pkg/helpers"
	"github.com/gamedb/gamedb/pkg/helpers/memcache"
	steamHelper "github.com/gamedb/gamedb/pkg/helpers/steam"
	"github.com/gamedb/gamedb/pkg/log"
	"github.com/gamedb/gamedb/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type AppAchievementsMessage struct {
	AppID   int    `json:"id"`
	AppName string `json:"app-name"`
}

func appAchievementsHandler(messages []*rabbit.Message) {

	for _, message := range messages {

		payload := AppAchievementsMessage{}

		err := helpers.Unmarshal(message.Message.Body, &payload)
		if err != nil {
			log.Err(err, message.Message.Body)
			sendToFailQueue(message)
			continue
		}

		// Get owners
		owners, err := mongo.CountDocuments(mongo.CollectionPlayerApps, bson.D{{"app_id", payload.AppID}}, 0)
		if err != nil {
			log.Err(err)
			sendToRetryQueue(message)
			continue
		}

		//
		schemaResponse, _, err := steamHelper.GetSteam().GetSchemaForGame(payload.AppID)
		err = steamHelper.AllowSteamCodes(err, 400, 403)
		if err != nil {
			steamHelper.LogSteamError(err)
			sendToRetryQueue(message)
			continue
		}

		globalResponse, _, err := steamHelper.GetSteam().GetGlobalAchievementPercentagesForApp(payload.AppID)
		err = steamHelper.AllowSteamCodes(err, 403, 500)
		if err != nil {
			steamHelper.LogSteamError(err)
			sendToRetryQueue(message)
			continue
		}

		// Build map of all global achievements
		var achievementsMap = map[string]mongo.AppAchievement{}

		for _, achievement := range globalResponse.GlobalAchievementPercentage {

			achievementsMap[achievement.Name] = mongo.AppAchievement{
				AppID:     payload.AppID,
				Key:       achievement.Name,
				Completed: achievement.Percent,
				Deleted:   false,
			}
		}

		// Add in data for achievements in schema
		var percentTotal float64
		var percentCount int
		var wait bool

		for _, achievement := range schemaResponse.AvailableGameStats.Achievements {

			if val, ok := achievementsMap[achievement.Name]; ok {

				percentTotal += val.Completed
				percentCount++

				val.Name = achievement.DisplayName
				val.SetIcon(achievement.Icon)
				val.Description = achievement.Description
				val.Hidden = bool(achievement.Hidden)
				val.Active = true
				val.AppOwners = owners

				achievementsMap[achievement.Name] = val

			} else {
				log.Info("Achevement in schema but not global", payload.AppID, achievement.Name)
				wait = true
				break
			}
		}

		// Wait for both API endpoints to sync up
		if wait {
			sendToRetryQueueWithDelay(message, 1*time.Minute)
			continue
		}

		// Save achievements to Mongo
		var achievementsSlice []mongo.AppAchievement
		for _, achievement := range achievementsMap {
			achievementsSlice = append(achievementsSlice, achievement)
		}

		// Sort by key to store the first 5
		sort.Slice(achievementsSlice, func(i, j int) bool {
			return achievementsSlice[i].Completed > achievementsSlice[j].Completed
		})

		err = mongo.SaveAppAchievements(achievementsSlice)
		if err != nil {
			log.Err(err, payload.AppID)
			sendToRetryQueue(message)
			continue
		}

		// Update in Elastic
		for _, v := range achievementsSlice {
			err = ProduceAchievementSearch(v, payload.AppName)
			log.Err(err)
		}

		// Update app row
		var achievementsCol []mongo.AppAchievement
		for _, achievement := range achievementsSlice {
			if achievement.Active && achievement.Icon != "" {
				if len(achievementsCol) < 5 {
					achievementsCol = append(achievementsCol, achievement)
				} else {
					break
				}
			}
		}

		var average float64
		if percentCount != 0 {
			average = percentTotal / float64(percentCount)
		}

		var stats []helpers.AppStat
		for _, v := range schemaResponse.AvailableGameStats.Stats {
			stats = append(stats, helpers.AppStat{
				Name:        v.Name,
				Default:     v.DefaultValue,
				DisplayName: v.DisplayName,
			})
		}

		var updateApp = []bson.E{
			{"version", schemaResponse.Version},
			{"achievements_count", len(schemaResponse.AvailableGameStats.Achievements)},
			{"achievements_count_total", len(globalResponse.GlobalAchievementPercentage)},
			{"achievements_5", achievementsCol},
			{"achievements_average_completion", average},
			{"stats", stats},
		}

		_, err = mongo.UpdateOne(mongo.CollectionApps, bson.D{{"_id", payload.AppID}}, updateApp)
		if err != nil {
			log.Err(err, payload.AppID)
			sendToRetryQueue(message)
			continue
		}

		// Mark apps in Mongo but not in global response as deleted
		var filter = bson.D{{"app_id", payload.AppID}}
		if len(achievementsMap) > 0 {
			var keys []string
			for k := range achievementsMap {
				keys = append(keys, k)
			}
			filter = append(filter, bson.E{Key: "key", Value: bson.M{"$nin": keys}})
		}

		_, err = mongo.UpdateManySet(mongo.CollectionAppAchievements, filter, bson.D{{"deleted", true}})
		if err != nil {
			log.Err(err, payload.AppID)
			sendToRetryQueue(message)
			continue
		}

		//
		var items = []string{
			memcache.MemcacheApp(payload.AppID).Key,
			memcache.MemcacheMongoCount(mongo.CollectionAppAchievements.String(), bson.D{{"app_id", payload.AppID}}).Key,
		}

		err = memcache.Delete(items...)
		if err != nil {
			log.Err(err, payload.AppID)
			sendToRetryQueue(message)
			continue
		}

		message.Ack(false)
	}
}
