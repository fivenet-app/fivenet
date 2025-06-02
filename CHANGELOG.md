# Changelog

All notable changes to this project will be documented in this file.

## [2025.5.5] - 2025-05-29

### 🐛 Fixes

- Perms events issues by reducing amount of events

## [2025.5.3] - 2025-05-24

### ⚙️ Miscellaneous Tasks

- Fix release-it issue with git-cliff caused by fetch-depth in

## [2025.5.2] - 2025-05-24

### 🐛 Fixes

- Add validation for negative value/valid value in job grade list
- Document view issues caused by content rendered nuxtimg component
- Sync api error handling improvements
- Issues with html nuxtimg tag and citizen documents list

### 🚀 Features

- Add basic `fivenet tools sync status` subcommand
- Use NuxtImg for lazy loading images
- Rework index page
- Add custom i18n logic for backend
- Citizen documents improvements
- Cleanup naming schema of apis, perms and more

### 📚 Documentation

- Add github issues templates

## [2025.5.1] - 2025-05-10

### 🐛 Fixes

- Flatten role attributes returning all attributes and not just the
- Job grades not being correctly "merged" on update from mstlystcdata

## [0.9.5] - 2025-04-16

### 🐛 Fixes

- Use modified useFileSelection function for editor image modal
- Downgrade nuxt to 3.14.59 for now
- Doc comments causing multiple notifications if an user created
- Timeclock list issue caused by dynamic table columns causing
- Adjust start point for timeclock timeline point with start but no
- Timeclock floating accuracy issue
- Add userLocale getter to workaround app config and user locale
- Housekeeper not running as expected for servers with cron agent
- Workaround corepack/pnpm install issue
- Replace some strconv.Atoi with strconv.ParseInt for int32
- Improve discord name construction to take 32 char limit into
- Add stream cancellation logic for websocket to be smarter
- Use jsoniter fastest for audit log and change croner scheduler
- Improve disable superuser mode flow and disabled dispatch center
- Use different user to run fivenet

### 🚀 Features

- Add "add image" feature to tiptap editor via modal
- Add FIVENET_SKIP_DB_MIGRATIONS env to disable db migrations
- Based on esx compat mode the users/fivnet_users table is switched
- Remove nickname regex bracket logic in favor of using the user's
- Start work for some smaller features
- Continue work on email message attachments
- Joblist and jobgradelist are sent to the client and use custom
- Improve livemap marker stream performance by using selective
- Make sure a fallback access is added to wiki pages when none is

### ⚙️ Miscellaneous Tasks

- Update

## [0.9.4] - 2024-12-22

### 🐛 Fixes

- Display qualification content when exam is enabled and can take
- Attempt to fix a discord role with KeepIfJobDifferent set being
- Page view edit/delete perms not checked
- Qualification editor requirements causing request validation to
- Make sure to drop any messenger service perms in email rework
- Citizen attributes, colleague labels create/update/delete check
- Unread email threads count issue by sending dummy state in thread
- Inactive colleagues list actions not showing and continue internet

### 🚀 Features

- Improve category display and fix issue with qualification
- Remove timeclock handling cron
- Add thread recipient email name to prevent email changes to be
- Start work on internet feature
- Add document owner override hatch

## [0.9.3] - 2024-10-24

### 🐛 Fixes

- Citizen info second back button and add audit log for delete
- Improve oauth2 connect experience
- Unit creation issue for good

### 📚 Documentation

- Add poeditor project link

## [0.9.2] - 2024-10-09

### 🐛 Fixes

- Livemap markers not updating their data
- Centrum manager dispatch nil panic
- Add user and password to docker-compose nats
- App config not triggering updates in own instance
- Date select popover mode for mobile
- Use wrapped date picker popover to fix touch issue
- Add grpcws ping packet
- Icon issue for good
- Tweak login issues for servers not having char1:... identifier

### 🚀 Features

- Add image preview to rector file list and adjust dispatch time
- Add color and icon fields to document category
- Clean up color select of tailwind css colors

### ⚙️ Miscellaneous Tasks

- Add missing grpcws localhost test certs

## [0.9.1] - 2024-06-11

### 🐛 Fixes

- "fix" splitpanes issue for now
- Iconify to use local api at /api/icons
- Change default gray color
- Downgrade jodit editor brokey
- `UInputMenu` and `USelectMenu`  recursive issue
- Replace NuxtLink with ULink
- Use a custom DashboardSidebarLinks component till I know why the
- Jobs colleauges id page titles
- Tweak token cookie times

### 🚀 Features

