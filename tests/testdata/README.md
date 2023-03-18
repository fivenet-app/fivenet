# testdata

This folder contains base data for tests related to querying, creating, updating or deleting data in MySQL.

## Dump Commands

**This assumes you are connecting against the local test database run through `docker` and you are in the `testdata/` folder.**

### Documents

```console
mysqldump -h 127.0.0.1 -u arpanet -pchangeme --no-create-info --skip-triggers arpanet arpanet_documents arpanet_documents_categories arpanet_documents_comments arpanet_documents_job_access arpanet_documents_references arpanet_documents_relations arpanet_documents_templates arpanet_documents_user_access > base_010_documents.sql
```

### RBAC

```console
mysqldump -h 127.0.0.1 -u arpanet -pchangeme --no-create-info --skip-triggers arpanet arpanet_permissions arpanet_role_permissions arpanet_roles arpanet_user_permissions arpanet_user_roles > base_000_rbac.sql
```
