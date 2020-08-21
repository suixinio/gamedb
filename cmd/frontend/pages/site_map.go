package pages

import (
	"net/http"
	"time"

	"github.com/Jleagle/sitemap-go/sitemap"
	"github.com/gamedb/gamedb/pkg/config"
	"github.com/gamedb/gamedb/pkg/helpers"
	"github.com/gamedb/gamedb/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

//noinspection GoUnusedParameter
func SiteMapIndexHandler(w http.ResponseWriter, r *http.Request) {

	var sitemaps = []string{
		"/sitemap-pages.xml",
		"/sitemap-games-by-score.xml",
		"/sitemap-games-by-players.xml",
		"/sitemap-games-new.xml",
		"/sitemap-games-upcoming.xml",
		"/sitemap-players-by-level.xml",
		"/sitemap-players-by-games.xml",
		"/sitemap-groups.xml",
		"/sitemap-badges.xml",
	}

	sm := sitemap.NewSiteMapIndex()

	for _, v := range sitemaps {
		sm.AddSitemap(config.Config.GameDBDomain.Get()+v, time.Time{})
	}

	_, err := sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}

//noinspection GoUnusedParameter
func SiteMapPagesHandler(w http.ResponseWriter, r *http.Request) {

	var pages = []string{
		"/",
		"/achievements",
		"/api",
		"/badges",
		"/changes",
		"/chat",
		"/commits",
		"/contact",
		"/developers",
		"/donate",
		"/experience",
		"/games",
		"/games/achievements",
		"/games/compare",
		"/games/coop",
		"/games/new-releases",
		"/games/random",
		"/games/sales",
		"/games/trending",
		"/games/upcoming",
		"/games/wishlist",
		"/genres",
		"/groups",
		"/info",
		"/login",
		"/news",
		"/packages",
		"/players",
		"/price-changes",
		"/product-keys",
		"/publishers",
		"/stats",
		"/steam-api",
		"/tags",
	}

	sm := sitemap.NewSitemap()

	for _, v := range pages {
		sm.AddLocation(config.Config.GameDBDomain.Get()+v, time.Time{}, sitemap.FrequencyHourly, 1)
	}

	_, err := sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}

func SiteMapGamesByPlayersHandler(w http.ResponseWriter, r *http.Request) {

	apps, err := mongo.GetApps(0, 1000, bson.D{{"player_peak_week", -1}}, bson.D{}, bson.M{"_id": 1, "name": 1, "updated_at": 1})
	if err != nil {
		zap.S().Error(err)
		return
	}

	sm := sitemap.NewSitemap()
	for _, app := range apps {
		sm.AddLocation(config.Config.GameDBDomain.Get()+app.GetPath(), app.UpdatedAt, sitemap.FrequencyWeekly, 0.9)
	}

	_, err = sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}

func SiteMapGamesByScoreHandler(w http.ResponseWriter, r *http.Request) {

	apps, err := mongo.GetApps(0, 1000, bson.D{{"reviews_score", -1}}, bson.D{}, bson.M{"_id": 1, "name": 1, "updated_at": 1})
	if err != nil {
		zap.S().Error(err)
		return
	}

	sm := sitemap.NewSitemap()
	for _, app := range apps {
		sm.AddLocation(config.Config.GameDBDomain.Get()+app.GetPath(), app.UpdatedAt, sitemap.FrequencyWeekly, 0.9)
	}

	_, err = sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}

func SiteMapGamesUpcomingHandler(w http.ResponseWriter, r *http.Request) {

	apps, err := mongo.GetApps(0, 1000, bson.D{{"release_date_unix", 1}}, upcomingFilter, bson.M{"_id": 1, "name": 1, "updated_at": 1})
	if err != nil {
		zap.S().Error(err)
		return
	}

	sm := sitemap.NewSitemap()
	for _, app := range apps {
		sm.AddLocation(config.Config.GameDBDomain.Get()+app.GetPath(), app.UpdatedAt, sitemap.FrequencyWeekly, 0.9)
	}

	_, err = sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}

func SiteMapGamesNewHandler(w http.ResponseWriter, r *http.Request) {

	var filter = bson.D{
		{"release_date_unix", bson.M{"$lt": time.Now().Unix()}},
		{"release_date_unix", bson.M{"$gt": time.Now().AddDate(0, 0, -config.Config.NewReleaseDays.GetInt()).Unix()}},
	}

	apps, err := mongo.GetApps(0, 1000, bson.D{{"release_date_unix", -1}}, filter, bson.M{"_id": 1, "name": 1, "updated_at": 1})
	if err != nil {
		zap.S().Error(err)
		return
	}

	sm := sitemap.NewSitemap()
	for _, app := range apps {
		sm.AddLocation(config.Config.GameDBDomain.Get()+app.GetPath(), app.UpdatedAt, sitemap.FrequencyWeekly, 0.9)
	}

	_, err = sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}

//noinspection GoUnusedParameter
func SiteMapPlayersByLevel(w http.ResponseWriter, r *http.Request) {

	sm := sitemap.NewSitemap()

	players, err := mongo.GetPlayers(0, 1000, bson.D{{Key: "level", Value: -1}}, nil, bson.M{"_id": 1, "persona_name": 1, "updated_at": 1})
	if err != nil {
		zap.S().Error(err)
	}

	for _, player := range players {
		sm.AddLocation(config.Config.GameDBDomain.Get()+player.GetPath(), player.UpdatedAt, sitemap.FrequencyWeekly, 0.9)
	}

	_, err = sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}

//noinspection GoUnusedParameter
func SiteMapPlayersByGamesCount(w http.ResponseWriter, r *http.Request) {

	sm := sitemap.NewSitemap()

	players, err := mongo.GetPlayers(0, 1000, bson.D{{Key: "games_count", Value: -1}}, nil, bson.M{"_id": 1, "persona_name": 1, "updated_at": 1})
	if err != nil {
		zap.S().Error(err)
	}

	for _, player := range players {
		sm.AddLocation(config.Config.GameDBDomain.Get()+player.GetPath(), player.UpdatedAt, sitemap.FrequencyWeekly, 0.9)
	}

	_, err = sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}

//noinspection GoUnusedParameter
func SiteMapGroups(w http.ResponseWriter, r *http.Request) {

	sm := sitemap.NewSitemap()

	groups, err := mongo.GetGroups(1000, 0, bson.D{{Key: "members", Value: -1}}, bson.D{{Key: "type", Value: helpers.GroupTypeGroup}}, bson.M{"_id": 1, "name": 1, "updated_at": 1})
	if err != nil {
		zap.S().Error(err)
	}

	for _, v := range groups {
		sm.AddLocation(config.Config.GameDBDomain.Get()+v.GetPath(), v.UpdatedAt, sitemap.FrequencyWeekly, 0.9)
	}

	_, err = sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}

func SiteMapBadges(w http.ResponseWriter, r *http.Request) {

	sm := sitemap.NewSitemap()

	for _, badge := range helpers.BuiltInSpecialBadges {
		sm.AddLocation(config.Config.GameDBDomain.Get()+badge.GetPath(false), time.Time{}, sitemap.FrequencyWeekly, 0.9)
	}
	for _, badge := range helpers.BuiltInEventBadges {
		sm.AddLocation(config.Config.GameDBDomain.Get()+badge.GetPath(false), time.Time{}, sitemap.FrequencyWeekly, 0.9)
	}

	_, err := sm.Write(w)
	if err != nil {
		zap.S().Error(err)
	}
}