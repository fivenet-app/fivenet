---
title: Templates
---

**Note**: A template must render out to valid HTML.

Golang templating is used. In addition to base [Golang html/template functions](https://pkg.go.dev/html/template) [`sprig` template functions](https://masterminds.github.io/sprig/) are available for convience.

* The whole template needs to be wrapped in `<p>` and `</p>`.
* Use `<br>` for new lines.

## Available Variables

* `.Documents` - List of documents that are in the user's clipboard.
    * `.Id`
    * `.CreatedAt`
    * `.Title`
    * `.State`
    * `.CreatorId`
    * `.Creator` - See [User Info Structure](#user-info-structure).
    * `.Closed` - Boolean.
    * `.CategoryId`
    * `.Category`
        * `.Name`
        * `.Description`
* `.Users` - List of citizens/ users that are in the user's clipboard.
    * See [User Info Structure](#user-info-structure).
* `.Vehicles` - List of vehicles that are in the user's clipboard.
    * `Plate`
    * `Model`
    * `Type`
    * `Owner` - See [User Info Structure](#user-info-structure).
* `.ActiveChar` - Author/Submitting user's info.
    * See [User Info Structure](#user-info-structure).

### User Info Structure

* `.UserId`
* `.Identifier`
* `.Job`* - Preferrably use `jobLabel`.
* `.JobLabel`*
* `.JobGrade`* - Preferrably use `jobGradeLabel`.
* `.JobGradeLabel`*
* `.Firstname`
* `.Lastname`
* `.Dateofbirth` - In `DD.MM.YYYY` format.
* `.PhoneNumber` - Optional, might not always be included.

(\*these fields are only available on the `.activeChar` variable)

## Snippets

### Access Creator Info

```gotemplate
{{ .ActiveChar.Firstname }}, {{ .ActiveChar.Lastname }}
```

### Get first Citizen

Get the first user in the list (first in the user's clipboard):

```gotemplate
{{- $citizen := first .Users -}}
```

Example access citizen info:

```gotemplate
{{ $citizen.Firstname }}, {{ $citizen.Lastname }} ({{ $citizen.Dateofbirth }})
```

### Current Date and Time

```gotemplate
{{ now | date "02.01.2006 15:04" }}
```

To learn more about different date and time formats, check out [the Golang `time` package documentation here](https://pkg.go.dev/time#pkg-constants).

### Showing a Timestamp (e.g., `CreatedAt`)

```gotemplate
{{ .CreatedAt | date "02.01.2006 15:04" }}
```

### Checkbox

```html
<label contenteditable="false"><input type="checkbox"><span> Yes </span></label>
```

## Examples

### Show List of Vehicles

```gotemplate
{{ if not .Vehicles }}
<p>
No Vehicles involved.
</p>
{{ else }}
<ul>
{{- range .Vehicles -}}
<li>{{ .Plate }} - {{ .Owner.Firstname }}, {{ .Owner.Lastname }}</li>
{{- end -}}
</ul>
{{ end }}
```
