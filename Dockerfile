FROM golang:1.13 as modules
ADD ./go.mod ./go.sum /m/
RUN cd /m && go mod download
FROM golang:1.13 as builder

RUN mkdir -p /opt/resource/

COPY --from=modules /go/pkg/ /go/pkg/

WORKDIR /opt/resource/
COPY cmd             cmd
COPY config          config
COPY kafka           kafka
COPY filter          filter
COPY matcher         matcher
COPY storage         storage
COPY web             web

WORKDIR /opt/resource/cmd/
RUN go build -o /opt/services/positioning-filter .

FROM alpine:3.7
COPY --from=builder /opt/services/positioning-filter /opt/services/positioning-filter
CMD /opt/services/positioning-filter