FROM golang:1.20.5-alpine as builder
ENV CGO_ENABLED 0

WORKDIR /src/fda-scrape

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN mkdir build && \
	go build -o ./build/ ./cmd/fda-scrape

FROM alpine

LABEL maintainer="daniel.theodoro@gmail.com"

COPY --from=builder /src/fda-scrape/build/fda-scrape bin/

ENTRYPOINT ["fda-scrape", "serve"]
