package mysql

import (
	"strconv"
	"time"

	"github.com/Jleagle/steam-go/steamapi"
	"github.com/dustin/go-humanize"
	"github.com/gamedb/gamedb/pkg/helpers"
)

type Tag struct {
	ID        int        `gorm:"not null;primary_key"`
	CreatedAt time.Time  `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt *time.Time `gorm:""`
	Name      string     `gorm:"not null;index:name"`
	Apps      int        `gorm:"not null"`
	MeanPrice string     `gorm:"not null"` // map[steamapi.CountryCode]float64
	MeanScore float64    `gorm:"not null"`
}

func (tag Tag) GetPath() string {
	return "/games?tags=" + strconv.Itoa(tag.ID)
}

func (tag Tag) GetName() (name string) {

	if tag.Name == "" {
		return "Tag " + humanize.Comma(int64(tag.ID))
	}

	return tag.Name
}

func (tag Tag) GetMeanPrice(code steamapi.ProductCC) (string, error) {
	return GetMeanPrice(code, tag.MeanPrice)
}

func (tag Tag) GetMeanScore() string {
	return helpers.FloatToString(tag.MeanScore, 2) + "%"
}

func GetAllTags() (tags []Tag, err error) {

	db, err := GetMySQLClient()
	if err != nil {
		return tags, err
	}

	db = db.Find(&tags)
	if db.Error != nil {
		return tags, db.Error
	}

	return tags, nil
}

func GetTagsByID(ids []int, columns []string) (tags []Tag, err error) {

	if len(ids) == 0 {
		return tags, err
	}

	db, err := GetMySQLClient()
	if err != nil {
		return tags, err
	}

	if len(columns) > 0 {
		db = db.Select(columns)
	}

	db = db.Where("id IN (?)", ids)
	db = db.Order("name ASC")
	db = db.Limit(100)
	db = db.Find(&tags)

	return tags, db.Error
}
