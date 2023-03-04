# arpanet

## Features

- [ ] Authentication
    - [x] Separate "accounts" table that allows users to login to the network
    - [ ] Create in-game "register" command to set username and password
- [ ] User Database - 1. Prio
    - [x] Search by name (`firstname`, `lastname`)
    - [x] Display a single user's info
        - [ ] Show a feed of activity of the user (e.g., documents created, documents mentioned in)
    - [ ] Wanted aka "additional UserProps"
        - [ ] Allow certain jobs to set a person as wanted
        - [ ] Add list to display only wanted people
- [ ] Documents
    - [ ] Templates (e.g., Haftbefehl)
    - [ ] Sharing
        - [ ] Sharing with the same job automatically
        - [ ] Sharing with the citizen affected (e.g., Patientenbefund is shared with the Patient, the lawyer and the DOJ)
        - [ ] People can request access by link
    - [ ] Different Styles/ Types (e.g., Arbeitsunfähigkeitsschein, Polizeireport)
    - [ ] Category System (no directories/ paths)
        - [ ] Sub-categories
- [ ] Dispatch System
    - [ ] Livemap display the dispatches
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
    - [ ] Warn Employees ("Führungsregister")
    - [ ] Promote and Demote Employees
    - [ ] Fire employees
- [ ] Permissions System
    - [ ] Based on Job + Job Rank/ Grade

## Development

### Required Tools

* Golang 1.19
* `yarn`
* [`protoc`](https://grpc.io/docs/protoc-installation/)
* `protoc-gen-go`:
    * `go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28`
    * `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2`
* `protoc-gen-js`: Run `yarn` (without any args)
* `protoc-gen-grpc-web`: Download and install the latest release from https://github.com/grpc/grpc-web/releases
* `protoc-gen-validate`: Download and install the latest release from https://github.com/bufbuild/protoc-gen-validate/releases

### Codium/ VSCode Users

Make sure to disable the builtin Typescript plugin.

> 1. In your project workspace, bring up the command palette with Ctrl + Shift + P (macOS: Cmd + Shift + P).
> 2. Type built and select "Extensions: Show Built-in Extensions".
> 3. Type typescript in the extension search box (do not remove @builtin prefix).
> 4. Click the little gear icon of "TypeScript and JavaScript Language Features", and select "Disable (Workspace)".
> 5. Reload the workspace. Takeover mode will be enabled when you open a Vue or TS file.

Copied from and for more information on "why you should do this", see: https://vuejs.org/guide/typescript/overview.html#volar-takeover-mode

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
