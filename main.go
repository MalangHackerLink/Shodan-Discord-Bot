package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ns3777k/go-shodan/shodan"
	log "github.com/sirupsen/logrus"
)

var (
	client *shodan.Client
	prefix = "!>"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Error(err)
	}
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Error(err)
	}

	client = shodan.NewClient(nil, os.Getenv("SHODANAPI"))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}