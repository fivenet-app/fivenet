# syntax=docker/dockerfile:1.26-labs

# Version helper, computed once for the whole image build.
FROM docker.io/library/alpine:3.24.1 AS version

WORKDIR /app

RUN apk --no-cache add git

COPY .git .git
COPY internal/scripts/get-version.sh ./internal/scripts/get-version.sh

RUN sh ./internal/scripts/get-version.sh > /version

# Livemap Tiles Layer for improved caching
FROM docker.io/library/alpine:3.24.1 AS livemaptiles

WORKDIR /app

COPY ./public/images/livemap/ ./public/images/livemap/

RUN find ./public/images/livemap/ \
        ! -path '*/tiles*' -and ! -path '*/overlays*' -and ! -path './public/images/livemap/' \
        -exec rm -rf {} +

# Iconify icon sets for backend server
FROM docker.io/library/alpine:3.24.1 AS iconsets

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

# Frontend Build
FROM docker.io/library/node:24.19.0-alpine3.24 AS frontendbuild

WORKDIR /app

COPY --from=version /version /version
COPY --exclude=public/images/livemap/ . ./

RUN apk add --no-cache git python3 make gcc g++ && \
    corepack enable && \
    corepack prepare pnpm@10.34.5 --activate && \
    version="$(cat /version)" && \
    COMMIT_REF="$version" pnpm install && \
    NODE_OPTIONS="--max-old-space-size=8192" \
        COMMIT_REF="$version" pnpm generate

# Backend Build
FROM docker.io/library/golang:1.27.0 AS backendbuild

WORKDIR /go/src/github.com/fivenet-app/fivenet/v2026/

COPY --from=version /version /version
COPY --exclude=public/images/livemap/ . ./

RUN apt-get update && \
    apt-get install -y git && \
    version="$(cat /version)" && \
    make build-go GIT_VERSION="$version"

# Final Image
FROM docker.io/library/alpine:3.24.1

WORKDIR /app

VOLUME ["/config", "/data"]

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
COPY --from=livemaptiles /app/public/images/livemap/ ./.output/public/images/livemap/
COPY --from=iconsets /app/icons/ ./icons/
COPY --from=frontendbuild /app/.output/public/ ./.output/public/
COPY --from=backendbuild /go/src/github.com/fivenet-app/fivenet/v2026/fivenet /usr/local/bin

USER 2000

EXPOSE 8080/tcp 7070/tcp

ENTRYPOINT ["tini", "--", "fivenet"]

CMD ["server"]
