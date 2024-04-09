# syntax=docker/dockerfile:1

# Frontend Build
FROM docker.io/library/node:20.12.0-alpine3.18 AS nodebuilder
ARG NUXT_UI_PRO_LICENSE
WORKDIR /app
COPY . ./
RUN rm -rf ./.nuxt/ && \
    apk add --no-cache git && \
    yarn && NUXT_UI_PRO_LICENSE=${NUXT_UI_PRO_LICENSE} yarn generate

# Backend Build
FROM docker.io/library/golang:1.21.8 AS gobuilder
WORKDIR /go/src/github.com/galexrt/fivenet/
COPY . ./
RUN apt-get update && \
    apt-get install -y git && \
    make build-go

# Final Image
FROM docker.io/library/alpine:3.19.1
WORKDIR /app
RUN apk --no-cache add ca-certificates tini tzdata && \
    mkdir -p ./.output/public
COPY --from=nodebuilder /app/.output/public ./.output/public
COPY --from=gobuilder /go/src/github.com/galexrt/fivenet/fivenet /usr/local/bin

EXPOSE 7070/tcp 8080/tcp 9090/tcp

ENTRYPOINT ["tini", "--", "fivenet"]

CMD ["server"]
