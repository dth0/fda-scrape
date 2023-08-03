FROM golang:1.20.7-alpine as builder
ENV CGO_ENABLED 0

WORKDIR /src/fda-scrape

RUN	apk update && \
	apk add ca-certificates && \
	rm -rf /var/cache/apk/*

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN	mkdir build && \
	go build -o ./build/fda-scrape

FROM alpine

LABEL maintainer="daniel.theodoro@gmail.com"

RUN	apk update && \
	apk add ca-certificates && \
	apk add curl && \
	rm -rf /var/cache/apk/*

COPY --from=builder /src/fda-scrape/build/fda-scrape bin/

CMD ["fda-scrape", "serve"]
