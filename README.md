# arpanet

## Features

- [ ] Authentication
    - [ ] Use Discord Login and an in-game "register" command of some sorts
- [ ] User Database
    - [ ] Search by `firstname` and `lastname`
- [ ] Documents
    - [ ] Different Styles/ Types (e.g., Arbeitsunfähigkeitsschein, Polizeireport)
    - [ ] Sharing with the same job automatically
    - [ ] Sharing with the citizen affected (e.g., Patientenbefund is shared with the Patient, the lawyer and the DOJ)
    - [ ] People can request access by link
- [ ] Dispatch System
    - [ ] Livemap
- [ ] Job Management
    - [ ] Warn Employees ("Führungsregister")
    - [ ] Promote and Demote Employees
    - [ ] Fire employees

## Development

### Required Tools

* Golang 1.19
* `swag` - Generate Swagger docs.
    ```console
    go install github.com/swaggo/swag/cmd/swag@latest
    ```
* `yarn`

### What data is currently missing from FiveM tables?

* `users`
    * Weitere Indexes
        * `firstname` und `lastname` Spalten:
            * `CREATE FULLTEXT INDEX IF NOT EXISTS users_firstname_IDX ON s4_fivem.users (firstname, lastname);`
        * `job` und `job_grade` Spalten:
            * `CREATE INDEX IF NOT EXISTS users_job_grade_IDX USING BTREE ON s4_fivem.users (job_grade, job);`
    * (Optional) Blood type
    * Rename `last_seen` to `updated_at`
