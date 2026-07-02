package main

import (
	"fmt"
	"maps"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/Masterminds/sprig/v3"
	permspb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/codegen/perms"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

// PermifyModule holds all state for this plugin.
type PermifyModule struct {
	*pgs.ModuleBase

	ctx      pgsgo.Context
	tpl      *template.Template
	constTpl *template.Template
	remapTpl *template.Template
}

// Permify returns an initialized PermifyModule.
func Permify() *PermifyModule { return &PermifyModule{ModuleBase: &pgs.ModuleBase{}} }

func (p *PermifyModule) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	p.ctx = pgsgo.InitContext(c.Parameters())

	serviceNameFn := func(s string) string {
		split := strings.Split(s, ".")
		return split[len(split)-1]
	}

	serviceNamespaceFn := func(s string) string {
		split := strings.Split(s, ".")
		return strings.Join(split[:len(split)-1], ".")
	}

	templateFns := sprig.FuncMap()
	templateFns["package"] = p.ctx.PackageName
	templateFns["name"] = p.ctx.Name
	templateFns["serviceName"] = serviceNameFn
	templateFns["serviceNamespace"] = serviceNamespaceFn
	templateFns["remapRef"] = remapRef
	templateFns["attrValueTypeName"] = attrValueTypeName
	templateFns["attrValueConstName"] = attrValueConstName

	tpl := template.New("permify").Funcs(templateFns)
	p.tpl = template.Must(tpl.Parse(permifyTpl))

	constTpl := template.New("permify_const").Funcs(templateFns)
	p.constTpl = template.Must(constTpl.Parse(permifyConstTpl))

	remapTpl := template.New("permify_remap").Funcs(templateFns)
	p.remapTpl = template.Must(remapTpl.Parse(permifyRemapTpl))
}

// Name satisfies the generator.Plugin interface.
func (p *PermifyModule) Name() string { return "permify" }

func qualifyService(namespace string, service string) string {
	if namespace == "" || service == "" || strings.Contains(service, ".") {
		return service
	}

	return namespace + "." + service
}

func shortServiceName(service string) string {
	split := strings.Split(service, ".")
	return split[len(split)-1]
}

func permPkgAlias(namespace string) string {
	replacer := strings.NewReplacer(".", "", "-", "_")
	return "perms" + replacer.Replace(namespace)
}

func remapImportPath(namespace string) string {
	return fmt.Sprintf(
		"github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/%s/perms",
		namespace,
	)
}

func remapTargetNamespace(defaultNamespace string, perm *Perm) string {
	if perm.Namespace != "" {
		return perm.Namespace
	}

	return defaultNamespace
}

func remapTargetService(defaultService string, perm *Perm) string {
	if perm.Service != nil && *perm.Service != "" {
		return *perm.Service
	}

	return defaultService
}

func remapRef(defaultNamespace string, defaultService string, perm *Perm) string {
	switch perm.Name {
	case "Any":
		return "perms.PermAnyRef"
	case "JobAdmin":
		return "perms.PermJobAdminRef"
	case "ConfigAdmin":
		return "perms.PermConfigAdminRef"

	default:
		namespace := remapTargetNamespace(defaultNamespace, perm)
		service := shortServiceName(remapTargetService(defaultService, perm))

		return fmt.Sprintf("%s.%s.%s.Perm", permPkgAlias(namespace), service, perm.Name)
	}
}

func sanitizeIdentifierPart(s string) string {
	var b strings.Builder
	upperNext := true

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			if b.Len() == 0 && unicode.IsDigit(r) {
				b.WriteRune('V')
			}

			if upperNext {
				b.WriteRune(unicode.ToUpper(r))
				upperNext = false
			} else {
				b.WriteRune(r)
			}
			continue
		}

		upperNext = true
	}

	if b.Len() == 0 {
		return "Value"
	}

	return b.String()
}

func attrValueTypeName(service string, perm *Perm, attr Attr) string {
	return sanitizeIdentifierPart(service) +
		sanitizeIdentifierPart(perm.Name) +
		sanitizeIdentifierPart(attr.Key) +
		"PermValue"
}

