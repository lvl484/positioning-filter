FROM golang:1.13 as modules
ADD go.mod go.sum /m/
RUN cd /m && go mod download

FROM golang:1.13 as builder
ENV CGO_ENABLED 0
RUN mkdir -p /opt/resource/
COPY --from=modules /go/pkg/ /go/pkg/
WORKDIR /opt/resource/
COPY . .
WORKDIR /opt/resource/cmd
RUN go build -o positioning-filter
RUN ./opt/resource/cmd/positioning-filter

FROM alpine:3.7
COPY --from=builder /opt/services/positioning-filter ./opt/services/positioning-filter
CMD ./opt/services/positioning-filter