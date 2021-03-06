package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gamedb/gamedb/pkg/chatbot"
	"github.com/gamedb/gamedb/pkg/config"
	"github.com/gamedb/gamedb/pkg/consumers"
	"github.com/gamedb/gamedb/pkg/helpers"
	influxHelper "github.com/gamedb/gamedb/pkg/influx"
	"github.com/gamedb/gamedb/pkg/log"
	"github.com/gamedb/gamedb/pkg/memcache"
	"github.com/gamedb/gamedb/pkg/mongo"
	"github.com/gamedb/gamedb/pkg/mysql"
	influx "github.com/influxdata/influxdb1-client"
	"go.uber.org/zap"
)

func main() {

	err := config.Init(helpers.GetIP())
	log.InitZap(log.LogNameChatbot)
	defer log.Flush()
	if err != nil {
		log.ErrS(err)
		return
	}

	if config.IsConsumer() {
		log.Err("Prod & local only")
		return
	}

	// Profiling
	// if config.IsLocal() {
	// 	go func() {
	// 		err := http.ListenAndServe(":6061", nil)
	// 		if err != nil {
	// 			log.ErrS(err)
	// 		}
	// 	}()
	// }

	err = mysql.GetConsumer("chatbot")
	if err != nil {
		log.ErrS(err)
		return
	}

	consumers.Init(consumers.ChatbotDefinitions)

	session, err := websocketServer()
	if err != nil {
		log.FatalS(err)
	}

	refreshCommands(session)

	go updateGuildsCount(session)

	helpers.KeepAlive(
		mysql.Close,
		mongo.Close,
		memcache.Close,
		closeWebcocketConn(session),
	)
}

func closeWebcocketConn(session *discordgo.Session) func() {
	return func() {
		err := session.Close()
		if err != nil {
			log.Err("disconnecting from discord", zap.Error(err))
		}
	}
}

func refreshCommands(session *discordgo.Session) {

	apiCommands, err := session.ApplicationCommands(config.C.DiscordChatBotClientID, "")
	if err != nil {
		log.Err("ApplicationCommands()", zap.Error(err))
		return
	}

	// Delete removed commands
	for _, apiCommand := range apiCommands {
		if _, ok := chatbot.CommandCache[apiCommand.Name]; !ok {

			log.Info("Deleting dommand", zap.String("id", apiCommand.Name))
			err = session.ApplicationCommandDelete(config.C.DiscordChatBotClientID, "", apiCommand.ID)
			if err != nil {
				log.Err("ApplicationCommandDelete()", zap.Error(err), zap.String("name", apiCommand.Name))
				return
			}
		}
	}

	// Update updated commands
	for _, apiCommand := range apiCommands {
		if localCommand, ok := chatbot.CommandCache[apiCommand.Name]; ok {

			if apiCommand.Options == nil {
				apiCommand.Options = []*discordgo.ApplicationCommandOption{}
			}

			b1, _ := json.Marshal(apiCommand.Options)
			b2, _ := json.Marshal(localCommand.Slash())
			if string(b1) != string(b2) {

				log.Info("Updating command", zap.String("id", localCommand.ID()))
				command := &discordgo.ApplicationCommand{
					Name:        localCommand.ID(),
					Description: strings.ToUpper(string(localCommand.Type())) + ": " + localCommand.Description(),
					Options:     localCommand.Slash(),
				}
				_, err = session.ApplicationCommandCreate(config.C.DiscordChatBotClientID, "", command)
				if err != nil {
					log.Err("ApplicationCommandCreate()", zap.Error(err), zap.String("id", localCommand.ID()))
					return
				}
			}
		}
	}

	// Add missing commands
	for k, localCommand := range chatbot.CommandCache {
		func() {

			// Check if already exists
			for _, apiCommand := range apiCommands {
				if apiCommand.Name == k {
					return
				}
			}

			log.Info("Adding command", zap.String("id", localCommand.ID()))
			command := &discordgo.ApplicationCommand{
				Name:        localCommand.ID(),
				Description: strings.ToUpper(string(localCommand.Type())) + ": " + localCommand.Description(),
				Options:     localCommand.Slash(),
			}
			_, err = session.ApplicationCommandCreate(config.C.DiscordChatBotClientID, "", command)
			if err != nil {
				log.Err("ApplicationCommandCreate()", zap.String("id", localCommand.ID()))
				return
			}
		}()
	}
}

func updateGuildsCount(session *discordgo.Session) {

	if config.IsProd() {
		for {
			func() {
				var after = ""
				var count = 0
				for {

					guilds, err := session.UserGuilds(100, "", after)
					if err != nil {
						log.ErrS(err)
						return
					}

					for _, guild := range guilds {
						count++
						after = guild.ID
					}

					if len(guilds) < 100 {
						break
					}
				}

				// Save to Influx
				point := influx.Point{
					Measurement: influxHelper.InfluxMeasurementChatBot.String(),
					Fields: map[string]interface{}{
						"guilds": count,
					},
					Precision: "h",
				}

				_, err := influxHelper.InfluxWrite(influxHelper.InfluxRetentionPolicyAllTime, point)
				if err != nil {
					log.ErrS(err)
				}
			}()

			time.Sleep(time.Hour)
		}
	}
}
