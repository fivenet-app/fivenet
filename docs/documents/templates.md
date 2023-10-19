---
title: Document Templates
---

**Note**: A template must render out to valid HTML.

Golang templating is used. In addition to base [Golang html/template functions](https://pkg.go.dev/html/template) [`sprig` template functions](https://masterminds.github.io/sprig/) are available for convience.

* The whole template needs to be wrapped in `<p>` and `</p>`.
* Use `<br>` for new lines.

## Available Variables

* `.documents` - Documents that are in the user's clipboard.
    * `id`
    * `createdAt`
    * `title`
    * `state`
    * `creatorId`
    * `creator` - See [User Info Structure](#user-info-structure).
    * `closed` - Boolean.
    * `categoryId`
    * `category`
        * `name`
        * `description`
* `.users` - List of citizens/ users that are in the user's clipboard.
    * See [User Info Structure](#user-info-structure).
* `.vehicles` - Vehicles that are in the user's clipboard.
    * `plate`
    * `model`
    * `type`
    * `owner` - See [User Info Structure](#user-info-structure).
* `.activeChar` - Submitting user's info.
    * See [User Info Structure](#user-info-structure).

### User Info Structure

* `userId`
* `identifier`
* `job`* - Preferrably use `jobLabel`.
* `jobLabel`*
* `jobGrade`* - Preferrably use `jobGradeLabel`.
* `jobGradeLabel`*
* `firstname`
* `lastname`
* `dateofbirth` - In `DD.MM.YYYY` format.

(\*these fields are only available on the `.activeChar` variable)

## Snippets

### Access Creator Info

```gotemplate
{{ .activeChar.firstname }}, {{ .activeChar.lastname }}
```

### Get first Citizen

Get the first user in the list (first in the user's clipboard):

```gotemplate
{{- $citizen := first .users -}}
```

Example access citizen info:

```gotemplate
{{ $citizen.firstname }}, {{ $citizen.lastname }} ({{ $citizen.dateofbirth }})
```

### Current Date and Time

```gotemplate
{{ now | date "02.01.2006 15:04" }}
```

To learn more about different date and time formats, check out [the Golang `time` package documentation here](https://pkg.go.dev/time#pkg-constants).

## Examples

### Create List of Vehicles

```gotemplate
{{ if not .vehicles }}
<p>
No Vehicles involved.
</p>
{{ else }}
<ul>
{{- range .vehicles -}}
<li>{{ .plate }} - {{ .owner.firstname }}, {{ .owner.lastname }}</li>
{{- end -}}
</ul>
{{ end }}
```
