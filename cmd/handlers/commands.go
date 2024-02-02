package handlers

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/erazr/omnibot/config"
	"github.com/erazr/omnibot/internal/weather"
	weatherWidget "github.com/erazr/omnibot/internal/weather_widget"
)

var (
	integerOptionMinValue = 1.0

	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "current-weather",
			Description: "Current weather information and forecast up to 5 days",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "query",
					Description: "City name, zipcode, IP address, or Latitude/Longitude (decimal degree).",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days-to-forcast",
					Description: "Number of days to forecast. Value ranges from 1 to 5",
					MinValue:    &integerOptionMinValue,
					MaxValue:    5,
					Required:    false,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"current-weather": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options
			optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
			for _, opt := range options {
				optionMap[opt.Name] = opt
			}

			var weatherData *weather.WeatherResponse

			query := optionMap["query"].StringValue()
			days := optionMap["days-to-forcast"]

			if days == nil {
				weatherData, _ = weather.GetWeather(query, 4)
			} else {
				weatherData, _ = weather.GetWeather(query, days.IntValue())
			}

			reader, _ := weatherWidget.DrawWeatherWidget(weatherData)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Files: []*discordgo.File{{
						Name:        "assets/images/widget.png",
						ContentType: "image/png",
						Reader:      reader,
					}},
				},
			})
		},
	}
)

// Create new session and register commands
func RegisterCommands() error {

	cfg, err := config.LoadConfig()

	dgo, err := discordgo.New("Bot " + cfg.TOKEN)

	if err != nil {
		return err
	}

	err = dgo.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return err
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := dgo.ApplicationCommandCreate(dgo.State.User.ID, dgo.State.Application.GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
	dgo.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	dgo.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentGuilds | discordgo.IntentGuildMessages

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dgo.Close()

	return nil
}
