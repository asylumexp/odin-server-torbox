From golang:alpine as builder
WORKDIR /build/
COPY go.mod .
COPY go.sum .
RUN CGO_ENABLED=0 GOOS=linux go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o odin-backend


FROM alpine
RUN apk --update add ca-certificates curl mailcap
WORKDIR /
COPY --from=builder /build/odin-backend /odin-backend


CMD ["/odin-backend", "serve", "--http=0.0.0.0:8090"]