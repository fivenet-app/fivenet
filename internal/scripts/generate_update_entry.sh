#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd "$SCRIPT_DIR" || { echo "Failed to enter script directory."; exit 1; }

cd ../../ || { echo "Failed to enter root directory."; exit 1; }

OLD_VERSION="$(git describe --abbrev=0 --tags --exclude='fivenet-*' "$(git rev-list --tags --skip=1 --max-count=1 --exclude='fivenet-*')")"
VERSION="${VERSION:-$(cat VERSION)}"

COMMITS="$(git --no-pager log "${OLD_VERSION}".."${VERSION}" --pretty=format:'* %s' | sort | uniq)"

cat <<EOF
# FiveNet Update \`${VERSION}\`

TODO add updatesummary

Danke an alle die FiveNet mit Übersetzungen, Feedback, und detaillierten Bug Reports unterstützen!

## Highlights

* :fivenet_star: TODO_FILL_IN_HIGHLIGHT
* :fivenet_fixed: TODO_FILL_IN_HIGHLIGHT
* :fivenet_maintenance: TODO_FILL_IN_HIGHLIGHT
* :fivenet_chore: TODO_FILL_IN_HIGHLIGHT

|| @🔔-FiveNet-Notifications ||

EOF

echo "---"
read -r -p "Copied, filled out and posted above message to the announcements channel? "
echo "---"
echo

cat <<EOF
# FiveNet Changes \`${OLD_VERSION}\` -> \`${VERSION}\`

## Commits

### Features

$(echo "${COMMITS}" | grep "feat: " | sed 's/^feat: /:fivenet_star: /g')

### Fixes

$(echo "${COMMITS}" | grep "fix: " | sed 's/^fix: /:fivenet_fixed: /g')

### Chores

$(echo "${COMMITS}" | grep "docs: " | sed 's/^docs: /:fivenet_maintenance: /g')
$(echo "${COMMITS}" | grep "ci: " | sed 's/^ci: /:fivenet_maintenance: /g')
$(echo "${COMMITS}" | grep "build: " | sed 's/^build: /:fivenet_maintenance: /g')
$(echo "${COMMITS}" | grep "chore: " | sed 's/^chore: /:fivenet_chore: /g')
EOF
