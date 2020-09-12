FROM golang:alpine

RUN apk add git nmap nmap-scripts tor
RUN go get -u github.com/bwmarrin/discordgo
RUN go get -u github.com/ns3777k/go-shodan/shodan
RUN go get -u github.com/sirupsen/logrus
RUN go get -u github.com/JustHumanz/simplepaste
RUN go get -u github.com/olekukonko/tablewriter
RUN go get -u github.com/Ullaakut/nmap

RUN mkdir /app
COPY . /app
WORKDIR /app
ENTRYPOINT ["go","run","."]
CMD ["nohup","tor","&"]