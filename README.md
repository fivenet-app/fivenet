# arpanet

## Features

- [ ] Authentication
    - [ ] Use Discord Login and an in-game "register" command of some sorts
- [ ] Citizen Database
    - [ ] Search
    - [ ] Auto-generate blood type for LSMD
- [ ] Documents
    - [ ] Different Styles/ Types (e.g., Arbeitsunfähigkeitsschein, Polizeireport)
    - [ ] Sharing with the same job automatically
    - [ ] Sharing with the citizen affected (e.g., Patientenbefund is shared with the Patient, the lawyer and the DOJ)
    - [ ] People can request access by link
- [ ] Job Management
    - [ ] Fire employees
    - [ ] Warn Employees
    - [ ] Promote and Demote Employees 

## Development

```console
go install github.com/swaggo/swag/cmd/swag@latest
```

### What data is currently missing from FiveM tables?

* `users`
    * ID (Auto increment ID + Index)
    * Blood type
    * Rename `last_seen` to `updated_at`
* `jobs`
    * ID (Auto increment ID + Index)
* `job_grades`
    * ID (Auto increment ID + Index)

### Flow

1. Load `jobs` and `job_grades` from the FiveM database.
