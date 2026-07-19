package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/lyft/protoc-gen-star/v2/testutils"
)

type frontendGenerator interface {
	InitContext(pgs.BuildContext)
	Execute(map[string]pgs.File, map[string]pgs.Package) []pgs.Artifact
}

func writeTestProto(t *testing.T, root, relPath, contents string) {
	t.Helper()

	fullPath := filepath.Join(root, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("create proto dir %q: %v", filepath.Dir(fullPath), err)
	}
	if err := os.WriteFile(fullPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("write proto %q: %v", fullPath, err)
	}
}

func loadFrontendFixture(t *testing.T) (pgs.AST, map[string]pgs.File) {
	t.Helper()

	if _, err := exec.LookPath("protoc"); err != nil {
		t.Skipf("protoc not available: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get cwd: %v", err)
	}

	root := t.TempDir()
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir to temp dir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(cwd)
	})

	writeTestProto(t, root, "resources/permissions/attributes/attributes.proto", testAttributesProto)
	writeTestProto(t, root, "codegen/perms/perms.proto", testPermissionsProto)
	writeTestProto(t, root, "services/documents/documents.proto", testDocumentsProto)
	writeTestProto(t, root, "services/jobs/conduct.proto", testConductProto)
	writeTestProto(t, root, "services/qualifications/qualifications.proto", testQualificationsProto)

	loader := testutils.Loader{ImportPaths: []string{"."}}
	ast := loader.LoadProtos(
		t,
		"services/documents/documents.proto",
		"services/jobs/conduct.proto",
		"services/qualifications/qualifications.proto",
	)

	targets := map[string]pgs.File{}
	for _, name := range []string{
		"services/documents/documents.proto",
		"services/jobs/conduct.proto",
		"services/qualifications/qualifications.proto",
	} {
		ent, ok := ast.Lookup(name)
		if !ok {
			t.Fatalf("lookup target %q", name)
		}
		file, ok := ent.(pgs.File)
		if !ok {
			t.Fatalf("lookup target %q did not return a file", name)
		}
		targets[name] = file
	}

	return ast, targets
}

func renderGeneratorOutput(t *testing.T, mod frontendGenerator, ast pgs.AST, targets map[string]pgs.File) string {
	t.Helper()

	mod.InitContext(pgs.Context(pgs.InitMockDebugger(), pgs.Parameters{}, "."))

	arts := mod.Execute(targets, ast.Packages())
	if len(arts) != 1 {
		t.Fatalf("expected one artifact, got %d", len(arts))
	}

	gtf, ok := arts[0].(pgs.GeneratorTemplateFile)
	if !ok {
		t.Fatalf("unexpected artifact type %T", arts[0])
	}

	resp, err := gtf.ProtoFile()
	if err != nil {
		t.Fatalf("render artifact: %v", err)
	}

	return resp.GetContent()
}

func mustContain(t *testing.T, haystack, needle string) {
	t.Helper()

	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected output to contain %q", needle)
	}
}
