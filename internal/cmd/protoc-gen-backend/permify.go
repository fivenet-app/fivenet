package main

import (
	"fmt"
	"maps"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"text/template"

	permspb "github.com/fivenet-app/fivenet/v2025/gen/go/proto/codegen/perms"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
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
		_, after, _ := strings.Cut(s, ".")
		return after
	}

	funcs := map[string]any{
		"package":     p.ctx.PackageName,
		"name":        p.ctx.Name,
		"serviceName": serviceNameFn,
	}

	tpl := template.New("permify").Funcs(funcs)
	p.tpl = template.Must(tpl.Parse(permifyTpl))

	constTpl := template.New("permify_const").Funcs(funcs)
	p.constTpl = template.Must(constTpl.Parse(permifyConstTpl))

	remapTpl := template.New("permify_remap").Funcs(funcs)
	p.remapTpl = template.Must(remapTpl.Parse(permifyRemapTpl))
}

// Name satisfies the generator.Plugin interface.
func (p *PermifyModule) Name() string { return "permify" }

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

	remaps := map[string]map[string][]*Perm{}
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

	p.AddGeneratorTemplateFile(
		path.Join("perms_remap.go"),
		p.remapTpl,
		struct {
			FS              []pgs.File
			PermissionRemap map[string]map[string][]*Perm
		}{
			FS:              fs,
			PermissionRemap: remaps,
		},
	)

	return p.Artifacts()
}

func (p *PermifyModule) generate(fs []pgs.File) map[string]map[string][]*Perm {
	f := fs[0]

	fqn := strings.Split(f.FullyQualifiedName(), ".")
	data := struct {
		FS                    []pgs.File
		F                     pgs.File
		GoPath                string
		PermissionServiceKeys []string
		Permissions           map[string]map[string]*Perm
		PermissionRemap       map[string]map[string][]*Perm
		Attributes            map[string]map[string]*Attr
	}{
		FS: fs,
		F:  f,
		GoPath: fmt.Sprintf(
			"github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/%s/perms",
			fqn[len(fqn)-1],
		),
		PermissionServiceKeys: []string{},
		Permissions:           map[string]map[string]*Perm{},
		PermissionRemap:       map[string]map[string][]*Perm{},
	}

	slices.SortFunc(fs, func(a, b pgs.File) int {
		return strings.Compare(a.File().InputPath().String(), b.File().InputPath().String())
	})

	for _, f := range fs {
		for _, s := range f.Services() {
			sName := strings.TrimPrefix(string(s.FullyQualifiedName()), ".services.")

			data.PermissionServiceKeys = append(data.PermissionServiceKeys, sName)
			p.Debugf("Service: %s (%s)", sName, data.PermissionServiceKeys)

			var order int32
			var icon *string

			var serviceOpts permspb.ServiceOptions
			_, err := s.Extension(permspb.E_PermsSvc, &serviceOpts)
			if err != nil {
				p.Fail("error reading perms option:", err)
			}

			if serviceOpts.Icon != nil {
				icon = serviceOpts.Icon
			}
			if serviceOpts.Order != 0 {
				order = serviceOpts.Order
			}

			for _, m := range s.Methods() {
				mName := string(m.Name())
				mName = strings.TrimPrefix(mName, "services.")

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
					p.Fail("perms option not enabled for method:", sName, mName)
					continue
				}

				perm = &Perm{
					Name: mName,
				}
				names := []string{}
				if len(val.Names) > 0 {
					names = val.Names
				} else if val.Name != nil && *val.Name != "" {
					names = append(names, *val.Name)
				}
				if val.Service != nil && *val.Service != "" {
					perm.Service = val.Service
				}
				if val.Order != 0 {
					perm.Order = val.Order
				} else {
					perm.Order = order
				}
				perm.Order *= 100
				perm.Icon = icon

				perm.Attrs = make([]Attr, len(val.Attrs))
				for i, a := range val.Attrs {
					atype := "StringList"
					switch a.Type {
					case permissions.AttributeType_ATTRIBUTE_TYPE_JOB_LIST:
						atype = "JobList"
					case permissions.AttributeType_ATTRIBUTE_TYPE_JOB_GRADE_LIST:
						atype = "JobGradeList"
					}

					perm.Attrs[i] = Attr{
						Key:  a.Key,
						Type: atype,
					}
					if a.ValidStringList != nil {
						perm.Attrs[i].Valid += "[]string{"
						for _, v := range a.ValidStringList {
							perm.Attrs[i].Valid += fmt.Sprintf("%q, ", v)
						}
						perm.Attrs[i].Valid += "}"
					}
				}

				for _, name := range names {
					if name != mName || len(names) > 1 {
						p.Debugf("names %+v; name: %s, mName: %s", names, name, mName)
						remapServiceName := strings.TrimPrefix(
							string(s.FullyQualifiedName()),
							".services.",
						)

						pm := &Perm{
							Service: perm.Service,
							Name:    name,
							Attrs:   perm.Attrs,
							Order:   perm.Order,
							Icon:    perm.Icon,
						}
						if _, ok := data.PermissionRemap[remapServiceName]; !ok {
							data.PermissionRemap[remapServiceName] = map[string][]*Perm{}
						}

						data.PermissionRemap[remapServiceName][mName] = append(
							data.PermissionRemap[remapServiceName][mName],
							pm,
						)
						svc := sName
						if perm.Service != nil {
							svc = *pm.Service
						}
						p.Debugf("Permission Remap added: %q → %q/%q\n", mName, svc, pm.Name)
					}
				}

				if perm.Name == "Superuser" || perm.Name == "Any" || perm.Service != nil {
					continue
				}

				if _, ok := data.Permissions[sName]; !ok {
					data.Permissions[sName] = map[string]*Perm{}
				}
				if _, ok := data.Permissions[sName][perm.Name]; !ok {
					data.Permissions[sName][perm.Name] = perm
					p.Debugf("Permission added: %q - %+v\n", mName, perm)
				} else {
					p.Debugf("Permission already in list, updating: %q - %+v\n", mName, perm)
					if len(perm.Attrs) > 0 {
						data.Permissions[sName][perm.Name].Attrs = append(data.Permissions[sName][perm.Name].Attrs, perm.Attrs...)
					}
					perm.Order = data.Permissions[sName][perm.Name].Order
				}
			}
		}
	}

	if len(data.Permissions) == 0 && len(data.PermissionRemap) == 0 {
		return nil
	}

	sort.Strings(data.PermissionServiceKeys)

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
    "github.com/fivenet-app/fivenet/v2025/pkg/config"
    "github.com/fivenet-app/fivenet/v2025/pkg/perms"
    "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
    permkeys "{{ .GoPath }}"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
	{{- range $sName, $service := .Permissions }}

		// Service: {{ $sName }}
		{{ range $perm := $service -}}
		{
			Category: permkeys.{{ serviceName $sName }}Perm,
			Name: permkeys.{{ serviceName $sName }}{{ $perm.Name }}Perm,
            Attrs: []perms.Attr{
            {{- range $attr := $perm.Attrs }}
                {
                    Key: permkeys.{{ serviceName $sName }}{{ $perm.Name }}{{ $attr.Key }}PermField,
                    Type: permissions.{{ $attr.Type }}AttributeType,
                    {{ with $attr.Valid -}}ValidValues: {{ $attr.Valid }},{{ end }}
                },
            {{- end }}
            },
            Order: {{ $perm.Order }},
            {{ with $perm.Icon }}Icon: "{{ $perm.Icon }}",{{ end }}
		},
		{{ end }}
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
    "github.com/fivenet-app/fivenet/v2025/pkg/perms"
)