func attrValueConstName(service string, perm *Perm, attr Attr, value string) string {
	return sanitizeIdentifierPart(service) +
		sanitizeIdentifierPart(perm.Name) +
		sanitizeIdentifierPart(attr.Key) +
		"PermValue" +
		sanitizeIdentifierPart(value)
}

func mergeAttrsUnique(existing []Attr, extra []Attr) []Attr {
	if len(extra) == 0 {
		return existing
	}

	out := slices.Clone(existing)
	for _, candidate := range extra {
		dup := false
		for i := range out {
			if out[i].Key == candidate.Key && out[i].Type == candidate.Type {
				dup = true
				break
			}
		}
		if !dup {
			out = append(out, candidate)
		}
	}

	return out
}

type RemapImport struct {
	Alias string
	Path  string
}

func collectRemapImports(remaps map[string]map[string]map[string][]*Perm) []RemapImport {
	namespaces := map[string]struct{}{}
	for namespace, remapNs := range remaps {
		for _, remap := range remapNs {
			for _, targets := range remap {
				for _, target := range targets {
					if target.Name == "Any" || target.Name == "Superuser" ||
						target.Name == "JobAdmin" {
						continue
					}

					namespaces[remapTargetNamespace(namespace, target)] = struct{}{}
				}
			}
		}
	}

	imports := make([]RemapImport, 0, len(namespaces))
	for namespace := range namespaces {
		imports = append(imports, RemapImport{
			Alias: permPkgAlias(namespace),
			Path:  remapImportPath(namespace),
		})
	}
	slices.SortFunc(imports, func(a, b RemapImport) int {
		return strings.Compare(a.Path, b.Path)
	})

	return imports
}

func (p *PermifyModule) Execute(
	targets map[string]pgs.File,
	pkgs map[string]pgs.Package,
) []pgs.Artifact {
	visited := map[string][]pgs.File{}
	for _, t := range targets {
		key := t.File().InputPath().Dir().String()
		if _, ok := visited[key]; ok {
			visited[key] = append(visited[key], t)
			continue
		}

		visited[key] = []pgs.File{t}
	}

	remaps := map[string]map[string]map[string][]*Perm{}
	for _, fs := range visited {
		remap := p.generate(fs)
		maps.Copy(remaps, remap)
	}

	fs := []pgs.File{}
	for _, v := range visited {
		for _, f := range v {
			if !strings.Contains(f.File().InputPath().String(), "services/") {
				continue
			}
			fs = append(fs, f)
		}
	}

	slices.SortFunc(fs, func(a, b pgs.File) int {
		return strings.Compare(a.File().InputPath().String(), b.File().InputPath().String())
	})

	p.AddGeneratorTemplateFile(
		path.Join("perms_remap.go"),
		p.remapTpl,
		struct {
			FS              []pgs.File
			PermissionRemap map[string]map[string]map[string][]*Perm
			RemapImports    []RemapImport
		}{
			FS:              fs,
			PermissionRemap: remaps,
			RemapImports:    collectRemapImports(remaps),
		},
	)

	return p.Artifacts()
}

