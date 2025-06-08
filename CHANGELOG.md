# Changelog

All notable changes to this project will be documented in this file.

## [2025.6.0] - 2025-06-07

### 🚀 Features

- Add auto grading to qualifications
- Store discord user token for guild list and other functionality in

### 🐛 Bug Fixes

- Oauth2 tokens migration
- Image inserted as base64 when file upload is enabled and other

### ⚙️ Miscellaneous Tasks

- Improve github actions image buid by using arm64 runner for arm64

## [2025.5.5] - 2025-05-29

### 🚀 Features

- Add collab editor to wiki pages
- Add filestore migrations command and cleanup otlp usage

### 🐛 Bug Fixes

- Perms events issues by reducing amount of events
- Db migrations using "wrong" fivenet_users table when esxcompat

## [2025.5.4] - 2025-05-28

### 🐛 Bug Fixes

- Perms events issues by reducing amount of events

## [2025.5.3] - 2025-05-24

### ⚙️ Miscellaneous Tasks

- Fix release-it issue with git-cliff caused by fetch-depth in

## [2025.5.2] - 2025-05-24

### 🚀 Features

- Add basic `fivenet tools sync status` subcommand
- Use NuxtImg for lazy loading images
- Rework index page
- Add custom i18n logic for backend
- Citizen documents improvements
- Cleanup naming schema of apis, perms and more

### 🐛 Bug Fixes

- Add validation for negative value/valid value in job grade list
- Document view issues caused by content rendered nuxtimg component
- Sync api error handling improvements
- Issues with html nuxtimg tag and citizen documents list

### 📚 Documentation

- Add github issues templates

## [2025.5.1] - 2025-05-10

### 🐛 Bug Fixes

- Flatten role attributes returning all attributes and not just the
- Job grades not being correctly "merged" on update from mstlystcdata

## [2025.4.7] - 2025-04-22

### 🚀 Features

- Add "basic" oauth2 client system
- Improve audit logging
- Use vueuse time ago function
- Add first draft of penalty calculator
- Use vee-validate i18n module instead of yup
- Start work on traffic infraction points
- Use more go generics in utils pkg
- Change table pagination change previous/next buttons if disabled
- Continue centrum feature work
- Add licenses page
- Fix centrum components issues
- Update generated protoc code and nuxt to latest 3.7
- Use citizen info popover more
- Remove copied hasher code
- Show unit initials above player markers
- Rework details for centrum units and dispatches
- Improve timeclock daily pagination issues
- Show dates in dispatch list
- Add basic s3 storage logic
- Remove job roles and "employee" role
- Add attributes to units
- Add server time correction
- Add icon marker when creating new marker or dispatch
- Take unit users location into account when auto assigning
- Add no dispatch auto assign unit attribute
- Remove comma between first and last name
- Add predefined status for dispatches and units
- Remove sentry from back- and frontend
- Improve storage system to use filehash to deduplicate files in
- Replace absence_date with absence_begin and absence_end
- Add upload file function for "jobassets"
- Use custom pagination component
- Allow updating marker marker info
- Restructure citizen info/profile layout
- Rework citizen attributes to store the user/attributes in a table
- Cleanup protobuf generator logic
- Add basic notes field to colleague info and split the jobs user
- Add image preview to rector file list and adjust dispatch time
- Add color and icon fields to document category
- Clean up color select of tailwind css colors
- Improve category display and fix issue with qualification
- Remove timeclock handling cron
- Add thread recipient email name to prevent email changes to be
- Start work on internet feature
- Add document owner override hatch
- Add "add image" feature to tiptap editor via modal
- Add FIVENET_SKIP_DB_MIGRATIONS env to disable db migrations
- Based on esx compat mode the users/fivnet_users table is switched
- Remove nickname regex bracket logic in favor of using the user's
- Start work for some smaller features
- Continue work on email message attachments
- Joblist and jobgradelist are sent to the client and use custom
- Improve livemap marker stream performance by using selective
- Make sure a fallback access is added to wiki pages when none is

### 🐛 Bug Fixes

- Notifications causing ambiguous import
- Correct query for retrieving user documents
- Remove pwa update logic not needed
- Use bigint as we started having some issues
- Clipboard issue after first doc creation
- Docstore user documents listing
- Some buttons not being correctly disabled
- Move quick buttons more to bottom right
- Issues with discord bot
- Improve auth middleware a tad bit
- Add state to template
- Add missing delete dispatch function
- Try to improve responsiveness for units assigned and unit users
- Continue work on jobs conduct system
- Dbmanager migration issue
- Improve notification lang issue
- Job being shown where it shouldn't
- Use text-neutral instead of text-white for better theming in the
- Improve unit user id mapping
- Notificator not restarting
- Issues with centrum state
- Clone state when computing update
- Try to improve unassigned dispatches not correctly being set to be
- Add debug log to take dispatch entry to see why the timestamps
- Apply of job perms removing job grade perms completely
- Make sure dispatches are at least 5 seconds old before auto
- Add ticker to update marker markers for now
- Tweak picture modal button conditions and start work on predefined
- Move consume err handler logic into events package
- Resolve units translation issue and use state for quick buttons
- Tweak login and logout page behavior
- Helm chart releaser issuer
- Disable ssr and help pages for now
- "fix" splitpanes issue for now
- Iconify to use local api at /api/icons
- Change default gray color
- Downgrade jodit editor brokey
- `UInputMenu` and `USelectMenu`  recursive issue
- Replace NuxtLink with ULink
- Use a custom DashboardSidebarLinks component till I know why the
- Jobs colleauges id page titles
- Tweak token cookie times
- Livemap markers not updating their data
- Centrum manager dispatch nil panic
- Add user and password to docker-compose nats
- App config not triggering updates in own instance
- Date select popover mode for mobile
- Use wrapped date picker popover to fix touch issue
- Add grpcws ping packet
- Icon issue for good
- Tweak login issues for servers not having char1:... identifier
- Citizen info second back button and add audit log for delete
- Improve oauth2 connect experience
- Unit creation issue for good
- Display qualification content when exam is enabled and can take
- Attempt to fix a discord role with KeepIfJobDifferent set being
- Page view edit/delete perms not checked
- Qualification editor requirements causing request validation to
- Make sure to drop any messenger service perms in email rework
- Citizen attributes, colleague labels create/update/delete check
- Unread email threads count issue by sending dummy state in thread
- Inactive colleagues list actions not showing and continue internet
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

### 💼 Other

- Only display dispatches that aren't older than 20 minutes
- Continue work on audit log
- Add check token api
- Add last char field to account retrieval
- Update

### 📚 Documentation

- Add poeditor project link

### ⚙️ Miscellaneous Tasks

- Add commitlint check
- Tweak user rank access control
- Use `navigateTo` instead of `useRouter()` + `push()`
- Remove auto imported vue imports
- Rework templates api
- Update deps
- Reorganize vue components
- Cleanup prettier ignore config
- Run prettier-plugin-tailwindcss
- Use @raffaelesgarro/vue-use-sound as a replacement due to bug
- Remove nuxt eslint module
- Update js and go deps
- Update js deps
- Add missing grpcws localhost test certs
- Downgrade vue to workaround v-calendar issue
- Update go prot generated files
- Switch to calver via release-it

<!-- generated by git-cliff -->
