FROM golang:1.16.5-alpine

WORKDIR /gokeeper

COPY app/ ./app
COPY config.yml ./
COPY gqlgen.yml ./
COPY run.sh ./
COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go build -o gokeeper app/main.go

EXPOSE 8080

CMD [ "./gokeeper" ]