{{ with .PermissionServiceKeys }}
const (
{{ range $key, $sName := . -}}
    {{ serviceName $sName }}Perm perms.Category = "{{ $sName }}"
{{ end }}

{{ range $sName, $service := $.Permissions -}}
	    {{- range $perm := $service }}
	{{ serviceName $sName }}{{ $perm.Name }}Perm perms.Name = "{{ $perm.Name }}"
            {{- range $attr := $perm.Attrs }}
    {{ serviceName $sName }}{{ $perm.Name }}{{ $attr.Key }}PermField perms.Key = "{{ $attr.Key }}"
            {{- end }}
		{{- end }}
	{{- end }}
)
{{ end }}
`

const permifyRemapTpl = `// Code generated by protoc-gen-backend. DO NOT EDIT.
{{- range $f := .FS }}
// source: {{ $f.File.InputPath }}
{{- end }}

package goproto

{{ with .PermissionRemap }}
var PermsRemap = map[string][]string{
    {{- range $service, $remap := . }}
	// Service: {{ $service }}
	{{ range $key, $target := $remap -}}
	"{{ $service }}/{{ $key }}": []string{
        {{ range $t := $target -}}
        "{{- if and (ne $t.Name "Superuser") (ne $t.Name "Any") }}{{ or $t.Service $service }}/{{ end }}{{ $t.Name }}",
        {{- end }}
    },
    {{ end }}
    {{ end }}
}
{{ end }}
`

type Perm struct {
	Service *string
	Name    string
	Attrs   []Attr
	Order   int32
	Icon    *string
}

type Attr struct {
	Key   string
	Type  string
	Valid string
}
