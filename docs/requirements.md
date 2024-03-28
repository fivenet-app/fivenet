---
title: "Requirements"
---

## Must have

* MySQL or MariaDB database aka the "gameserver database", the database the FiveM server is using.
    * Make sure your tables have at least the [below structure](#database).
* NATS message queue server or cluster (prefered), with JetStream and memory storage enabled (you probably also want to have at least 1-20MB of memory storage available).
* Storage space: Either local filesystem directory or S3 bucket storage.
* Leaflet Map tiles: Generated using `gdal2tiles-leaflet` or similar. The map source image is expected to be `16384x16384` in resolution.
    * To be able to generate the tiles, you must have the map file in the `./internal/maps/` directory. You can use `make gen-tiles` to generate the tiles.

## Optional

* Tracing: For OpenTelemetry based tracing support.
    * Currently only Jaeger is supported as an exporter target.

***

## Database

This is a list of expected tables and their columns:

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

### `user_licenses` Table

* `type`
* `owner` - `varchar(64) NOT NULL`

### `owned_vehicles` Table

* `owner` - `varchar(64) NOT NULL`
* `plate`
* `type`
* `model` (Optional, can be overriden via `database.custom.columns.user.visum`)

### `users` Table

* `id` - `int(11) NOT NULL AUTO_INCREMENT`
* `identifier` - `varchar(64) NOT NULL`
* `group`
* `firstname`
* `lastname`
* `dateofbirth`
* `job`
* `job_grade`
* `sex`
* `height`
* `phone_number`
* `visum` (Optional, can be overriden via `database.custom.columns.user.visum`)
* `playtime` (Optional, can be overriden via `database.custom.columns.user.playtime`)
