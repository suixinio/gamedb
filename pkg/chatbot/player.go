package chatbot

import (
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"github.com/gamedb/gamedb/pkg/helpers"
	"github.com/gamedb/gamedb/pkg/helpers/memcache"
	"github.com/gamedb/gamedb/pkg/log"
	"github.com/gamedb/gamedb/pkg/mongo"
	"github.com/gamedb/gamedb/pkg/queue"
	"go.mongodb.org/mongo-driver/bson"
)

type CommandPlayer struct {
}

func (CommandPlayer) Regex() *regexp.Regexp {
	return regexp.MustCompile(`^[.|!](player|user) (.{2,32})$`)
}

func (c CommandPlayer) Output(msg *discordgo.MessageCreate) (message discordgo.MessageSend, err error) {

	matches := c.Regex().FindStringSubmatch(msg.Message.Content)

	player, q, err := mongo.SearchPlayer(matches[2], bson.M{"_id": 1, "persona_name": 1, "avatar": 1, "level": 1, "games_count": 1, "play_time": 1, "friends_count": 1})
	if err == mongo.ErrNoDocuments {

		message.Content = "Player **" + matches[2] + "** not found, please enter a user's vanity URL"
		return message, nil

	} else if err != nil {
		return message, err
	}

	if q {
		err = queue.ProducePlayer(queue.PlayerMessage{ID: player.ID})
		err = helpers.IgnoreErrors(err, memcache.ErrInQueue)
		log.Err(err)
	}

	avatar := player.GetAvatar()
	if strings.HasPrefix(avatar, "/") {
		avatar = "https://gamedb.online" + avatar
	}

	message.Content = "<@" + msg.Author.ID + ">"
	message.Embed = &discordgo.MessageEmbed{
		Title: player.GetName(),
		URL:   "https://gamedb.online" + player.GetPath(),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatar,
		},
		Footer: getFooter(),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Level",
				Value: humanize.Comma(int64(player.Level)),
			},
			{
				Name:  "Games",
				Value: humanize.Comma(int64(player.GamesCount)),
			},
			{
				Name:  "Playtime",
				Value: helpers.GetTimeLong(player.PlayTime, 3),
			},
			{
				Name:  "Friends",
				Value: humanize.Comma(int64(player.FriendsCount)),
			},
		},
	}

	return message, nil
}

func (CommandPlayer) Example() string {
	return ".player {player_name}"
}

func (CommandPlayer) Description() string {
	return "Get info on a player"
}

func (CommandPlayer) Type() CommandType {
	return TypePlayer
}