- Allow updating marker marker info
- Restructure citizen info/profile layout
- Rework citizen attributes to store the user/attributes in a table
- Cleanup protobuf generator logic
- Add basic notes field to colleague info and split the jobs user

### ⚙️ Miscellaneous Tasks

- Add last char field to account retrieval

## [0.9.0] - 2024-04-10

### 🐛 Fixes

- Helm chart releaser issuer
- Disable ssr and help pages for now

## [0.8.21] - 2024-04-09

### 🐛 Fixes

- Tweak picture modal button conditions and start work on predefined
- Move consume err handler logic into events package
- Resolve units translation issue and use state for quick buttons
- Tweak login and logout page behavior

### 🚀 Features

- Add predefined status for dispatches and units
- Remove sentry from back- and frontend
- Improve storage system to use filehash to deduplicate files in
- Replace absence_date with absence_begin and absence_end
- Add upload file function for "jobassets"
- Use custom pagination component

## [0.8.20] - 2024-02-28

### 🐛 Fixes

- Add ticker to update marker markers for now

### 🚀 Features

- Remove comma between first and last name

## [0.8.19] - 2024-01-27

### 🐛 Fixes

- Apply of job perms removing job grade perms completely
- Make sure dispatches are at least 5 seconds old before auto

### 🚀 Features

- Take unit users location into account when auto assigning
- Add no dispatch auto assign unit attribute

## [0.8.16] - 2024-01-05

### 🐛 Fixes

- Try to improve unassigned dispatches not correctly being set to be
- Add debug log to take dispatch entry to see why the timestamps

### 🚀 Features

- Add server time correction
- Add icon marker when creating new marker or dispatch

## [0.8.14] - 2023-11-26

### 🐛 Fixes

- Clone state when computing update

### 🚀 Features

- Add attributes to units

## [0.8.13] - 2023-11-24

### 🐛 Fixes

- Notificator not restarting
- Issues with centrum state

## [0.8.11] - 2023-11-15

### 🚀 Features

- Show dates in dispatch list
- Add basic s3 storage logic
- Remove job roles and "employee" role

## [0.8.10] - 2023-10-31

### 🚀 Features

- Improve timeclock daily pagination issues

## [0.8.8] - 2023-10-16

### 🐛 Fixes

- Improve unit user id mapping

### 🚀 Features

- Rework details for centrum units and dispatches

## [0.8.7] - 2023-10-12

### 🐛 Fixes

- Use text-neutral instead of text-white for better theming in the

## [0.8.6] - 2023-10-06

### 🐛 Fixes

- Improve notification lang issue
- Job being shown where it shouldn't

### 🚀 Features

- Show unit initials above player markers

## [0.8.5] - 2023-09-28

### 🚀 Features

- Remove copied hasher code

## [0.8.4] - 2023-09-24

### 🐛 Fixes

- Dbmanager migration issue

## [0.8.1] - 2023-09-13

### 🐛 Fixes

- Improve auth middleware a tad bit
- Add state to template
- Add missing delete dispatch function
- Try to improve responsiveness for units assigned and unit users
- Continue work on jobs conduct system

### 🚀 Features

- Update generated protoc code and nuxt to latest 3.7
- Use citizen info popover more

## [0.8.0] - 2023-08-22

### 🐛 Fixes

- Issues with discord bot

### 🚀 Features

- Continue centrum feature work
- Add licenses page
- Fix centrum components issues

## [0.7.3] - 2023-06-30

### 🐛 Fixes

- Docstore user documents listing
- Some buttons not being correctly disabled
- Move quick buttons more to bottom right

### 🚀 Features

- Use more go generics in utils pkg
- Change table pagination change previous/next buttons if disabled

## [0.7.0] - 2023-06-18

### 🐛 Fixes

- Clipboard issue after first doc creation

## [0.6.5] - 2023-06-14

### 🚀 Features

- Start work on traffic infraction points

## [0.5.2] - 2023-05-31

### 🐛 Fixes

- Use bigint as we started having some issues

## [0.4.3] - 2023-05-22

### 🐛 Fixes

- Notifications causing ambiguous import
- Correct query for retrieving user documents
- Remove pwa update logic not needed

### 🚀 Features

- Add "basic" oauth2 client system
- Improve audit logging
- Use vueuse time ago function
- Add first draft of penalty calculator
- Use vee-validate i18n module instead of yup

### ⚙️ Miscellaneous Tasks

- Only display dispatches that aren't older than 20 minutes
- Add commitlint check
- Continue work on audit log
- Tweak user rank access control
- Add check token api

<!-- generated by git-cliff -->
