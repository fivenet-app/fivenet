---
title: "Requirements"
---

## Must have

* MySQL or MariaDB database aka the "gameserver database", the database the FiveM server is using.
* NATS message queue server or cluster (prefered), with JetStream enabled.
* Storage space: Either local filesystem or S3 bucket storage.
* Leaflet tiles: Generated using `gdal2tiles-leaflet` or similar. The map source image is expected to be `16384x16384` in resolution.

## Optional

* For OpenTelemetry based tracing support, a Jaeger instance.

## Database

### `jobs` Table

* `name`
* `label`

### `job_grades` Table

* `job_name`
* `grade`
* `name`
* `label`

### `licenses` Table

* `type`
* `label`

### `owned_vehicles` Table

* `owner`
* `plate`
* `type`
* `model` (Optional)

### `user_licenses` Table

* `type`
* `owner`

### `users` Table

* `id` int(11) NOT NULL AUTO_INCREMENT
* `identifier` varchar(64) NOT NULL
* `group`
* `firstname`
* `lastname`
* `dateofbirth`
* `job`
* `job_grade`
* `sex`
* `height`
* `phone_number`
