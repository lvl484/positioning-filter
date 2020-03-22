FROM golang:1.13 as builder
RUN mkdir -p /go/src/github.com/lvl484
ENV GO111MODULE on
ENV CGO_ENABLED 0
WORKDIR /go/src/github.com/lvl484/positioning-filter
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
WORKDIR /go/src/github.com/lvl484/positioning-filter/cmd
RUN mkdir -p /opt/services/ && go build -o /opt/services/positioning-filter

FROM alpine:3.7
COPY --from=builder /opt/services/positioning-filter /opt/services/positioning-filter/positioning-filter
COPY config/viper.config.json /opt/services/positioning-filter/config/
WORKDIR /opt/services/positioning-filter