func (p *PermifyModule) generate(fs []pgs.File) map[string]map[string]map[string][]*Perm {
	f := fs[0]

	fqn := strings.Split(f.FullyQualifiedName(), ".")
	data := struct {
		FS                    []pgs.File
		F                     pgs.File
		GoPath                string
		PermissionServiceKeys []string
		Permissions           map[string]map[string]map[string]*Perm
		PermissionRemap       map[string]map[string]map[string][]*Perm
		Attributes            map[string]map[string]*Attr
	}{
		FS: fs,
		F:  f,
		GoPath: fmt.Sprintf(
			"github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/%s/perms",
			fqn[len(fqn)-1],
		),
		PermissionServiceKeys: []string{},
		Permissions:           map[string]map[string]map[string]*Perm{},
		PermissionRemap:       map[string]map[string]map[string][]*Perm{},
		Attributes:            map[string]map[string]*Attr{},
	}

	slices.SortFunc(fs, func(a, b pgs.File) int {
		return strings.Compare(a.File().InputPath().String(), b.File().InputPath().String())
	})

	for _, f := range fs {
		for _, s := range f.Services() {
			var order int32
			var icon string

			var serviceOpts permspb.ServiceOptions
			_, err := s.Extension(permspb.E_PermsSvc, &serviceOpts)
			if err != nil {
				p.Fail("error reading perms option:", err)
			}

			// Get full service name and remove `.services.` prefix.
			actualSvcKey := strings.TrimPrefix(string(s.FullyQualifiedName()), ".services.")
			actualSvcNameSplit := strings.Split(actualSvcKey, ".")
			actualNamespace := strings.Join(actualSvcNameSplit[:len(actualSvcNameSplit)-1], ".")
			actualSvcName := actualSvcNameSplit[len(actualSvcNameSplit)-1]

			permNamespace := actualNamespace
			if serviceOpts.GetNamespace() != "" {
				permNamespace = *serviceOpts.Namespace
			}

			permSvcKey := actualSvcKey
			permSvcName := actualSvcName
			if serviceOpts.GetService() != "" {
				permSvcKey = qualifyService(permNamespace, *serviceOpts.Service)
				permSvcName = shortServiceName(permSvcKey)
			}
			if !slices.Contains(data.PermissionServiceKeys, permSvcKey) {
				data.PermissionServiceKeys = append(data.PermissionServiceKeys, permSvcKey)
			}

			p.Debugf("Service: %s", actualSvcName)

			if serviceOpts.GetIcon() != "" {
				icon = serviceOpts.GetIcon()
			}
			if serviceOpts.GetOrder() != 0 {
				order = serviceOpts.GetOrder()
			}

			if len(serviceOpts.AdditionalPerms) > 0 {
				for _, v := range serviceOpts.AdditionalPerms {
					if _, ok := data.Permissions[permNamespace]; !ok {
						data.Permissions[permNamespace] = map[string]map[string]*Perm{}
					}
					if _, ok := data.Permissions[permNamespace][permSvcName]; !ok {
						data.Permissions[permNamespace][permSvcName] = map[string]*Perm{}
					}
					if v.Order >= 0 {
						v.Order = order
					}
					perm := &Perm{
						Namespace: permNamespace,
						Name:      v.Name,
						Service:   &permSvcName,
						Order:     order,
					}
					if icon != "" {
						perm.Icon = new(icon)
					}
					perm.Order *= 100

					perm.Attrs = make([]Attr, len(v.Attrs))
					for i, a := range v.Attrs {
						atype := "StringList"
						switch a.Type {
						case permissionsattributes.AttributeType_ATTRIBUTE_TYPE_JOB_LIST:
							atype = "JobList"
						case permissionsattributes.AttributeType_ATTRIBUTE_TYPE_JOB_GRADE_LIST:
							atype = "JobGradeList"
						}

						perm.Attrs[i] = Attr{
							Key:         a.Key,
							Type:        atype,
							ValidValues: slices.Clone(a.ValidStringList),
						}
						if a.ValidStringList != nil {
							perm.Attrs[i].Valid += "[]string{"
							for _, v := range a.ValidStringList {
								perm.Attrs[i].Valid += strconv.Quote(v) + ", "
							}
							perm.Attrs[i].Valid += "}"
						}
					}
					data.Permissions[permNamespace][permSvcName][v.Name] = perm
				}
			}

			for _, m := range s.Methods() {
				methodName := string(m.Name())
				methodName = strings.TrimPrefix(methodName, "services.")

				// Check if the field option is present and true
				var val permspb.PermsOptions
				ok, err := m.Extension(permspb.E_Perms, &val)
				if err != nil {
					p.Fail("error reading perms option:", err)
				}

				var perm *Perm
				if !ok {
					continue
				}

				if !val.Enabled {
					p.Fail("perms option not enabled for method:", actualSvcName, methodName)
					continue
				}

				perm = &Perm{
					Namespace: permNamespace,
					Name:      methodName,
				}
				names := []string{}
				if len(val.Names) > 0 {
					names = val.Names
					perm.Name = val.Names[0]
				} else if val.Name != nil && *val.Name != "" {
					names = append(names, *val.Name)
					perm.Name = *val.Name
				}
				if len(names) == 0 {
					names = append(names, perm.Name)
				}
				if val.GetNamespace() != "" {
					perm.Namespace = val.GetNamespace()
				}
				if val.GetService() != "" {
					perm.Service = val.Service
				}
				if val.GetOrder() != 0 {
					perm.Order = val.Order
				} else {
					perm.Order = order
				}
				perm.Order *= 100
				if icon != "" {
					perm.Icon = new(icon)
				}

				perm.Attrs = make([]Attr, len(val.Attrs))
				for i, a := range val.Attrs {
					atype := "StringList"
					switch a.Type {
					case permissionsattributes.AttributeType_ATTRIBUTE_TYPE_JOB_LIST:
						atype = "JobList"
					case permissionsattributes.AttributeType_ATTRIBUTE_TYPE_JOB_GRADE_LIST:
						atype = "JobGradeList"
					}

					perm.Attrs[i] = Attr{
						Key:         a.Key,
						Type:        atype,
						ValidValues: slices.Clone(a.ValidStringList),
					}
					if a.ValidStringList != nil {
						perm.Attrs[i].Valid += "[]string{"
						for _, v := range a.ValidStringList {
							perm.Attrs[i].Valid += strconv.Quote(v) + ", "
						}
						perm.Attrs[i].Valid += "}"
					}
				}

				serviceTargetOverridden := permSvcKey != actualSvcKey
				methodTargetOverridden := perm.Namespace != permNamespace || perm.Service != nil
				for _, name := range names {
					if name != methodName || len(names) > 1 ||
						serviceTargetOverridden || methodTargetOverridden {
						p.Debugf("names %+v; name: %s, mName: %s", names, name, methodName)

						var service *string
						if perm.Service != nil {
							s := *perm.Service
							service = &s
						} else {
							service = &permSvcName
						}

						pm := &Perm{
							Namespace: perm.Namespace,
							Service:   service,
							Name:      name,
							Attrs:     perm.Attrs,
							Order:     perm.Order,
							Icon:      perm.Icon,
						}
						if _, ok := data.PermissionRemap[actualNamespace]; !ok {
							data.PermissionRemap[actualNamespace] = map[string]map[string][]*Perm{}
						}
						if _, ok := data.PermissionRemap[actualNamespace][actualSvcName]; !ok {
							data.PermissionRemap[actualNamespace][actualSvcName] = map[string][]*Perm{}
						}

						data.PermissionRemap[actualNamespace][actualSvcName][methodName] = append(
							data.PermissionRemap[actualNamespace][actualSvcName][methodName],
							pm,
						)

						p.Debugf(
							"Permission Remap added: %q -> %q/%q\n",
							methodName,
							*pm.Service,
							pm.Name,
						)
					}
				}

				if slices.Contains(names, "JobAdmin") || slices.Contains(names, "ConfigAdmin") ||
					slices.Contains(names, "Any") ||
					methodTargetOverridden {
					continue
				}

				if _, ok := data.Permissions[permNamespace]; !ok {
					data.Permissions[permNamespace] = map[string]map[string]*Perm{}
				}
				if _, ok := data.Permissions[permNamespace][permSvcName]; !ok {
					data.Permissions[permNamespace][permSvcName] = map[string]*Perm{}
				}
				if _, ok := data.Permissions[permNamespace][permSvcName][perm.Name]; !ok {
					data.Permissions[permNamespace][permSvcName][perm.Name] = perm
					p.Debugf("Permission added: %q - %+v\n", methodName, perm)
				} else {
					p.Debugf("Permission already in list, updating: %q - %+v\n", methodName, perm)
					if len(perm.Attrs) > 0 {
						data.Permissions[permNamespace][permSvcName][perm.Name].Attrs = mergeAttrsUnique(
							data.Permissions[permNamespace][permSvcName][perm.Name].Attrs,
							perm.Attrs,
						)
					}
					perm.Order = data.Permissions[permNamespace][permSvcName][perm.Name].Order
				}
			}
		}
	}

	if len(data.Permissions) == 0 && len(data.PermissionRemap) == 0 {
		return nil
	}

	slices.Sort(data.PermissionServiceKeys)

	name := p.ctx.OutputPath(f)
	p.AddGeneratorTemplateFile(
		path.Join(filepath.Dir(name.String()), "service_perms.go"),
		p.tpl,
		data,
	)

	constPath := path.Join(filepath.Dir(name.String()), "perms", "perms.go")
	p.AddGeneratorTemplateFile(constPath, p.constTpl, data)

	return data.PermissionRemap
}

