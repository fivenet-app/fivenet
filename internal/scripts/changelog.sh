#!/usr/bin/env bash

NODE_VERSION=$(node -p -e "require('./package.json').version")

if [ "$1" = "stdout" ]; then
    pnpx git-cliff -o - --unreleased --tag "$NODE_VERSION"
else
    pnpx git-cliff -o './CHANGELOG.md'
fi
