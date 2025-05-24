# testdata/sql

This folder contains base data for tests related to querying, creating, updating or deleting data in MySQL.

## Dump Commands

**This assumes you are connecting against the local test database run through `docker` and you are in the [`internal/tests/testdata/sql/`](/internal/tests/testdata/sql/) folder.**

### 1. RBAC

```console
mysqldump -h 127.0.0.1 -u fivenet -pchangeme --no-create-info --skip-triggers --complete-insert fivenet fivenet_attrs fivenet_permissions fivenet_roles fivenet_role_attrs fivenet_role_permissions > base_000_rbac.sql
```

### 2. Documents

```console
mysqldump -h 127.0.0.1 -u fivenet -pchangeme --no-create-info --skip-triggers --complete-insert fivenet fivenet_documents fivenet_documents_categories fivenet_documents_comments fivenet_documents_job_access fivenet_documents_references fivenet_documents_relations fivenet_documents_templates fivenet_documents_templates_job_access fivenet_documents_citizen_access > base_010_documents.sql
```
