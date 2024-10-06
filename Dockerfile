# syntax=docker/dockerfile:1

# Frontend Build
FROM docker.io/library/node:20.18.0-alpine3.20 AS nodebuilder
ARG NUXT_UI_PRO_LICENSE
WORKDIR /app
COPY . ./
RUN rm -rf ./.nuxt/ && \
    find ./public/images/livemap/ \
        ! -path '*/tiles*' -and ! -path './public/images/livemap/' \
        -exec rm -rf {} + && \
    apk add --no-cache git && \
    corepack enable && \
    corepack prepare pnpm@latest --activate && \
    pnpm install && \
    NUXT_UI_PRO_LICENSE=${NUXT_UI_PRO_LICENSE} pnpm generate

# Backend Build
FROM docker.io/library/golang:1.23.1 AS gobuilder
WORKDIR /go/src/github.com/fivenet-app/fivenet/
COPY . ./
RUN apt-get update && \
    apt-get install -y git && \
    make build-go

# Final Image
FROM docker.io/library/alpine:3.20.3
WORKDIR /app
RUN apk --no-cache add ca-certificates tini tzdata && \
    mkdir -p ./.output/public
COPY --from=nodebuilder /app/.output/public ./.output/public
COPY --from=gobuilder /go/src/github.com/fivenet-app/fivenet/fivenet /usr/local/bin

EXPOSE 7070/tcp 8080/tcp

ENTRYPOINT ["tini", "--", "fivenet"]

CMD ["server"]
