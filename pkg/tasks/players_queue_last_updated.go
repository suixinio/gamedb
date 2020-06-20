package tasks

import (
	"time"

	"github.com/Jleagle/rabbit-go"
	"github.com/gamedb/gamedb/pkg/helpers"
	"github.com/gamedb/gamedb/pkg/helpers/memcache"
	"github.com/gamedb/gamedb/pkg/mongo"
	"github.com/gamedb/gamedb/pkg/queue"
	"go.mongodb.org/mongo-driver/bson"
)

type PlayersQueueLastUpdated struct {
	BaseTask
}

func (c PlayersQueueLastUpdated) ID() string {
	return "update-last-updated-players"
}

func (c PlayersQueueLastUpdated) Name() string {
	return "Update last updated players"
}

func (c PlayersQueueLastUpdated) Cron() string {
	return CronTimeUpdateLastUpdatedPlayers
}

const toQueue = 10
const cronTime = time.Minute

func (c PlayersQueueLastUpdated) work() (err error) {

	// Skip if queues have activity
	limits := map[rabbit.QueueName]int{
		queue.QueueApps:     50,
		queue.QueuePackages: 50,
		queue.QueuePlayers:  5,
	}

	queues, err := helpers.RabbitClient.GetQueues()
	if err != nil {
		return err
	}

	var consumers int
	for _, q := range queues {
		if val, ok := limits[rabbit.QueueName(q.Name)]; ok && q.Messages > val {
			return nil
		}
		if q.Name == string(queue.QueuePlayers) {
			consumers = q.Consumers
		}
	}

	if consumers == 0 {
		return nil
	}

	// Queue last updated players
	players, err := mongo.GetPlayers(0, int64(toQueue*consumers), bson.D{{"updated_at", 1}}, nil, bson.M{"_id": 1})
	if err != nil {
		return err
	}

	for _, player := range players {

		err = queue.ProducePlayer(queue.PlayerMessage{ID: player.ID, SkipPlayerGroups: true, SkipAchievements: true})
		err = helpers.IgnoreErrors(err, memcache.ErrInQueue)
		if err != nil {
			return err
		}

		time.Sleep((cronTime / time.Duration(toQueue)) * time.Duration(consumers))
	}

	return err
}