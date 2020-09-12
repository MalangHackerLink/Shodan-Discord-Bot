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
	if strings.HasPrefix(m.Content, prefix) {
		log.WithFields(log.Fields{
			"UserID": m.Author.ID,
		}).Info(m.Content[len(prefix):])

		array := strings.Split(m.Content, " ")
		if array[0] == prefix+"sub" {
			info, err := client.GetDomain(ctx, array[1])
			if err != nil {
				log.Error(err)
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
			_, err = s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\nSubs domain for `"+array[1]+"` : "+strings.Join(info.Subdomains, ","))
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\nSubs domain for `"+array[1]+"` : "+PushPastebin(array[1], []byte(strings.Join(info.Subdomains, ","))))
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
				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\nReverse ip for `"+ips[i].String()+"` : "+strings.Join(result, ","))
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
				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\nResolve domain for `"+domainlist[i]+"` : "+info[domainlist[i]].String())
			}
		} else if array[0] == prefix+"host" {
			s.ChannelMessageSend(m.ChannelID, "Wait for 1 minute,if the result still not upper that's mean i fucked")
			hostlist := strings.Split(array[1], ",")
			for i := 0; i < len(hostlist); i++ {
				info, err := client.GetServicesForHost(ctx, hostlist[i], nil)
				if err != nil {
					log.Error(err)
					s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\n"+err.Error())
					return
				}
				e, err := json.Marshal(info)
				if err != nil {
					log.Error(err)
				}
				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\n"+PushPastebin(hostlist[i], e))
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

type IPPORT struct {
	IP   string
	Port []string
}

func Map(s *discordgo.Session, m *discordgo.MessageCreate) {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	if strings.HasPrefix(m.Content, prefix2) {
		log.WithFields(log.Fields{
			"UserID": m.Author.ID,
		}).Info(m.Content[len(prefix2):])
		array := strings.Split(m.Content, " ")
		if array[0] == prefix2+"scan-tcp" || array[0] == prefix2+"scan-udp" {
			s.ChannelMessageSend(m.ChannelID, "Wait for 3-5 minute,if the result still not upper that's mean i fucked")
			var (
				host     = array[1]
				dat      [][]string
				warnings []string
				err      error
			)
			portlist := strings.Split(array[2], ",")
			Data := IPPORT{
				IP:   host,
				Port: portlist,
			}
			if array[0][len(prefix2+"scan-"):] == "tcp" {
				dat, warnings, err = Data.ScanGoBrrrr("tcp")
				if err != nil {
					log.Error(err)
					s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\n"+err.Error()+"\nWarnings "+strings.Join(warnings, " "))
					return
				}
			} else {
				dat, warnings, err = Data.ScanGoBrrrr("udp")
				if err != nil {
					log.Error(err)
					s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\n"+err.Error()+"\nWarnings "+strings.Join(warnings, " "))
					return
				}
			}
			table.SetHeader([]string{"Port", "Status", "Service", "Reason"})
			for _, v := range dat {
				table.Append(v)
			}
			table.Render()
			if warnings != nil {
				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\nWARNINGS `"+strings.Join(warnings, " ")+"`\nPort for `"+host+"`\n```\r"+tableString.String()+"```")
			} else {
				s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\nPort for `"+host+":`\n```\r"+tableString.String()+"```")
			}
		} else if array[0] == prefix2+"script" {
			if len(array) > 1 {
				s.ChannelMessageSend(m.ChannelID, "Wait for 3-5 minute,if the result still not upper that's mean i fucked")
				Data := IPPORT{
					IP: array[2],
				}
				dat, warn := Data.ScanScriptBrrr(array[1])
				text := strings.Join(dat, "\n")
				if warn != nil {
					log.Warn(warn)
				} else {
					warn = append(warn, "null")
				}
				if len(text) > 2000 {
					s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\n```"+PushPastebin(Data.IP, []byte(text))+"```")
				} else {
					s.ChannelMessageSend(m.ChannelID, "<@"+m.Author.ID+">\nWarning\n "+strings.Join(warn, " ")+"```"+text+"```")
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "invalid script args")
			}
		}
	}
}

func (Data IPPORT) ScanGoBrrrr(pkttype string) ([][]string, []string, error) {
	var (
		tbl     [][]string
		scanner *nmap.Scanner
		err     error
	)
	if pkttype == "tcp" {
		scanner, err = nmap.NewScanner(
			nmap.WithTargets(Data.IP),
			nmap.WithPorts(strings.Join(Data.Port, ",")),
			nmap.WithProxies(TORSOCKS),
			nmap.WithReason(),
			nmap.WithDefaultScript(),
			nmap.WithContext(ctx),
		)
		if err != nil {
			log.Error(err)
			return nil, nil, err
		}

		result, warnings, err := scanner.Run()
		if err != nil {
			log.Error(err, warnings)
			return nil, warnings, err
		}

		for _, host := range result.Hosts {
			if len(host.Ports) == 0 || len(host.Addresses) == 0 {
				continue
			}
			for _, port := range host.Ports {
				tmp := [][]string{[]string{strconv.Itoa(int(port.ID)) + "/" + port.Protocol, port.State.String(), port.Service.Name, host.Status.Reason}}
				tbl = append(tbl, tmp...)
			}
		}
		return tbl, warnings, nil
	} else if pkttype == "udp" {
		scanner, err = nmap.NewScanner(
			nmap.WithUDPScan(),
			nmap.WithTargets(Data.IP),
			nmap.WithPorts(strings.Join(Data.Port, ",")),
			nmap.WithProxies(TORSOCKS),
			nmap.WithReason(),
			nmap.WithDefaultScript(),
			nmap.WithContext(ctx),
		)
		if err != nil {
			log.Error(err)
			return nil, nil, err
		}

		result, warnings, err := scanner.Run()
		if err != nil {
			log.Error(err, warnings)
			return nil, warnings, err
		}

		for _, host := range result.Hosts {
			if len(host.Ports) == 0 || len(host.Addresses) == 0 {
				continue
			}
			for _, port := range host.Ports {
				tmp := [][]string{[]string{strconv.Itoa(int(port.ID)) + "/" + port.Protocol, port.State.String(), port.Service.Name, host.Status.Reason}}
				tbl = append(tbl, tmp...)

			}
		}
		return tbl, warnings, nil
	}
	return nil, nil, nil
}

func (Data IPPORT) ScanScriptBrrr(script string) ([]string, []string) {
	var fix []string
	scanner, err := nmap.NewScanner(
		nmap.WithScripts(script),
		nmap.WithTargets(Data.IP),
		nmap.WithProxies(TORSOCKS),
		nmap.WithContext(ctx),
	)
	if err != nil {
		log.Error(err)
	}

	result, warnings, err := scanner.Run()
	for _, res := range result.Hosts {
		for _, pors := range res.Ports {
			for _, scr := range pors.Scripts {
				fix = append(fix, scr.Output)
			}
		}
	}
	return fix, warnings
}
