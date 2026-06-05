package main

import (
	"path"
	"slices"
	"strings"
	"text/template"

	permspb "github.com/fivenet-app/fivenet/v2026/gen/go/proto/codegen/perms"
	permissionsattributes "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/permissions/attributes"
	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

// PermifyPlugin is a protoc-gen-star module that generates a TypeScript file
// defining permission types based on comments in the protobuf service methods.
type PermifyModule struct {
	*pgs.ModuleBase

	ctx pgsgo.Context
	tpl *template.Template
}

// Permify returns an initialized PermifyPlugin.
func Permify() *PermifyModule {
	return &PermifyModule{ModuleBase: &pgs.ModuleBase{}}
}

func (p *PermifyModule) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	p.ctx = pgsgo.InitContext(c.Parameters())

	tpl := template.New("permify").Funcs(fns)

	p.tpl = template.Must(tpl.Parse(permsTpl))
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
		if _, ok := visited[key]; !ok {
			visited[key] = []pgs.File{t}
			continue
		}

		visited[key] = append(visited[key], t)
	}

	data := struct {
		FS          []pgs.File
		Permissions map[string]map[string]*Perm
	}{
		FS:          []pgs.File{},
		Permissions: map[string]map[string]*Perm{},
	}
	for _, fs := range visited {
		slices.SortFunc(fs, func(a, b pgs.File) int {
			return strings.Compare(a.File().InputPath().String(), b.File().InputPath().String())
		})

		for _, f := range fs {
			if len(f.Services()) == 0 {
				continue
			}

			data.FS = append(data.FS, f)

			for _, s := range f.Services() {
				svcName := strings.TrimPrefix(string(s.FullyQualifiedName()), ".services.")

				var serviceOpts permspb.ServiceOptions
				_, err := s.Extension(permspb.E_PermsSvc, &serviceOpts)
				if err != nil {
					p.Fail("error reading perms option:", err)
				}

				for _, v := range serviceOpts.AdditionalPerms {
					if _, ok := data.Permissions[svcName]; !ok {
						data.Permissions[svcName] = map[string]*Perm{}
					}

					perm := &Perm{
						Name:    v.Name,
						Service: &svcName,
						Attrs:   make([]Attr, len(v.Attrs)),
					}

					for i, a := range v.Attrs {
						atype := "StringList"
						kind := "stringList"
						switch a.Type {
						case permissionsattributes.AttributeType_ATTRIBUTE_TYPE_JOB_LIST:
							atype = "JobList"
							kind = "jobList"
						case permissionsattributes.AttributeType_ATTRIBUTE_TYPE_JOB_GRADE_LIST:
							atype = "JobGradeList"
							kind = "jobGradeList"
						}

						perm.Attrs[i] = Attr{
							Key:    a.Key,
							Type:   atype,
							Kind:   kind,
							Values: slices.Clone(a.ValidStringList),
						}
					}

					data.Permissions[svcName][perm.Name] = perm
				}
			}

			for _, s := range f.Services() {
				sName := strings.TrimPrefix(string(s.FullyQualifiedName()), ".services.")

				var serviceOpts permspb.ServiceOptions
				_, err := s.Extension(permspb.E_PermsSvc, &serviceOpts)
				if err != nil {
					p.Fail("error reading perms option:", err)
				}

				p.Debugf("Service: %s", sName)

				for _, m := range s.Methods() {
					mName := string(m.Name())
					mName = strings.TrimPrefix(mName, "services.")

					// check if the field option is present and true
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
					if val.Name != nil && *val.Name != "" {
						perm.Name = *val.Name
					}
					if val.Service != nil && *val.Service != "" {
						perm.Service = val.Service
					}

					perm.Attrs = make([]Attr, len(val.Attrs))
					for i, a := range val.Attrs {
						atype := "StringList"
						kind := "stringList"
						switch a.Type {
						case permissionsattributes.AttributeType_ATTRIBUTE_TYPE_JOB_LIST:
							atype = "JobList"
							kind = "jobList"
						case permissionsattributes.AttributeType_ATTRIBUTE_TYPE_JOB_GRADE_LIST:
							atype = "JobGradeList"
							kind = "jobGradeList"
						}

						perm.Attrs[i] = Attr{
							Key:    a.Key,
							Type:   atype,
							Kind:   kind,
							Values: slices.Clone(a.ValidStringList),
						}
					}

					if perm.Name != mName {
						continue
					}

					if perm.Name == "Superuser" || perm.Name == "Any" {
						continue
					}

					if _, ok := data.Permissions[sName]; !ok {
						data.Permissions[sName] = map[string]*Perm{}
					}
					if _, ok := data.Permissions[sName][perm.Name]; !ok {
						data.Permissions[sName][perm.Name] = perm
						p.Debugf("Permission added: %q - %+v\n", mName, perm)
					} else {
						p.Debugf("Permission already in list: %q - %+v\n", mName, perm)
					}
				}
			}
		}
	}

	slices.SortStableFunc(data.FS, func(a, b pgs.File) int {
		return strings.Compare(a.FullyQualifiedName(), b.FullyQualifiedName())
	})

	p.AddGeneratorTemplateFile(path.Join("perms.ts"), p.tpl, data)

	return p.Artifacts()
}

