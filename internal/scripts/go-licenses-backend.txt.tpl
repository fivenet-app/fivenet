THE FOLLOWING SETS FORTH ATTRIBUTION NOTICES FOR THIRD PARTY SOFTWARE THAT MAY BE CONTAINED IN PORTIONS OF THE FIVENET PRODUCT.
{{ range . }}
-----

The following software may be included in this product: {{ .Name }} (Version: {{ .Version }}). This software contains the following license and notice below:

{{ .LicenseName }} License - License URL: {{ .LicenseURL }}

{{ .LicenseText -}}
{{ end }}
