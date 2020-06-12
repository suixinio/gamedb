package chatbot

import (
	"regexp"

	"github.com/Jleagle/steam-go/steamapi"
	"github.com/bwmarrin/discordgo"
	"github.com/gamedb/gamedb/pkg/elastic"
	"github.com/gamedb/gamedb/pkg/mongo"
)

type CommandApp struct {
}

func (CommandApp) Regex() *regexp.Regexp {
	return regexp.MustCompile(`^[.|!](app|game) (.*)`)
}

func (CommandApp) Example() string {
	return ".game GameName"
}

func (CommandApp) Description() string {
	return "Get info on a game"
}

func (CommandApp) Type() CommandType {
	return TypeGame
}

func (c CommandApp) Output(msg *discordgo.MessageCreate) (message discordgo.MessageSend, err error) {

	matches := RegexCache[c.Regex()].FindStringSubmatch(msg.Message.Content)

	apps, _, err := elastic.SearchApps(1, 0, matches[2], nil, false, false, false)
	if err != nil {
		return message, err
	} else if len(apps) == 0 {
		message.Content = "App **" + matches[2] + "** not found"
		return message, nil
	}

	app, err := mongo.GetApp(apps[0].ID)
	if err != nil {
		return message, err
	}

	message.Content = "<@" + msg.Author.ID + ">"
	message.Embed = &discordgo.MessageEmbed{
		Title: app.GetName(),
		URL:   "https://gamedb.online" + app.GetPath(),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: app.GetHeaderImage(),
		},
		Footer: getFooter(),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Release Date",
				Value: app.GetReleaseDateNice(),
			},
			{
				Name:  "Price",
				Value: app.Prices.Get(steamapi.ProductCCUS).GetFinal(),
			},
			{
				Name:  "Review Score",
				Value: app.GetReviewScore(),
			},
			{
				Name:  "Followers",
				Value: app.GetFollowers(),
			},
		},
	}

	return message, nil
}
