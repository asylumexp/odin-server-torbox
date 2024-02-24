From golang:alpine as builder
WORKDIR /go/build/
COPY . .
RUN CGO_ENABLED=0 go build




FROM alpine
RUN apk --update add ca-certificates curl
WORKDIR /
COPY --from=builder /go/build/odin-scraper /odin-scraper


CMD ["/odin-scraper"]