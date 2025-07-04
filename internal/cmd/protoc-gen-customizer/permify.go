package main

import (
	"bufio"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"text/template"

	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"

	pgs "github.com/lyft/protoc-gen-star/v2"
)

// PermifyModule
type PermifyModule struct {
	*pgs.ModuleBase
	ctx      pgsgo.Context
	tpl      *template.Template
	constTpl *template.Template
}

// Permify returns an initialized PermifyModule
func Permify() *PermifyModule { return &PermifyModule{ModuleBase: &pgs.ModuleBase{}} }

func (p *PermifyModule) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	p.ctx = pgsgo.InitContext(c.Parameters())

	serviceNameFn := func(s string) string {
		_, after, _ := strings.Cut(s, ".")
		return after
	}

	tpl := template.New("permify").Funcs(map[string]any{
		"package":     p.ctx.PackageName,
		"name":        p.ctx.Name,
		"serviceName": serviceNameFn,
	})

	p.tpl = template.Must(tpl.Parse(permifyTpl))

	constTpl := template.New("permify_const").Funcs(map[string]any{
		"package":     p.ctx.PackageName,
		"name":        p.ctx.Name,
		"serviceName": serviceNameFn,
	})

	p.constTpl = template.Must(constTpl.Parse(permifyConstTpl))
}

// Name satisfies the generator.Plugin interface.
func (p *PermifyModule) Name() string { return "permify" }

func (p *PermifyModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	visited := map[string][]pgs.File{}
	for _, t := range targets {
		key := t.File().InputPath().Dir().String()
		if _, ok := visited[key]; ok {
			visited[key] = append(visited[key], t)
			continue
		}

		visited[key] = []pgs.File{t}
	}

	for _, fs := range visited {
		p.generate(fs)
	}

	return p.Artifacts()
}

