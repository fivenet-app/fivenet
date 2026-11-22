# syntax=docker/dockerfile:1.29-labs

# Version helper, computed once for the whole image build.
FROM docker.io/library/alpine:3.24.4 AS version

WORKDIR /app

RUN apk --no-cache add git

COPY .git .git
COPY internal/scripts/get-version.sh ./internal/scripts/get-version.sh

RUN sh ./internal/scripts/get-version.sh > /version

# Livemap Tiles Layer for improved caching
FROM docker.io/library/alpine:3.24.4 AS livemaptiles

WORKDIR /app

COPY ./public/images/livemap/tiles/ ./public/images/livemap/tiles/
COPY ./public/images/livemap/overlays/ ./public/images/livemap/overlays/

# Iconify icon sets for backend server
FROM docker.io/library/alpine:3.24.4 AS iconsets

WORKDIR /app

# Load only the used icon sets (see pkg/server/icons/iconify.go).
# Update this ref from https://github.com/iconify/icon-sets/commits/master.
ARG ICONIFY_ICON_SETS_REF=d1fdcf54272b299075b4fbf21851147e8c3846f9
RUN apk add --no-cache curl && \
    mkdir -p icons && \
    for icon_set in mdi simple-icons flagpack; do \
        curl --fail --location --silent --show-error \
            "https://raw.githubusercontent.com/iconify/icon-sets/${ICONIFY_ICON_SETS_REF}/json/${icon_set}.json" \
            --output "icons/${icon_set}.json"; \
    done

# Frontend Build
FROM docker.io/library/node:26.13.1-alpine3.24 AS frontendbuild

WORKDIR /app

COPY --from=version /version /version
COPY --exclude=public/images/livemap/ . ./

RUN apk add --no-cache git python3 make gcc g++ pkgconfig && \
    npm install -g corepack && \
    corepack enable && \
    corepack install && \
    version="$(cat /version)" && \
    COMMIT_REF="$version" pnpm install --frozen-lockfile && \
    NODE_OPTIONS="--max-old-space-size=8192" \
        COMMIT_REF="$version" pnpm generate

# Backend Build
FROM docker.io/library/golang:1.27.4 AS backendbuild

WORKDIR /go/src/github.com/fivenet-app/fivenet/v2026/

COPY --from=version /version /version
COPY --exclude=public/images/livemap/ . ./

RUN apt-get update && \
    apt-get install -y git && \
    version="$(cat /version)" && \
    make build-go GIT_VERSION="$version"

# Final Image
FROM docker.io/library/alpine:3.24.4

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
