package main

import (
	"encoding/json"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/JustHumanz/simplepaste"
	"github.com/Ullaakut/nmap"
	"github.com/bwmarrin/discordgo"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

func Msg(s *discordgo.Session, m *discordgo.MessageCreate) {
	log.WithFields(log.Fields{
		"UserID": m.Author.ID,
	}).Info(m.Content)

	if strings.HasPrefix(m.Content, prefix) {
		array := strings.Split(m.Content, " ")
		if array[0] == prefix+"sub" {
			info, err := client.GetDomain(ctx, array[1])
			if err != nil {
				log.Error(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			_, err = s.ChannelMessageSend(m.ChannelID, "Subs domain for `"+array[1]+"` : "+strings.Join(info.Subdomains, ","))
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "Subs domain for `"+array[1]+"` : "+PushPastebin(array[1], []byte(strings.Join(info.Subdomains, ","))))
				return
			}
		} else if array[0] == prefix+"rev" {
			iplist := strings.Split(array[1], ",")
			var (
				ips []net.IP
			)
			for i := 0; i < len(iplist); i++ {
				ip := net.ParseIP(iplist[i])
				ips = append(ips, ip)
			}
			info2, err := client.GetDNSReverse(ctx, ips)
			if err != nil {
				log.Error(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
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
			info, err := client.GetDNSResolve(ctx, domainlist)
			if err != nil {
				log.Error(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			for i := 0; i < len(domainlist); i++ {
				s.ChannelMessageSend(m.ChannelID, "Resolve domain for `"+domainlist[i]+"` : "+info[domainlist[i]].String())
			}
		} else if array[0] == prefix+"host" {
			s.ChannelMessageSend(m.ChannelID, "Wait for 1 minute,if the result still not upper that's mean i fucked")
			hostlist := strings.Split(array[1], ",")
			for i := 0; i < len(hostlist); i++ {
				info, err := client.GetServicesForHost(ctx, hostlist[i], nil)
				if err != nil {
					log.Error(err)
					s.ChannelMessageSend(m.ChannelID, err.Error())
					return
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
	api := simplepaste.NewAPI(PASTEBIN)
	paste := simplepaste.NewPaste(title, string(body))
	paste.ExpireDate = simplepaste.Day
	link, err := api.SendPaste(paste)
	if err != nil {
		log.Error(err)
	}
	return link
}

func Map(s *discordgo.Session, m *discordgo.MessageCreate) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	if strings.HasPrefix(m.Content, prefix2) {
		array := strings.Split(m.Content, " ")
		if array[0] == prefix2+"scan" {
			s.ChannelMessageSend(m.ChannelID, "Wait for 3-5 minute,if the result still not upper that's mean i fucked")
			host := array[1]
			portlist := strings.Split(array[2], ",")
			Data := IPPORT{
				IP:   host,
				Port: portlist,
			}
			dat, warnings, err := Data.ScanGobrrrr()
			if err != nil {
				log.Error(err)
				s.ChannelMessageSend(m.ChannelID, err.Error()+"\nWarnings "+strings.Join(warnings, " "))
				return
			}
			table.SetHeader([]string{"Port", "Status", "Service"})
			for _, v := range dat {
				table.Append(v)
			}
			table.Render()
			if warnings != nil {
				s.ChannelMessageSend(m.ChannelID, "WARNINGS `"+strings.Join(warnings, " ")+"`\nPort for `"+host+"`\n```\r"+tableString.String()+"```")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Port for `"+host+":`\n```\r"+tableString.String()+"```")
			}
		}
	}
}

type IPPORT struct {
	IP   string
	Port []string
}

func (Data IPPORT) ScanGobrrrr() ([][]string, []string, error) {
	var tbl [][]string
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(Data.IP),
		nmap.WithPorts(strings.Join(Data.Port, ",")),
		nmap.WithProxies(TORSOCKS),
		nmap.WithContext(ctx),
	)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}
		for _, port := range host.Ports {
			tmp := [][]string{[]string{strconv.Itoa(int(port.ID)) + "/" + port.Protocol, port.State.String(), port.Service.Name}}
			tbl = append(tbl, tmp...)

		}
	}
	return tbl, warnings, nil
}
