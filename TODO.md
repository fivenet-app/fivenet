- [ ] Leitstelle
    - [x] Player tracker in backend (`pkg/tracker/tracker.go`)
    - [ ] Centrum
        - [x] Convert GKSPhone Job Messages and delete them afterwards - started in `gen/go/proto/services/centrum/converter.go`
    - [ ] Frontend for employees to see their dispatches and manage their status
        - [ ] Notifications on Livemap when they are assigned a dispatch
        - [ ] Updating their own dispatch status
        - [ ] Updating their unit status
        - [ ] Sidebar for Livemap
    - [ ] Expire dispatch assignment after 15 seconds if not accepted
        * The instance that sent the assignment update should be the one to take care of checking if the dispatch assignment has been expired.
        * Loop over the dispatch assignments every 2 seconds.
    - [ ] TakeControl for leitstellen controller needed + "No one is active" when no controller is active
    - [x] Rector - Add Centrum Settings Page
    - [ ] Translations

***

- [x] TASK Work on Unit self assignment logic
- [ ] TASK Work on Unit disponent unit assignment logic
- [ ] TEST Make sure the unit updating, checking if non empty, etc., works
- [ ] DISCUSS What tasks and tests are needed for the management of dispatches in backend?
- [ ] DISCUSS What tasks and tests are needed for the management of dispatches in frontend?
