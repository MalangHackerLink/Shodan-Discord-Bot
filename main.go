package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	"github.com/ns3777k/go-shodan/shodan" // go modules required
)

var (
	client *shodan.Client
	prefix = "!>"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println(err)
	}
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	client = shodan.NewClient(nil, os.Getenv("SHODANAPI"))
	/*
		info, err := client.GetDomain(context.Background(), "justhumanz.me")
		if err != nil {
			log.Println(err)
		}
		for i := 0; i < len(info.Subdomains); i++ {
			fmt.Println(info.Subdomains[i])
		}
		var ips []net.IP
		ip := net.ParseIP("45.76.205.91")
		ip2 := net.ParseIP("74.125.200.139")
		ips = []net.IP{ip, ip2}
		info2, err := client.GetDNSReverse(context.Background(), ips)
		if err != nil {
			log.Println(err)
		}
		for i := 0; i < len(ips); i++ {
			fmt.Println(info2[ips[i].String()])
		}
	*/

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, prefix) {
		array := strings.Split(m.Content, " ")
		if array[0] == prefix+"pong" {
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		} else if array[0] == prefix+"sub" {
			info, err := client.GetDomain(context.Background(), array[1])
			if err != nil {
				log.Error(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			var domain []string
			for i := 0; i < len(info.Subdomains); i++ {
				domain = append(domain, info.Subdomains[i])
			}
			s.ChannelMessageSend(m.ChannelID, strings.Join(domain, ","))
		} else if array[0] == prefix+"rev" {
			iplist := strings.Split(array[1], ",")
			var (
				ips    []net.IP
				result []string
			)
			for i := 0; i < len(iplist); i++ {
				ip := net.ParseIP(iplist[i])
				ips = append(ips, ip)
			}
			info2, err := client.GetDNSReverse(context.Background(), ips)
			if err != nil {
				log.Error(err)
			}
			for i := 0; i < len(ips); i++ {
				for _, res := range *info2[ips[i].String()] {
					result = append(result, res)
				}
			}
			s.ChannelMessageSend(m.ChannelID, strings.Join(result, ","))
		}
	}
}
