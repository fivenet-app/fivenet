---
title: Document Templates
---

**Note**: A template must render out to valid HTML.

Golang templating is used. In addition to base [Golang html/template functions](https://pkg.go.dev/html/template) [`sprig` template functions](https://masterminds.github.io/sprig/) are available for convience.

* The whole template needs to be wrapped in `<p>` and `</p>`.
* Use `<br>` for new lines.

## Snippets

### Access Creator Info

```gotemplate
{{ .activeChar.firstname }}, {{ .activeChar.lastname }}
```

### Get first Citizen

```gotemplate
{{- $citizen := first .usersList -}}
```

Example access citizen info:

```gotemplate
{{ $citizen.firstname }}, {{ $citizen.lastname }} ({{ $citizen.dateofbirth }})
```

### Current Date and Time

```gotemplate
{{ now | date "02.01.2006 15:04" }}
```
