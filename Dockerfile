# syntax=docker/dockerfile:1

# Frontend Build
FROM docker.io/library/node:20.9.0-alpine3.17 AS nodebuilder
WORKDIR /app
COPY . ./
RUN rm -rf ./.nuxt/ && \
    apk add --no-cache git && \
    yarn && yarn generate


# Backend Build
FROM docker.io/library/golang:1.21 AS gobuilder
WORKDIR /go/src/github.com/galexrt/fivenet/
COPY . ./
RUN apt-get update && \
    apt-get install -y git && \
    make build-go

# Final Image
FROM docker.io/library/alpine:3.18.4
WORKDIR /app
RUN apk --no-cache add ca-certificates tini tzdata && \
    mkdir -p ./.output/public
COPY --from=nodebuilder /app/.output/public ./.output/public
COPY --from=gobuilder /go/src/github.com/galexrt/fivenet/fivenet /usr/local/bin

EXPOSE 7070/tcp 8080/tcp 9090/tcp

ENTRYPOINT ["tini", "--", "fivenet"]

CMD ["server"]
