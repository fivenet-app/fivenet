# arpanet

## Features

- [ ] Authentication
    - [ ] Use Discord Login and an in-game "register" command of some sorts
- [ ] User Database - 1. Prio
    - [ ] Search by `firstname`, `lastname` and `job`
    - [ ] Display a single user's info
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
    - [ ] Livemap that shows the dispatches
    - [ ] Tools to coordinate dispatches
        - [ ] Manually by user input
        - [ ] Automatically
- [ ] Livemap
    - [ ] See your colleagues
- [ ] Employee Management
    - [ ] Warn Employees ("Führungsregister")
    - [ ] Promote and Demote Employees
    - [ ] Fire employees

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

### What data is currently missing from FiveM tables?

* `users`
    * More Indexes
        * `firstname` and `lastname` Columns:
            * `CREATE FULLTEXT INDEX IF NOT EXISTS idx_users_firstname_lastname ON s4_fivem.users (firstname, lastname);`
        * `job` and `job_grade` Spalten:
            * `CREATE INDEX IF NOT EXISTS users_job_grade_IDX USING BTREE ON s4_fivem.users (job_grade, job);`

## Livemap

Based upon https://gist.github.com/NelsonMinar/6600524#file-maketiles-sh and CopNet/ MedicNet code.
