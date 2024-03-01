From golang:alpine as builder
WORKDIR /go/build/
COPY . .
RUN CGO_ENABLED=0 go build -o odin-backend




FROM alpine
RUN apk --update add ca-certificates curl
WORKDIR /
COPY --from=builder /go/build/odin-backend /odin-backend


CMD ["/odin-backend"]