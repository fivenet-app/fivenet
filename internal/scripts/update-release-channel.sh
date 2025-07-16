#!/usr/bin/env bash
set -euo pipefail

# Usage: update-channel.sh <channel> <image-tag> [<chart-version>]

cd "$(dirname "$0")" || { echo "Failed to enter script dir."; exit 1; }
cd ../../ || { echo "Failed to enter project root."; exit 1; }

if [[ $# -lt 2 || $# -gt 3 ]]; then
  echo "Usage: $0 <channel> <image_tag> [<chart_version>]"
  echo "  <channel>       : 'stable' or 'dev'"
  echo "  <image_tag>     : the new appVersion / image tag"
  echo "  <chart_version> : (optional) if omitted, reuse existing .<channel>.chartVersion"
  exit 1
fi

CHANNEL="$1"
IMAGE_TAG="$2"
CHART_VER="${3:-}"

JSON_FILE="channels.json"
TMP_FILE="$(mktemp)"

# Fallback to the JSON entry for the channel if chart version isn't provided
if [[ -z "$CHART_VER" ]]; then
    CHART_VER=$(jq -r --arg ch "$CHANNEL" '.[$ch].chartVersion' "$JSON_FILE")
    if [[ "$CHART_VER" == "null" || -z "$CHART_VER" ]]; then
        echo "❌ Could not read .${CHANNEL}.chartVersion from $JSON_FILE"
        exit 1
    fi
fi

TIMESTAMP=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

# Update channel JSON entry
jq --arg ch "$CHANNEL" \
   --arg cv "$CHART_VER" \
   --arg av "$IMAGE_TAG" \
   --arg ts "$TIMESTAMP" \
   '
    # Ensure channel object exists
    if .[$ch] == null then
        error("Channel \"" + $ch + "\" not found in JSON")
    else
        .
    end

    # Only update if chartVersion or appVersion differ
    | if .[$ch].chartVersion == $cv and .[$ch].appVersion == $av then
        .
    else
        # Update channel fields
        .[$ch].chartVersion = $cv
        | .[$ch].appVersion   = $av
        | .[$ch].publishedAt  = $ts
    end
   ' "$JSON_FILE" > "$TMP_FILE"
mv "$TMP_FILE" "$JSON_FILE"

echo "✅ Updated '$CHANNEL' → chartVersion=$CHART_VER, appVersion=$IMAGE_TAG, publishedAt=$TIMESTAMP"
