# tracker

`tracker` is the service that keeps every client up-to-date on who is on duty and where they are.
It relies on **two NATS Key-Value buckets** so look-ups and fan-out are both *O(1)*.

| Bucket | Stream name (backing) | Key format | Stored value (protobuf) | Purpose |
|--------|-----------------------|------------|-------------------------|---------|
| **USERLOC** | `user_locations` | `JOB.GRADE.USER_ID` | [`UserMarker`](../../proto/resources/livemap/user_marker.proto) – lat/lon, heading, timestamp, … | Large, chatty dataset used by the live map. Keys embed *job* and *grade* so JetStream can filter by ACL. |
| **UNITMAP** | `user_mappings` | `USER_ID` | [`UserUnitMapping`](../../proto/resources/centrum/user_unit.proto) – unit_id, job, grade, created_at | Small index that answers: “Which unit / job / grade does user 123 belong to?” |

All updates go through the KV API, so you still get atomic `Put`, `Delete`, history, and conditional updates.

---

## How it works

```
writers ──► KV/JetStream ──► fan-out deltas ──► clients
               ▲
               │
               │  (compressed roll-up every 30 s)
           snapshot daemon
```

### 1 — Snapshot daemon

* **Bootstraps** with an *ephemeral* consumer using `DeliverLastPerSubjectPolicy` on **both** buckets → exactly one message per key, regardless of churn.
* Maintains two in-memory maps:
  * `locs  map[key]*UserMarker`
  * `units map[user_id]*UserUnitMapping`
* Every *N* seconds it publishes a **roll-up** message (`Nats-Rollup: all`, `KV-Operation: ROLLUP`) on a special subject inside `user_locations`. A late-joining client downloads only this one compressed blob.

### 2 — Client session

1. The server clones the big snapshot and applies the user’s ACL:
   * *Job* must match.
   * Either `grade ≤ maxGrade` (up-to mode) **or** `grade ∈ explicit list` (fine-grained mode).
2. Sends that filtered slice as the **initial payload**.
3. Binds a **durable pull consumer** to `user_locations` with `FilterSubjects` derived from the same ACL (`job.grade.*` or `job.*`).
4. If the operator later changes permissions, the service calls `UpdateConsumer` with the new filter list – no reconnect, no replay storm.

### 3 — O(1) look-ups by user id

```go
// find current marker when only USER_ID is known
u   := units[USER_ID]                         // UNITMAP first
key := fmt.Sprintf("%s.%d.%s", u.Job, u.Grade, USER_ID)
m   := locs[key]                              // then USERLOC
```

Two hash-table hits in memory, zero network calls.

### 4 — Deletes, purges, roll-ups

| Writer action | Header seen | tracker effect |
|---------------|-------------|----------------|
| `kv.Delete()` | `KV-Operation: DEL`   | tombstone → remove entry, keep history |
| `kv.Purge()`  | `KV-Operation: PURGE` | hard delete → remove entry + history |
| roll-up helper | `KV-Operation: ROLLUP` + `Nats-Rollup: all` | replace entire snapshot |

Clients receive the same headers; their durable consumer stays consistent without polling.

---

## Quick reference

| Need | Where / how |
|------|-------------|
| **Full list** a user may see | Clone snapshot → ACL filter |
| **Live updates** | Durable consumer on `user_locations` (`DeliverNewPolicy`) with multi-filter subjects |
| **Lookup “where is user 123?”** | Read `UNITMAP[123]` → derive key → read `USERLOC[key]` |
| **Cold-start time** | O(# current markers) thanks to `DeliverLastPerSubjectPolicy` |
| **Late-joiner bandwidth** | One compressed roll-up (~few kB) |

This structure lets `tracker` scale linearly with concurrent sessions, keeps look-ups instant, and avoids replay storms even when thousands of location updates per second flow through the system.
