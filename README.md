<p align="center">
    <img alt="FiveNet Logo" src="public/images/social-card.png" width="640" />
</p>

# FiveNet

[![Container Images on GHCR.io](https://img.shields.io/badge/Container%20Images%20on-GHCR.io-blue)](https://github.com/fivenet-app/fivenet/pkgs/container/fivenet) [![Helm Logo](https://img.shields.io/badge/Helm%20Chart%20-available?logo=Helm&labelColor=0F1689)](https://github.com/FiveNet-app/charts) [![Discord](https://img.shields.io/badge/Discord-%235865F2.svg?&logo=discord&logoColor=white)](https://discord.gg/sWvkHuVQA5) [![fivenet.app Website](https://img.shields.io/badge/Website-fivenet.app-purple)](https://fivenet.app)

[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](https://go.dev/) [![Nuxt.js](https://img.shields.io/badge/Nuxt.js-00DC82?logo=nuxtdotjs&logoColor=fff)](https://nuxt.com/) [![Nuxt UI Pro](https://img.shields.io/badge/Made%20with-Nuxt%20UI%20Pro-00DC82?logo=nuxt.js&labelColor=020420)](https://ui.nuxt.com/pro) [![TypeScript](https://img.shields.io/badge/TypeScript-3178C6?logo=typescript&logoColor=fff)](#) [![NATS.io](https://img.shields.io/badge/nats.io-gray.svg?logo=natsdotio)](https://nats.io/) [![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=fff)](https://www.docker.com/)

## Roadmap

Things on the roadmap may or may not be implemented/changed/removed without warning.
For the roadmap [click here](https://github.com/users/galexrt/projects/2/views/1).

## Features

Be sure to checkout the website for an uptodate list of features with screenshots [here](https://fivenet.app/).

<details>
  <summary>Show partial feature List</summary>

- [x] Authentication
    - [x] Separate "accounts" table that allows users to log in to the network
    - [x] Last Char lock, prevent users from switching to chars other than the last one active on the server
- [x] "Content Moderation" access for server admins
    - [x] Use a list of ESX user groups in the config
    - [x] Allow them to switch jobs on the fly to always the highest job rank
    - [x] Allow them to edit/ delete any user content
- [x] Livemap
    - [x] See your colleagues (for now using Copnet VPC Connector's data)
        - [x] Create a table model for our player location table
    - [x] Multiple different designs
    - [x] Display dispatches (from GKS phone for now)
    - [x] See other jobs' positions and/ or dispatches
    - [x] Animated Marker when they move
    - [x] Search markers
    - [x] Postal Search
- [x] Permissions System
    - [x] Based on Job + Job Rank/ Grade
- [x] User Database - 1. Prio
    - [x] Search by
        - [x] Name
        - [x] Wanted State
    - [x] Display a single user's info
        - [x] Show a feed of the activity of the user (e.g., documents created, documents mentioned in)
    - [x] Wanted aka "additional UserProps"
        - [x] Allow certain jobs to set a person as wanted
        - [x] Add toggle to display only wanted people
- [x] Vehicles Search
    - [x] By Plate
    - [x] By Citizen on the citizen profile
- [x] Documents ("Akten")
    - [x] Each document is independent and has no direct parent or responses
        - [x] Users can leave Comments on documents
    - [x] Documents can reference each other ("document activity feed"), e.g., DOJ asks for a blood test on a patient, LSMD responds by creating the patient blood test result document and references the DOJ response
    - [x] Templates
        - [x] Add requirements for templates
    - [x] Sharing
        - [x] Sharing with the same job automatically
        - [x] Sharing with users/ citizens (e.g., Patientenbefund is shared with the Patient, the lawyer and the DOJ)
    - [x] Category System (no directories/ paths)
        - [x] ~~Sub-categories~~  - One level of categories that are sorted by names
    - [x] Functionality
        - [x] Create Documents with access
        - [x] Edit Documents
            - [x] With access modifications
            - [x] Set/ Update document category
            - [x] Set Access for Jobs and Users
        - [x] Document Comments
            - [x] View Document Comments
            - [x] Post Document Comments
            - [x] Edit Document Comments
- [x] "Completor" Service
    - [x] Uses [Bleve search](https://blevesearch.com/)
- [x] Breadcrumbs
    - [x] Use the closest thing to a page title (e.g., when viewing a user or editing a document) to build the breadcrumbs
- [x] "Faction Leader Control Panel" aka "Rector Service"
    - [x] Permission Editor for the job ranks (Rector)
        - [x] Can view the permissions
        - [x] Can edit the permissions
    - [x] Templates (DocStore)
        - [x] Create templates
        - [x] Edit templates
    - [x] Category (DocStore)
        - [x] Create Categories
        - [x] Edit Categories
        - [x] Delete categories
- [x] FiveM Integration plugin
    - [x] Livemap - Player position tracker plugin

And more features, check out [Getting Started - FiveNet website](https://fivenet.app/getting-started).

</details>

## Installation

### Container Images

Container images are available via [GHCR.io fivenet-app/fivenet image](https://github.com/orgs/fivenet-app/packages/container/package/fivenet).

### Helm Chart

Helm chart is available in the separate [GitHub fivenet-app/charts repository](https://github.com/fivenet-app/charts).

### Plugins

Gameserver plugins are available in the separate [GitHub fivenet-app/plugins repository](https://github.com/fivenet-app/plugins).

## Development

Please see [development documentation](https://fivenet.app/development).

## Translations (i18n)

Want to help translate FiveNet or fix a translation issue? You can join the translation group by heading over to [FiveNet's POEditor.com project here](https://poeditor.com/join/project/hszo85uo3K).

Thanks to [POEditor](https://poeditor.com/) for providing FiveNet with a free translation project!

![POEditor Logo](https://poeditor.com/public/images/ui/logos/logo_dark.svg)

## Security

If you find a vulnerability or a potential vulnerability in FiveNet, please see the [security release process](SECURITY.md).

## Credits

* Leaflet Livemap Code CRS: Based upon [NelsonMinar's Map Viewer Gist](https://gist.github.com/NelsonMinar/6600524) and VPC's CopNet/ MedicNet livemap code, and a lot of Leaflet CRS related Stackoverflow posts.

## License

Code licensed under Apache 2.0 license, see [LICENSE](LICENSE).

Licenses of used libraries, code and media can be found in the [`public/licenses/` folder](public/licenses/).
