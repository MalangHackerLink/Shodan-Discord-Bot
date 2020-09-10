package main

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"strings"
	"time"

	"github.com/JustHumanz/simplepaste"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if strings.HasPrefix(m.Content, prefix) {
		array := strings.Split(m.Content, " ")
		if array[0] == prefix+"sub" {
			info, err := client.GetDomain(context.Background(), array[1])
			if err != nil {
				log.Error(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			s.ChannelMessageSend(m.ChannelID, "Subs domain for `"+array[1]+"` : "+strings.Join(info.Subdomains, ","))
		} else if array[0] == prefix+"rev" {
			iplist := strings.Split(array[1], ",")
			var (
				ips []net.IP
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
				var (
					result []string
				)
				for _, res := range *info2[ips[i].String()] {
					result = append(result, res)

				}
				s.ChannelMessageSend(m.ChannelID, "Reverse ip for `"+ips[i].String()+"` : "+strings.Join(result, ","))
			}
		} else if array[0] == prefix+"res" {
			domainlist := strings.Split(array[1], ",")
			info, err := client.GetDNSResolve(context.Background(), domainlist)
			if err != nil {
				log.Error(err)
			}
			for i := 0; i < len(domainlist); i++ {
				s.ChannelMessageSend(m.ChannelID, "Resolve domain for `"+domainlist[i]+"` : "+info[domainlist[i]].String())
			}
		} else if array[0] == prefix+"host" {
			s.ChannelMessageSend(m.ChannelID, "Wait for 1 minute,if still not upper that's mean i fucked")
			hostlist := strings.Split(array[1], ",")
			for i := 0; i < len(hostlist); i++ {
				info, err := client.GetServicesForHost(context.Background(), hostlist[i], nil)
				if err != nil {
					log.Error(err)
				}
				e, err := json.Marshal(info)
				if err != nil {
					log.Error(err)
				}
				s.ChannelMessageSend(m.ChannelID, PushPastebin(hostlist[i], e))
				time.Sleep(30 * time.Second)
			}
		} else {
			s.ChannelMessageSend(m.ChannelID, "???")
		}
	}
}

func PushPastebin(title string, body []byte) string {
	api := simplepaste.NewAPI(os.Getenv("PASTEBIN"))
	paste := simplepaste.NewPaste(title, string(body))
	paste.ExpireDate = simplepaste.Day
	link, err := api.SendPaste(paste)
	if err != nil {
		log.Error(err)
	}
	return link
}