const permsTpl = `// Code generated by protoc-gen-fronthelper. DO NOT EDIT.
{{- range $f := .FS }}
// source: {{ $f.File.InputPath }}
{{- end }}

export type Perms = SystemPerms | GRPCServicePerms;

export type SystemPerms = 'internal.Superuser/CanBeSuperuser' | 'internal.Superuser/Superuser' | 'TODOService/TODOMethod';

export type GRPCServicePerms =
{{- range $sName, $service := $.Permissions -}}
	{{- range $i, $perm := $service }}
	| '{{ $sName }}/{{ $perm.Name }}'
	{{- end }}
{{- end -}};

export const GRPCServices = [
{{- range $sName, $service := $.Permissions }}
	'{{ $sName }}',
{{- end }}
];

export const GRPCServiceMethods = [
{{- range $sName, $service := $.Permissions -}}
	{{- range $i, $perm := $service }}
	'{{ $sName }}/{{ $perm.Name }}',
	{{- end }}
{{- end }}
];

export const PermAttributes = {
{{- range $sName, $service := $.Permissions }}
	{{- range $permName, $perm := $service }}
	'{{ $sName }}/{{ $perm.Name }}': {
		{{- range $i, $attr := $perm.Attrs }}
		'{{ $attr.Key }}': {
			type: '{{ $attr.Kind }}',
			{{- if $attr.Values }}
			values: [{{- range $j, $value := $attr.Values }}'{{ $value }}', {{- end }}] as const,
			{{- end }}
		},
		{{- end }}
	},
	{{- end }}
{{- end }}
} as const;

export type PermAttributesMap = typeof PermAttributes;
export type PermAttrKind = 'stringList' | 'jobList' | 'jobGradeList';

export type PermAttrPerm = keyof PermAttributesMap;

export type PermAttrKey<P extends Perms> = P extends keyof PermAttributesMap
	? keyof PermAttributesMap[P] & string
	: never;

export type PermAttrDescriptor<P extends Perms, K extends PermAttrKey<P>> = P extends keyof PermAttributesMap
	? K extends keyof PermAttributesMap[P]
		? PermAttributesMap[P][K]
		: never
	: never;

export type PermAttrType<P extends Perms, K extends PermAttrKey<P>> = PermAttrDescriptor<P, K>['type'];

export type PermAttrValue<P extends Perms, K extends PermAttrKey<P>> = PermAttrDescriptor<P, K> extends {
	values: readonly (infer V)[];
}
	? V
	: string;

export type PermAttrKeysByType<P extends Perms, T extends PermAttrKind> = P extends keyof PermAttributesMap
	? {
			[K in keyof PermAttributesMap[P] & string]: PermAttributesMap[P][K] extends { type: T } ? K : never;
		}[keyof PermAttributesMap[P] & string]
	: never;
`

type Perm struct {
	Service *string
	Name    string
	Attrs   []Attr
}

type Attr struct {
	Key    string
	Type   string
	Kind   string
	Values []string
}
