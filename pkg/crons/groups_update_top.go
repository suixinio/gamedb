package crons

import (
	"github.com/gamedb/gamedb/pkg/consumers"
	"github.com/gamedb/gamedb/pkg/helpers"
	"github.com/gamedb/gamedb/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

type GroupsUpdateTop struct {
	BaseTask
}

func (c GroupsUpdateTop) ID() string {
	return "queue-player-groups"
}

func (c GroupsUpdateTop) Name() string {
	return "Queue player groups"
}

func (c GroupsUpdateTop) Group() TaskGroup {
	return TaskGroupGroups
}

func (c GroupsUpdateTop) Cron() TaskTime {
	return CronTimeQueuePlayerGroups
}

func (c GroupsUpdateTop) work() (err error) {

	var filter = bson.D{
		{Key: "type", Value: helpers.GroupTypeGroup},
	}

	var sorts = []bson.D{
		{{"members", -1}},
		{{"trending", 1}},
		{{"trending", -1}},
	}

	var groupMap = map[string]bool{}

	for _, sort := range sorts {

		groups, err := mongo.GetGroups(0, 10000, sort, filter, bson.M{"_id": 1})
		if err != nil {
			return err
		}

		for _, group := range groups {
			groupMap[group.ID] = true
		}
	}

	for groupID := range groupMap {

		err = consumers.ProduceGroup(consumers.GroupMessage{ID: groupID})
		err = helpers.IgnoreErrors(err, consumers.ErrInQueue)
		if err != nil {
			return err
		}
	}

	return nil
}
