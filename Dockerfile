FROM golang:1.8

RUN mkdir -p /go/src/nos-modem-alcatel-mw40v-prometheus-exporther
WORKDIR /go/src/nos-modem-alcatel-mw40v-prometheus-exporther
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.BUILD_DATE=$(date -u '+%Y-%m-%d_%H:%M:%S') -X main.GIT_HASH=$(git rev-parse HEAD) -X main.GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD) -linkmode external -extldflags -static" -a nos-modem-alcatel-mw40v-prometheus-exporther.go

FROM scratch
COPY --from=0 /go/src/nos-modem-alcatel-mw40v-prometheus-exporther/nos-modem-alcatel-mw40v-prometheus-exporther /nos-modem-alcatel-mw40v-prometheus-exporther
ENTRYPOINT ["/nos-modem-alcatel-mw40v-prometheus-exporther"]