const permifyTpl = `// Code generated by protoc-gen-backend. DO NOT EDIT.
{{- range $f := .FS }}
// source: {{ $f.File.InputPath }}
{{- end }}

package {{ package .F }}

import (
    "github.com/fivenet-app/fivenet/v2026/pkg/config"
    "github.com/fivenet-app/fivenet/v2026/pkg/perms"
    permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
    permkeys "{{ .GoPath }}"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
	{{- range $namespace, $services := .Permissions }}
        // Namespace: {{ $namespace }}

        {{- range $svcName, $service := $services }}

		// Service: {{ $namespace }}.{{ $svcName }}
		    {{ range $perm := $service -}}
		{
			Namespace: permkeys.Namespace,
            Service: permkeys.{{ serviceName $svcName }}Perm,
			Name: permkeys.{{ serviceName $svcName }}{{ $perm.Name }}Perm,
            Attrs: []perms.Attr{
            {{- range $attr := $perm.Attrs }}
                {
                    Key: permkeys.{{ serviceName $svcName }}{{ $perm.Name }}{{ $attr.Key }}PermField,
                    Type: permissionsattributes.{{ $attr.Type }}AttributeType,
                    {{ with $attr.Valid -}}ValidValues: {{ $attr.Valid }},{{ end }}
                },
            {{- end }}
            },
            Order: {{ $perm.Order }},
            {{ with $perm.Icon }}Icon: "{{ $perm.Icon }}",{{ end }}
		},
            {{ end }}
        {{- end }}
    {{- end }}
	})
}
`

