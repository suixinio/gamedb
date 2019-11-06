package mongo

import (
	"strconv"
	"time"

	"github.com/gamedb/gamedb/pkg/helpers"
	"github.com/gamedb/gamedb/pkg/log"
	. "go.mongodb.org/mongo-driver/bson"
)

type Change struct {
	ID        int       `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	Apps      []int     `bson:"apps"`
	Packages  []int     `bson:"packages"`
}

func (change Change) BSON() D {

	// Apps
	var apps A
	for _, v := range change.Apps {
		apps = append(apps, v)
	}

	// Packages
	var packages A
	for _, v := range change.Packages {
		packages = append(packages, v)
	}

	// BSON
	return D{
		{"_id", change.ID},
		{"created_at", change.CreatedAt},
		{"apps", apps},
		{"packages", packages},
	}
}

func (change Change) GetName() (name string) {

	return "Change " + strconv.Itoa(change.ID)
}

func (change Change) GetPath() string {
	return "/changes/" + strconv.Itoa(change.ID)
}

func (change Change) GetTimestamp() int64 {
	return change.CreatedAt.Unix()
}

func (change Change) GetNiceDate() string {
	return change.CreatedAt.Format(helpers.DateYearTime)
}

func (change Change) OutputForJSON(allApps map[int]string, allPackages map[int]string) (output []interface{}) {

	var apps = map[int]changeProduct{}
	var packages = map[int]changeProduct{}

	for _, v := range change.Apps {
		if val, ok := allApps[v]; ok {
			apps[v] = changeProduct{
				Name: val,
				Path: helpers.GetAppPath(v, val),
			}
		}
	}

	for _, v := range change.Packages {
		if val, ok := allPackages[v]; ok {
			packages[v] = changeProduct{
				Name: val,
				Path: helpers.GetPackagePath(v, val),
			}
		}
	}

	return []interface{}{
		change.ID,
		change.CreatedAt.Unix(),
		change.CreatedAt.Format(helpers.DateYearTime),
		apps,
		packages,
		change.GetPath(),
		change.GetName(),
	}
}

type changeProduct struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

func GetChange(id int64) (change Change, err error) {

	var item = helpers.MemcacheChange(id)

	err = helpers.GetMemcache().GetSetInterface(item.Key, item.Expiration, &change, func() (interface{}, error) {

		var change Change

		err = FindOne(CollectionChanges, D{{"_id", id}}, nil, nil, &change)
		if change.ID == 0 {
			return change, ErrNoDocuments
		}

		return change, err
	})

	return change, err
}

func GetChanges(offset int64) (changes []Change, err error) {

	var sort = D{{"_id", -1}}

	cur, ctx, err := Find(CollectionChanges, offset, 100, sort, nil, nil, nil)
	if err != nil {
		return changes, err
	}

	defer func() {
		err = cur.Close(ctx)
		log.Err(err)
	}()

	for cur.Next(ctx) {

		var change Change
		err := cur.Decode(&change)
		log.Err(err)
		changes = append(changes, change)
	}

	return changes, cur.Err()
}
