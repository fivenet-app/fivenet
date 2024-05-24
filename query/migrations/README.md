# migrations

Uses [golang-migrate/migrate](https://github.com/golang-migrate/migrate) to run migrations if needed.

Migration file naming:

```console
# Update
{timestamp}_{title}.up.sql
# Rollback
{timestamp}_{title}.down.sql
```

(Timestamp is generated using `date +%s` command)

For more information about the migration library, see https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md

## Create ready-to-edit Migration files

Tested with GNU Bash (version `5.1.16`):

```console
REASON="yourreasonhere"
FILENAME_PREFIX="$(date +%s)_$REASON"
touch "${FILENAME_PREFIX}."{up,down}.sql
```
