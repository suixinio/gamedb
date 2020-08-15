package tasks

import (
	"encoding/json"

	"github.com/Jleagle/steam-go/steamapi"
	"go.uber.org/zap"
)

type statsRow struct {
	name       string
	count      int
	totalPrice map[steamapi.ProductCC]int
	totalScore float64
}

func (t statsRow) getMeanPrice() string {

	means := map[steamapi.ProductCC]float64{}

	for code, total := range t.totalPrice {
		means[code] = float64(total) / float64(t.count)
	}

	bytes, err := json.Marshal(means)
	zap.S().Error(err)

	return string(bytes)
}

func (t statsRow) getMeanScore() float64 {
	return t.totalScore / float64(t.count)
}
