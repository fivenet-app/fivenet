# testdata

This folder contains base data for tests related to querying, creating, updating or deleting data in MySQL.

## Dump Commands

**This assumes you are connecting against the local test database run through `docker` and you are in the `testdata/` folder.**

### Documents

```console
mysqldump -h localhost -u root -psecret-pw-for-root-user --no-create-info --skip-triggers arpanet arpanet_documents arpanet_documents_categories arpanet_documents_comments arpanet_documents_job_access arpanet_documents_references arpanet_documents_relations arpanet_documents_templates arpanet_documents_user_access > base_010_documents.sql
```
