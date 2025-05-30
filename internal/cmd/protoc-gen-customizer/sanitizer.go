package main

import (
	"bufio"
	"strings"
	"text/template"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

// SanitizerPlugin
type SanitizerModule struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
	tpl *template.Template
}

// Sanitizer returns an initialized SanitizerPlugin
func Sanitizer() *SanitizerModule { return &SanitizerModule{ModuleBase: &pgs.ModuleBase{}} }

func (p *SanitizerModule) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	p.ctx = pgsgo.InitContext(c.Parameters())

	tpl := template.New("sanitizer").Funcs(map[string]any{
		"package": p.ctx.PackageName,
		"name":    p.ctx.Name,
		"fType":   p.ctx.Type,
	})

	p.tpl = template.Must(tpl.Parse(sanitizerTpl))
}

// Name satisfies the generator.Plugin interface.
func (p *SanitizerModule) Name() string { return "sanitizer" }

func (p *SanitizerModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, t := range targets {
		p.generate(t)
	}

	return p.Artifacts()
}

func (p *SanitizerModule) generate(f pgs.File) {
	if len(f.Messages()) == 0 {
		return
	}

	data := struct {
		F    pgs.File
		FMap map[string]map[string]*Sanitize
	}{
		F:    f,
		FMap: map[string]map[string]*Sanitize{},
	}

	for _, m := range f.Messages() {
		data.FMap[string(m.Name())] = map[string]*Sanitize{}

		for _, f := range m.Fields() {
			// Skip numeric and enum fields
			if f.Type().ProtoType().IsNumeric() || f.Type().IsEnum() || f.Type().ProtoType() == pgs.BoolT {
				continue
			}

			comment := f.SourceCodeInfo().LeadingComments()
			comment = strings.TrimSpace(comment)

			// Skip string fields without the sanitize comment
			if f.Type().ProtoType() == pgs.StringT {
				if !strings.Contains(comment, "@sanitize") {
					continue
				}
			}

			// Find comment in multiline comment
			sc := bufio.NewScanner(strings.NewReader(comment))
			for sc.Scan() {
				text := strings.TrimSpace(sc.Text())
				if strings.HasPrefix(text, "@sanitize") {
					comment = text
					break
				}
			}

			s := p.parseComment(f, comment)

			data.FMap[string(m.Name())][string(f.Name().UpperCamelCase())] = s
		}
	}

	if len(data.FMap) > 0 {
		name := p.ctx.OutputPath(f).SetExt(".sanitizer.go")
		p.AddGeneratorTemplateFile(name.String(), p.tpl, data)
	}
}

func (p *SanitizerModule) parseComment(field pgs.Field, comment string) *Sanitize {
	comment = strings.TrimPrefix(comment, "@sanitize: ")
	comment = strings.TrimPrefix(comment, "@sanitize")

	s := &Sanitize{
		Name:   field.Name().UpperCamelCase().String(),
		Method: "Sanitize",
		F:      field,
	}

	split := strings.Split(comment, ";")

	for i := range split {
		k, v, _ := strings.Cut(split[i], "=")
		if v == "" {
			continue
		}

		switch strings.ToLower(k) {
		case "method":
			s.Method = v
			continue
		}
	}

	return s
}

const sanitizerTpl = `// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: {{ .F.InputPath }}

package {{ package .F }}

import (
    "github.com/fivenet-app/fivenet/v2025/pkg/html/htmlsanitizer"
)

{{ range $key, $fields := .FMap }}
func (m *{{ $key }}) Sanitize() error {
    if m == nil {
		return nil
	}

    {{- $lastOneOf := "" }}
    {{ range $f := $fields }}
        {{- if and (ne $lastOneOf "") (or (not $f.F.InRealOneOf) (ne $lastOneOf $f.F.OneOf.Message.Name)) }}
        }
        {{- $lastOneOf = "" }}
        {{ end }}

        // Field: {{ $f.Name }}
        {{- $fType := fType $f.F }}
        {{- if $f.F.Type.IsRepeated }}
        for idx, item := range m.{{ $f.Name }} {
            _, _ = idx, item

            {{ if eq $fType.Element "string" }}
            m.{{ $f.Name }}[idx] = htmlsanitizer.{{ $f.Method }}(m.{{ $f.Name }}[idx])
            {{ else if $f.F.Type.Element.IsEmbed }}
            if v, ok := any(item).(interface{ Sanitize() error }); ok {
                if err := v.Sanitize(); err != nil {
                    return err
                }
            }
            {{ else if not (and $f.F.Type.Element (or $f.F.Type.Element.IsEnum $f.F.Type.Element.ProtoType.IsNumeric)) }}
            // ! Repeated: Element type is not a string nor embed type ({{ $fType }})
            {{ end }}
        }
        {{ else if $f.F.Type.IsMap }}
         for idx, item := range m.{{ $f.Name }} {
            _, _ = idx, item

            {{ if eq $fType.Element "string" }}
            m.{{ $f.Name }}[idx] = htmlsanitizer.{{ $f.Method }}(m.{{ $f.Name }}[idx])
            {{ else if $f.F.Type.Element.IsEmbed }}
            if v, ok := any(item).(interface{ Sanitize() error }); ok {
                if err := v.Sanitize(); err != nil {
                    return err
                }
            }
            {{ else if not (and $f.F.Type.Element (or $f.F.Type.Element.IsEnum $f.F.Type.Element.ProtoType.IsNumeric)) }}
            // ! Map: Element type is not a string nor embed type ({{ $fType }})
            {{ end }}
        }
        {{ else if $f.F.InRealOneOf }}
            {{- if ne $lastOneOf $f.F.OneOf.Message.Name }}
	        switch v := m.{{ $f.F.OneOf.Name.UpperCamelCase }}.(type) {
            {{ end }}
                case *{{ $key }}_{{ $f.Name }}:
                    if v, ok := any(v).(interface{ Sanitize() error }); ok {
                        if err := v.Sanitize(); err != nil {
                            return err
                        }
                    }

            {{- $lastOneOf = $f.F.OneOf.Message.Name }}
        {{ else if $f.F.Type.IsEmbed }}
        if m.{{ $f.Name }} != nil {
            if v, ok := any(m.Get{{ $f.Name }}()).(interface{ Sanitize() error }); ok {
                if err := v.Sanitize(); err != nil {
                    return err
                }
            }
        }
        {{ else if or (eq $fType "string") (eq $fType "*string") }}
            {{ if $f.F.HasOptionalKeyword }}
            if m.{{ $f.Name }} != nil {
            {{ end -}}
            {{ if $f.F.HasOptionalKeyword }}*{{ end }}m.{{ $f.Name }} = htmlsanitizer.{{ $f.Method }}({{ if $f.F.HasOptionalKeyword }}*{{ end }}m.{{ $f.Name }})
            {{- if $f.F.HasOptionalKeyword }}
            }
            {{- end }}
        {{ else if or $f.F.Type.IsEnum (and $f.F.Type.Element $f.F.Type.Element.IsEnum) }}
         {{/* Skip (repeated) enums */}}
        {{ else if and (not $f.F.Type.ProtoType.IsNumeric) (and $f.F.Type.Element (not $f.F.Type.Element.ProtoType.IsNumeric)) }}
        // Unhandled Type: {{ $fType }} - {{ $f.F.HasOptionalKeyword }}

        {{ end }}

    {{ end }}

    {{- if ne $lastOneOf "" }}
    }
    {{ end }}

	return nil
}
{{ end }}
`

type Sanitize struct {
	Name   string
	Method string

	F pgs.Field
}
