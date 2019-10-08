FROM golang:1.13 AS build

WORKDIR /go/src/github.com/micnncim/spanner-monitor
COPY go.mod go.sum ./
RUN export GO111MODULE=on
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux \
    go build -o /go/bin/spanner-monitor cmd/spanner-monitor/main.go

FROM gcr.io/distroless/base
COPY --from=build /go/bin/spanner-monitor /
CMD ["/spanner-monitor"]
