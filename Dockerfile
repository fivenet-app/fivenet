# syntax=docker/dockerfile:1.18-labs

# Frontend Build
FROM docker.io/library/node:23.11.1-alpine3.22 AS nodebuilder

ARG NUXT_UI_PRO_LICENSE

WORKDIR /app

COPY --exclude=public/images/livemap/ . ./

RUN apk add --no-cache git && \
    corepack enable && \
    corepack prepare pnpm@10.15.1 --activate && \
    pnpm install && \
    NODE_OPTIONS="--max-old-space-size=8192" \
        NUXT_UI_PRO_LICENSE=${NUXT_UI_PRO_LICENSE} \
        pnpm generate

# Livemap Tiles Layer for improved caching
FROM docker.io/library/alpine:3.22.1 AS livemaptiles

WORKDIR /app

COPY ./public/images/livemap/ ./public/images/livemap/

RUN find ./public/images/livemap/ \
        ! -path '*/tiles*' -and ! -path './public/images/livemap/' \
        -exec rm -rf {} +

# Iconify icon sets for backend server
FROM docker.io/library/alpine:3.22.1 AS iconsets

WORKDIR /app

# Clone the icon sets repository and filter for JSON files only
RUN apk add --no-cache git && \
    mkdir -p icons && \
    git init icons && \
    cd icons && \
    git remote add -f origin https://github.com/iconify/icon-sets.git && \
    git config core.sparseCheckout true && \
    echo "json/" >> .git/info/sparse-checkout  && \
    git pull origin master && \
    rm -rf .git && \
    mv json/* . && \
    rm -rf json && \
    find . -type f ! -name '*.json' -delete

# Backend Build
FROM docker.io/library/golang:1.25.1 AS gobuilder

WORKDIR /go/src/github.com/fivenet-app/fivenet/v2025/

COPY --exclude=public/images/livemap/ . ./

RUN apt-get update && \
    apt-get install -y git && \
    make build-go

# Final Image
FROM docker.io/library/alpine:3.22.1

WORKDIR /app

VOLUME ["/config", "/data"]

COPY --from=livemaptiles /app/public/images/livemap/ ./.output/public/images/livemap/

## Install required packages and create a non-root user
RUN apk --no-cache add ca-certificates tini tzdata && \
    addgroup \
        --gid 2000 \
        fivenet && \
    adduser \
        --uid 2000 \
        --disabled-password \
        --gecos "" \
        --home "$(pwd)" \
        --ingroup fivenet \
        --no-create-home \
        fivenet && \
    mkdir -p ./.output/public

## Copy built files from the builder stages
COPY --from=iconsets /app/icons/ ./icons/
COPY --from=nodebuilder /app/.output/public ./.output/public
COPY --from=gobuilder /go/src/github.com/fivenet-app/fivenet/v2025/fivenet /usr/local/bin

USER 2000

EXPOSE 8080/tcp 7070/tcp

ENTRYPOINT ["tini", "--", "fivenet"]

CMD ["server"]
