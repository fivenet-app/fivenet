package main

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star/v2"
	pgsgo "github.com/lyft/protoc-gen-star/v2/lang/go"
)

// TesterPlugin
type TesterModule struct {
	*pgs.ModuleBase
	ctx pgsgo.Context
	tpl *template.Template
}

// Tester returns an initialized TesterPlugin
func Tester() *TesterModule { return &TesterModule{ModuleBase: &pgs.ModuleBase{}} }

func (p *TesterModule) InitContext(c pgs.BuildContext) {
	p.ModuleBase.InitContext(c)
	p.ctx = pgsgo.InitContext(c.Parameters())

	tpl := template.New("tester").Funcs(map[string]any{
		"package": p.ctx.PackageName,
		"name":    p.ctx.Name,
	})

	p.tpl = template.Must(tpl.Parse(testerTpl))
}

// Name satisfies the generator.Plugin interface.
func (p *TesterModule) Name() string { return "tester" }

func (p *TesterModule) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, t := range targets {
		p.generate(t)
	}

	return p.Artifacts()
}

func (p *TesterModule) generate(f pgs.File) {
	if len(f.Services()) == 0 {
		return
	}

	name := p.ctx.OutputPath(f).SetExt(".tester.go")
	p.AddGeneratorTemplateFile(name.String(), p.tpl, f)
}

const testerTpl = `// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: {{ .File.InputPath }}

package {{ package . }}

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

{{ range $service := .Services }}
func NewTest{{ $service.Name }}Client(srv {{ $service.Name }}Server) ({{ $service.Name }}Client, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

    buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	server := grpc.NewServer()
	Register{{ $service.Name }}Server(server, srv)
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("error serving test grpc server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to test grpc server: %v", err)
	}

	go func() {
        <-ctx.Done()
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		server.Stop()
	}()

	client := New{{ $service.Name }}Client(conn)
	return client, ctx, cancel
}
{{ end }}
`
