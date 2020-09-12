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
	client   *shodan.Client
	prefix   = "shodan>"
	prefix2  = "nmap>"
	TORSOCKS string
	PASTEBIN string
	DISCORD  string
	SHODAN   string
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.Info("Starting bot...")
	TORSOCKS = os.Getenv("TORSOCKS")
	PASTEBIN = os.Getenv("PASTEBIN")
	DISCORD = os.Getenv("DISCORD")
	SHODAN = os.Getenv("SHODAN")

	if TORSOCKS == "" || PASTEBIN == "" || DISCORD == "" || SHODAN == "" {
		log.Error("TORSOCKS,PASTEBIN,DISCORD,SHODAN nill")
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + DISCORD)
	if err != nil {
		log.Error(err)
	}
	go dg.AddHandler(Msg)
	go dg.AddHandler(Map)

	err = dg.Open()
	if err != nil {
		log.Error(err)
	}

	client = shodan.NewClient(nil, SHODAN)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
