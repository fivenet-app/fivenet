# syntax=docker/dockerfile:1
FROM docker.io/library/node:19.9.0-alpine3.17 AS nodebuilder
WORKDIR /app
COPY . ./
RUN rm -rf ./.nuxt/ && \
    yarn && yarn generate

FROM docker.io/library/golang:1.20 AS gobuilder
WORKDIR /go/src/github.com/galexrt/fivenet/
COPY . ./
RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o fivenet .

FROM docker.io/library/alpine:3.17.3
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata && \
    mkdir -p ./.output/public
COPY --from=nodebuilder /app/.output/public ./.output/public
COPY --from=gobuilder /go/src/github.com/galexrt/fivenet/fivenet /usr/local/bin

EXPOSE 8080/tcp
EXPOSE 9090/tcp

CMD ["fivenet", "server"]
