From golang:alpine as builder
WORKDIR /build/
COPY . .
RUN CGO_ENABLED=0 go mod tidy && go build -o odin-scraper




FROM alpine
RUN apk --update add ca-certificates curl
WORKDIR /
COPY --from=builder /build/odin-scraper /odin-scraper


CMD ["/odin-scraper"]