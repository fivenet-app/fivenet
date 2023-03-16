# arpanet

## Features

- [ ] Authentication
    - [x] Separate "accounts" table that allows users to login to the network
    - [ ] Create in-game "register" command to set username and password
- [x] User Database - 1. Prio
    - [x] Search by
        - [x] Name
        - [x] Wanted State
    - [x] Display a single user's info
        - [x] Show a feed of activity of the user (e.g., documents created, documents mentioned in)
    - [x] Wanted aka "additional UserProps"
        - [x] Allow certain jobs to set a person as wanted
        - [x] Add toggle to display only wanted people
- [ ] Documents ("Akten")
    - [x] Each document is independent and has no direct parent or responses
        - [x] Users can leave Comments on documents
    - [x] Documents can reference each other ("document activity feed"), e.g., DOJ asks for blood test on a patient, LSMD responds by creating the patient blood test result document and references the DOJ response
    - [ ] Templates (e.g., Haftbefehl)
        - [ ] Editable templates by the leader of the jobs
    - [ ] Sharing
        - [ ] Sharing with the same job automatically
        - [x] Sharing with users/ citizens (e.g., Patientenbefund is shared with the Patient, the lawyer and the DOJ
        - [ ] People can request access by link
    - [x] Category System (no directories/ paths)
        - [x] ~~Sub-categories~~  - One level of categories that is sorted by names
    - [ ] Different Styles/ Types (e.g., Arbeitsunfähigkeitsschein, Polizeireport)
- [ ] Dispatch System
    - [ ] Livemap to display dispatches
    - [ ] Tools to coordinate dispatches
        - [ ] Manually by user input
        - [ ] Automatically
- [x] Livemap
    - [x] See your colleagues (for now using Copnet VPC Connector's data)
        - [x] Create table model for our own player location table
        - [ ] Write FiveM plugin that writes into our own location table
    - [x] Multiple different designs
    - [ ] Future: See other jobs' positions and/ or dispatches
- [ ] Employee Management
    - [ ] Create Notes and Warns for Employees ("Führungsregister")
    - [ ] Promote and Demote Employees
    - [ ] Fire employees
- [ ] Permissions System
    - [ ] Based on Job + Job Rank/ Grade

## Development

### Required Tools

* Golang 1.20
* `yarn`
* [`protoc`](https://grpc.io/docs/protoc-installation/)
* `protoc-gen-go`:
    * `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28`
    * `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2`
* `protoc-gen-js`: Run `yarn` (without any args)
* `protoc-gen-grpc-web`: Download and install the latest release from https://github.com/grpc/grpc-web/releases
* `protoc-gen-validate`: Download and install the latest release from https://github.com/bufbuild/protoc-gen-validate/releases
* `protoc-go-inject-tag`: Run `go install github.com/favadi/protoc-go-inject-tag@latest`

### Codium/ VSCode Users

Make sure to disable the builtin Typescript plugin.

> 1. In your project workspace, bring up the command palette with Ctrl + Shift + P (macOS: Cmd + Shift + P).
> 2. Type built and select "Extensions: Show Built-in Extensions".
> 3. Type typescript in the extension search box (do not remove @builtin prefix).
> 4. Click the little gear icon of "TypeScript and JavaScript Language Features", and select "Disable (Workspace)".
> 5. Reload the workspace. Takeover mode will be enabled when you open a Vue or TS file.

Copied from and for more information on "why you should do this", see: https://vuejs.org/guide/typescript/overview.html#volar-takeover-mode

### GRPC Web Debugging in your Browser

You must use this forked version: [Github jrapoport/grpc-web-devtools](https://github.com/jrapoport/grpc-web-devtools).

## Database

### Indexes for existing Tables

* `users`
    * Indexes
        * `firstname` and `lastname` Columns:
            * `CREATE FULLTEXT INDEX IF NOT EXISTS idx_users_firstname_lastname ON s4_fivem.users (firstname, lastname);`
        * `job` and `job_grade` Spalten:
            * `CREATE INDEX IF NOT EXISTS users_job_grade_IDX USING BTREE ON s4_fivem.users (job_grade, job);`

## Livemap

Based upon https://gist.github.com/NelsonMinar/6600524#file-maketiles-sh and CopNet/ MedicNet code.
