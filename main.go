package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Token string `yaml:"token"`
}

func main() {
	config := loadConfig("config.yml")

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		fmt.Println(m, m.Content, "ping received")
		if m.Content == "ping" {
			s.ChannelMessageSend(m.ChannelID, "pong")
		}
		fmt.Println("message created at", m.Timestamp, "from", m.Author.Username)
	})

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
	}
	defer dg.Close()
	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func loadConfig(filename string) Config {
	var config Config
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening config file:", err)
		return config
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("error decoding config file:", err)
		return config
	}

	return config
}
