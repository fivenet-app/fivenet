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

type backendGenerator interface {
	InitContext(pgs.BuildContext)
	Execute(map[string]pgs.File, map[string]pgs.Package) []pgs.Artifact
}

func writeBackendTestProto(t *testing.T, root, relPath, contents string) {
	t.Helper()

	fullPath := filepath.Join(root, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("create proto dir %q: %v", filepath.Dir(fullPath), err)
	}
	if err := os.WriteFile(fullPath, []byte(contents), 0o644); err != nil {
		t.Fatalf("write proto %q: %v", fullPath, err)
	}
}

func loadBackendFixture(t *testing.T) (pgs.AST, map[string]pgs.File) {
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

	writeBackendTestProto(t, root, "resources/permissions/attributes/attributes.proto", backendTestAttributesProto)
	writeBackendTestProto(t, root, "codegen/perms/perms.proto", backendTestPermissionsProto)
	writeBackendTestProto(t, root, "services/documents/documents.proto", backendTestDocumentsProto)
	writeBackendTestProto(t, root, "services/qualifications/qualifications.proto", backendTestQualificationsProto)

	loader := testutils.Loader{ImportPaths: []string{"."}}
	ast := loader.LoadProtos(
		t,
		"services/documents/documents.proto",
		"services/qualifications/qualifications.proto",
	)

	targets := map[string]pgs.File{}
	for _, name := range []string{
		"services/documents/documents.proto",
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

func loadBackendOverrideFixture(t *testing.T) (pgs.AST, map[string]pgs.File) {
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

	writeBackendTestProto(t, root, "resources/permissions/attributes/attributes.proto", backendTestAttributesProto)
	writeBackendTestProto(t, root, "codegen/perms/perms.proto", backendTestPermissionsProto)
	writeBackendTestProto(t, root, "services/settings/overrides.proto", backendTestOverridesProto)

	loader := testutils.Loader{ImportPaths: []string{"."}}
	ast := loader.LoadProtos(
		t,
		"services/settings/overrides.proto",
	)

	targets := map[string]pgs.File{}
	ent, ok := ast.Lookup("services/settings/overrides.proto")
	if !ok {
		t.Fatalf("lookup target %q", "services/settings/overrides.proto")
	}
	file, ok := ent.(pgs.File)
	if !ok {
		t.Fatalf("lookup target %q did not return a file", "services/settings/overrides.proto")
	}
	targets["services/settings/overrides.proto"] = file

	return ast, targets
}

func renderBackendGeneratorOutputs(t *testing.T, mod backendGenerator, ast pgs.AST, targets map[string]pgs.File) map[string]string {
	t.Helper()

	mod.InitContext(pgs.Context(pgs.InitMockDebugger(), pgs.Parameters{}, "."))

	arts := mod.Execute(targets, ast.Packages())
	out := make(map[string]string, len(arts))
	for _, art := range arts {
		gtf, ok := art.(pgs.GeneratorTemplateFile)
		if !ok {
			t.Fatalf("unexpected artifact type %T", art)
		}

		resp, err := gtf.ProtoFile()
		if err != nil {
			t.Fatalf("render artifact: %v", err)
		}

		out[resp.GetName()] = resp.GetContent()
	}

	return out
}

func mustContainBackend(t *testing.T, haystack, needle string) {
	t.Helper()

	if !strings.Contains(haystack, needle) {
		t.Fatalf("expected output to contain %q", needle)
	}
}

func mapsKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
