FROM golang:alpine

RUN apk add git 
RUN go get -u github.com/bwmarrin/discordgo
RUN go get -u github.com/ns3777k/go-shodan/shodan
RUN go get -u github.com/sirupsen/logrus
RUN go get -u github.com/JustHumanz/simplepaste

RUN mkdir /app
COPY . /app
WORKDIR /app

CMD ["go","run","."]