const permifyConstTpl = `// Code generated by protoc-gen-backend. DO NOT EDIT.
{{- range $f := .FS }}
// source: {{ $f.File.InputPath }}
{{- end }}

package perms{{ package .F }}

import (
    "github.com/fivenet-app/fivenet/v2026/pkg/perms"
)

const (
{{ with .PermissionServiceKeys }}
{{- $namespace := serviceNamespace (first .) }}
    Namespace perms.Namespace = "{{ $namespace }}"

{{ range $key, $svcName := . -}}
    {{ serviceName $svcName }}Perm perms.Service = "{{ serviceName $svcName }}"
{{ end }}
{{ end }}

{{ range $namespace, $services := $.Permissions }}
    {{ range $svcName, $service := $services }}
    // Service: {{ $namespace }}.{{ $svcName }}
	    {{- range $perm := $service }}
	{{ serviceName $svcName }}{{ $perm.Name }}Perm perms.Name = "{{ $perm.Name }}"
            {{- range $attr := $perm.Attrs }}
    {{ serviceName $svcName }}{{ $perm.Name }}{{ $attr.Key }}PermField perms.Key = "{{ $attr.Key }}"
            {{- end }}
		{{- end }}
	{{ end }}
{{ end }}
)

{{ range $namespace, $services := $.Permissions }}
    {{- range $svcName, $service := $services }}
        {{- range $perm := $service }}
            {{- range $attr := $perm.Attrs }}
                {{- if eq $attr.Type "StringList" }}
type {{ attrValueTypeName $svcName $perm $attr }} string

                    {{- if gt (len $attr.ValidValues) 0 }}
const (
                        {{- range $val := $attr.ValidValues }}
    {{ attrValueConstName $svcName $perm $attr $val }} {{ attrValueTypeName $svcName $perm $attr }} = {{ printf "%q" $val }}
                        {{- end }}
)
                    {{- end }}

                {{- end }}
            {{- end }}
        {{- end }}
    {{- end }}
{{ end }}

{{ range $namespace, $services := $.Permissions }}
    {{- range $svcName, $service := $services }}
type {{ serviceName $svcName }}Perms struct {
        {{- range $perm := $service }}
    {{ $perm.Name }} {{ serviceName $svcName }}{{ $perm.Name }}PermRef
        {{- end }}
}

        {{- range $perm := $service }}
type {{ serviceName $svcName }}{{ $perm.Name }}PermRef struct {
    Perm perms.PermissionRef
            {{- range $attr := $perm.Attrs }}
    {{ $attr.Key }} perms.AttrRef[perms.{{ $attr.Type }}Attr]
                {{- if eq $attr.Type "StringList" }}
    {{ $attr.Key }}Typed perms.StringListAttrRef[{{ attrValueTypeName $svcName $perm $attr }}]
                {{- end }}
            {{- end }}
}

        {{- end }}
var {{ serviceName $svcName }} = {{ serviceName $svcName }}Perms{
        {{- range $perm := $service }}
    {{ $perm.Name }}: {{ serviceName $svcName }}{{ $perm.Name }}PermRef{
        Perm: perms.NewPermissionRef(Namespace, {{ serviceName $svcName }}Perm, {{ serviceName $svcName }}{{ $perm.Name }}Perm),
            {{- range $attr := $perm.Attrs }}
        {{ $attr.Key }}: perms.New{{ $attr.Type }}AttrRef(
            perms.NewPermissionRef(Namespace, {{ serviceName $svcName }}Perm, {{ serviceName $svcName }}{{ $perm.Name }}Perm),
            {{ serviceName $svcName }}{{ $perm.Name }}{{ $attr.Key }}PermField,
        ),
                {{- if eq $attr.Type "StringList" }}
        {{ $attr.Key }}Typed: perms.NewTypedStringListAttrRef[{{ attrValueTypeName $svcName $perm $attr }}](
            perms.NewPermissionRef(Namespace, {{ serviceName $svcName }}Perm, {{ serviceName $svcName }}{{ $perm.Name }}Perm),
            {{ serviceName $svcName }}{{ $perm.Name }}{{ $attr.Key }}PermField,
        ),
                {{- end }}
            {{- end }}
    },
        {{- end }}
}

    {{- end }}
{{ end }}
`

const permifyRemapTpl = `// Code generated by protoc-gen-backend. DO NOT EDIT.
{{- range $f := .FS }}
// source: {{ $f.File.InputPath }}
{{- end }}

package goproto

import (
{{- range $import := .RemapImports }}
    {{ $import.Alias }} "{{ $import.Path }}"
{{- end }}
    perms "github.com/fivenet-app/fivenet/v2026/pkg/perms"
)

var PermsRemap = map[string][]perms.PermissionRef{
{{- range $namespace, $remapNs := .PermissionRemap }}
{{- range $service, $remap := $remapNs }}
	// Service: {{ $namespace }}.{{ $service }}
	    {{ range $key, $target := $remap -}}
	"{{ $namespace }}.{{ $service }}/{{ $key }}": {
            {{ range $t := $target -}}
        {{ remapRef $namespace $service $t }},
            {{- end }}
    },
        {{ end }}
    {{ end }}
{{ end }}
}
`

type Perm struct {
	Namespace string
	Service   *string
	Name      string
	Attrs     []Attr
	Order     int32
	Icon      *string
}

type Attr struct {
	Key         string
	Type        string
	Valid       string
	ValidValues []string
}
