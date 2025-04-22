#!/usr/bin/env bash

NODE_VERSION=$(node -p -e "require('./package.json').version")

if [ "$1" = "stdout" ]; then
    pnpm git-cliff -o - --unreleased --tag "$NODE_VERSION"
else
    pnpm git-cliff -o './CHANGELOG.md' --tag "$NODE_VERSION"
fi