func (p *PermifyModule) generate(fs []pgs.File) {
	f := fs[0]

	fqn := strings.Split(f.FullyQualifiedName(), ".")
	data := struct {
		FS                    []pgs.File
		F                     pgs.File
		GoPath                string
		PermissionServiceKeys []string
		Permissions           map[string]map[string]*Perm
		PermissionRemap       map[string]map[string]*Perm
		Attributes            map[string]map[string]*Attr
	}{
		FS:                    fs,
		F:                     f,
		GoPath:                fmt.Sprintf("github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/%s/perms", fqn[len(fqn)-1]),
		PermissionServiceKeys: []string{},
		Permissions:           map[string]map[string]*Perm{},
		PermissionRemap:       map[string]map[string]*Perm{},
	}

	slices.SortFunc(fs, func(a, b pgs.File) int {
		return strings.Compare(a.File().InputPath().String(), b.File().InputPath().String())
	})

	for _, f := range fs {
		for _, s := range f.Services() {
			sName := strings.TrimPrefix(string(s.FullyQualifiedName()), ".services.")

			data.PermissionServiceKeys = append(data.PermissionServiceKeys, sName)
			p.Debugf("Service: %s (%s)", sName, data.PermissionServiceKeys)

			for _, m := range s.Methods() {
				mName := string(m.Name())
				mName = strings.TrimPrefix(mName, "services.")

				comment := m.SourceCodeInfo().LeadingComments()
				comment = strings.TrimSpace(comment)
				if !strings.Contains(comment, "@perm") {
					continue
				}

				// Find comment in multiline comment
				sc := bufio.NewScanner(strings.NewReader(comment))
				for sc.Scan() {
					text := strings.TrimSpace(sc.Text())
					if strings.HasPrefix(text, "@perm") {
						comment = text
						break
					}
				}

				perm, err := p.parseComment(sName, mName, comment)
				if err != nil {
					p.Failf("failed to parse comment in %s method %s (comment: '%s'), error. %w", f.InputPath(), mName, comment, err)
					return
				}
				if perm == nil {
					p.Failf("failed to parse comment in %s method %s (comment: '%s')", f.InputPath(), mName, comment)
					return
				}

				if perm.Name != mName {
					remapServiceName := strings.TrimPrefix(string(s.FullyQualifiedName()), ".services.")

					if _, ok := data.PermissionRemap[remapServiceName]; !ok {
						data.PermissionRemap[remapServiceName] = map[string]*Perm{}
					}
					if _, ok := data.PermissionRemap[remapServiceName][mName]; !ok {
						data.PermissionRemap[remapServiceName][mName] = perm
						svc := sName
						if perm.Service != nil {
							svc = *perm.Service
						}
						p.Debugf("Permission Remap added: %q → %q/%q\n", mName, svc, perm.Name)
					} else {
						p.Debugf("Permission Remap already exists: %q → %q\n", mName, perm.Name)
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
		return
	}

	sort.Strings(data.PermissionServiceKeys)

	name := p.ctx.OutputPath(f)
	p.AddGeneratorTemplateFile(path.Join(filepath.Dir(name.String()), "service_perms.go"), p.tpl, data)

	constPath := path.Join(filepath.Dir(name.String()), "perms", "perms.go")
	p.AddGeneratorTemplateFile(constPath, p.constTpl, data)
}

func (p *PermifyModule) parseComment(_ string, method string, comment string) (*Perm, error) {
	comment = strings.TrimPrefix(comment, "@perm: ")
	comment = strings.TrimPrefix(comment, "@perm")

	perm := &Perm{
		Name: method,
	}

	if comment == "" {
		p.Debugf("No permission comment found for method %s, skipping", method)
		return perm, nil
	}

	split := strings.Split(comment, ";")

	for i := range split {
		k, v, _ := strings.Cut(split[i], "=")
		if v == "" {
			p.Debugf("Skipping empty permission key %s in method %s", k, method)
			continue
		}

		switch strings.ToLower(k) {
		case "name":
			if strings.Contains(v, "/") {
				split := strings.Split(v, "/")
				if len(split) != 2 {
					p.Failf("Invalid name value found: %s", v)
				}
				perm.Service = &split[0]
				perm.Name = split[1]
			} else {
				perm.Name = v
			}
			p.Debug("Parsing permission name:", v)

		case "order":
			order, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return nil, err
			}

			perm.Order = int32(order)

		case "attrs":
			for v := range strings.SplitSeq(v, "|") {
				attrSplit := strings.Split(v, "/")
				if len(attrSplit) <= 1 {
					p.Fail("Invalid attrs value found:", v)
				}

				attrType := attrSplit[1]
				validValue := ""
				validList := strings.Split(attrSplit[1], ":")
				if len(validList) > 1 {
					attrType = validList[0]
					validValue = strings.Join(validList[1:], ":")
				}

				perm.Attrs = append(perm.Attrs, Attr{
					Key:   attrSplit[0],
					Type:  attrType,
					Valid: validValue,
				})
			}
		}
		if perm.Attrs != nil {
			p.Debug("Parsing attr:", perm.Attrs)
		}
	}

	return perm, nil
}

const permifyTpl = `// Code generated by protoc-gen-customizer. DO NOT EDIT.
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

{{ with .PermissionRemap }}
var PermsRemap = map[string]string{
    {{- range $service, $remap := . }}
	// Service: {{ $service }}
	{{ range $key, $target := $remap -}}
	"{{ $service }}/{{ $key }}": "{{- if and (ne $target.Name "Superuser") (ne $target.Name "Any") }}{{ or $target.Service $service }}/{{ end }}{{ $target.Name }}",
    {{ end }}
    {{ end }}
}
{{ end }}

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
		},
		{{ end }}
	{{- end }}
	})
}
`

const permifyConstTpl = `// Code generated by protoc-gen-customizer. DO NOT EDIT.
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

type Perm struct {
	Service *string
	Name    string
	Attrs   []Attr
	Order   int32
}

type Attr struct {
	Key   string
	Type  string
	Valid string
}
