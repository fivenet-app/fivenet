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

const backendTestSanitizerProto = `syntax = "proto3";

package codegen.sanitizer;

import "google/protobuf/descriptor.proto";

extend google.protobuf.FieldOptions {
  FieldOptions sanitizer = 51003;
}

message FieldOptions {
  bool enabled = 1;
  optional string method = 2;
  optional bool strip_html_tags = 3;
  optional bool tiptap_json = 4;
  optional uint32 max_bytes = 5;
}
`

const backendTestOneofSanitizerProto = `syntax = "proto3";

package services.sanitizertest;

import "codegen/sanitizer/sanitizer.proto";

message OneofPayload {
  string title = 1 [(codegen.sanitizer.sanitizer) = {enabled: true}];
}

message OneofContainer {
  oneof data {
    OneofPayload payload = 1 [(codegen.sanitizer.sanitizer) = {enabled: true}];
    string label = 2 [(codegen.sanitizer.sanitizer) = {
      enabled: true
      strip_html_tags: true
    }];
  }
}
`

func loadBackendSanitizerFixture(t *testing.T) (pgs.AST, map[string]pgs.File) {
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

	writeBackendTestProto(t, root, "codegen/sanitizer/sanitizer.proto", backendTestSanitizerProto)
	writeBackendTestProto(t, root, "services/sanitizertest/test.proto", backendTestOneofSanitizerProto)

	loader := testutils.Loader{ImportPaths: []string{"."}}
	ast := loader.LoadProtos(t, "services/sanitizertest/test.proto")

	ent, ok := ast.Lookup("services/sanitizertest/test.proto")
	if !ok {
		t.Fatalf("lookup target %q", "services/sanitizertest/test.proto")
	}
	file, ok := ent.(pgs.File)
	if !ok {
		t.Fatalf("lookup target %q did not return a file", "services/sanitizertest/test.proto")
	}

	return ast, map[string]pgs.File{"services/sanitizertest/test.proto": file}
}

func TestSanitizerModule_Execute_SanitizesOneofPayloadValues(t *testing.T) {
	ast, targets := loadBackendSanitizerFixture(t)
	out := renderBackendGeneratorOutputs(t, Sanitizer(), ast, targets)

	got, ok := out[filepath.ToSlash("services/sanitizertest/test.pb.sanitizer.go")]
	if !ok {
		t.Fatalf("missing sanitizer output: %v", mapsKeys(out))
	}

	mustContainBackend(t, got, "case *OneofContainer_Payload:")
	mustContainBackend(t, got, "if v.Payload != nil {")
	mustContainBackend(t, got, "if s, ok := any(v.Payload).(interface{ Sanitize() error }); ok {")
	mustContainBackend(t, got, "if err := s.Sanitize(); err != nil {")
	mustContainBackend(t, got, "case *OneofContainer_Label:")
	mustContainBackend(t, got, "v.Label = htmlsanitizer.StripHTMLTags(v.Label)")

	if strings.Contains(got, "any(v).(interface{ Sanitize() error })") {
		t.Fatalf("oneof wrapper should not be type asserted as sanitizable:\n%s", got)
	}
}
