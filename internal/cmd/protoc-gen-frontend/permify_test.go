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

const testPermissionsProto = `syntax = "proto3";

package codegen.perms;

import "google/protobuf/descriptor.proto";
import "resources/permissions/attributes/attributes.proto";

extend google.protobuf.MethodOptions {
  PermsOptions perms = 51002;
}

message PermsOptions {
  bool enabled = 1;
  optional string namespace = 2;
  optional string service = 3;
  optional string name = 4;
  repeated string names = 5;
  int32 order = 6;
  repeated Attr attrs = 7;
  bool internal = 8;
}

message Attr {
  string key = 1;
  string value = 2;
  resources.permissions.attributes.AttributeType type = 3;
  repeated string valid_string_list = 4;
}

extend google.protobuf.ServiceOptions {
  ServiceOptions perms_svc = 51005;
}

message ServiceOptions {
  int32 order = 1;
  optional string icon = 2;
  optional string namespace = 3;
  optional string service = 4;
  repeated AdditionalServicePerm additional_perms = 5;
}

message AdditionalServicePerm {
  string name = 3;
  int32 order = 4;
  repeated Attr attrs = 5;
}
`

const testAttributesProto = `syntax = "proto3";

package resources.permissions.attributes;

enum AttributeType {
  ATTRIBUTE_TYPE_UNSPECIFIED = 0;
  ATTRIBUTE_TYPE_STRING_LIST = 1;
  ATTRIBUTE_TYPE_JOB_LIST = 2;
  ATTRIBUTE_TYPE_JOB_GRADE_LIST = 3;
}
`

const testDocumentsProto = `syntax = "proto3";

package services.documents;

import "codegen/perms/perms.proto";

message Empty {}

service DocumentsService {
  option (codegen.perms.perms_svc) = {
    additional_perms: [
      {
        name: "ListDocuments"
      },
      {
        name: "ListDocuments"
        attrs: [
          {
            key: "Access"
            type: 1
            valid_string_list: "Own"
            valid_string_list: "Any"
          }
        ]
      }
    ]
  };

  rpc CreateDocument(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      name: "UpdateDocument"
    };
  }

  rpc UpdateDocument(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      name: "UpdateDocument"
      attrs: [
        {
          key: "Access"
          type: 1
          valid_string_list: "Own"
          valid_string_list: "Lower_Rank"
          valid_string_list: "Same_Rank"
          valid_string_list: "Any"
        }
      ]
    };
  }
}
`

const testQualificationsProto = `syntax = "proto3";

package services.qualifications;

import "codegen/perms/perms.proto";

message Empty {}

service QualificationsService {
  rpc CreateQualification(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      name: "UpdateQualification"
    };
  }

  rpc UpdateQualification(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      name: "UpdateQualification"
      attrs: [
        {
          key: "Access"
          type: 1
          valid_string_list: "Own"
          valid_string_list: "Lower_Rank"
          valid_string_list: "Same_Rank"
          valid_string_list: "Any"
        },
        {
          key: "Fields"
          type: 1
          valid_string_list: "Public"
        }
      ]
    };
  }
}
`

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
	writeTestProto(t, root, "services/qualifications/qualifications.proto", testQualificationsProto)

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

func TestPermifyModule_Execute_RendersMergedAttrs(t *testing.T) {
	ast, targets := loadFrontendFixture(t)
	out := renderGeneratorOutput(t, Permify(), ast, targets)

	mustContain(t, out, "'documents.DocumentsService/ListDocuments': {")
	mustContain(t, out, "'documents.DocumentsService/UpdateDocument': {")
	mustContain(t, out, "'qualifications.QualificationsService/UpdateQualification': {")
	mustContain(t, out, "'Access': {")
	mustContain(t, out, "'Fields': {")
	mustContain(t, out, "values: ['Own','Any',] as const,")
	mustContain(t, out, "values: ['Own','Lower_Rank','Same_Rank','Any',] as const,")
	mustContain(t, out, "values: ['Public',] as const,")

	if got := strings.Count(out, "'documents.DocumentsService/ListDocuments': {"); got != 1 {
		t.Fatalf("expected one ListDocuments block, got %d", got)
	}
	if got := strings.Count(out, "'documents.DocumentsService/UpdateDocument': {"); got != 1 {
		t.Fatalf("expected one UpdateDocument block, got %d", got)
	}
	if got := strings.Count(out, "'qualifications.QualificationsService/UpdateQualification': {"); got != 1 {
		t.Fatalf("expected one UpdateQualification block, got %d", got)
	}
}

func TestListSvcMethodsModule_Execute_RendersServicesAndMethods(t *testing.T) {
	ast, targets := loadFrontendFixture(t)
	out := renderGeneratorOutput(t, ListSvcMethods(), ast, targets)

	mustContain(t, out, "'documents.DocumentsService',")
	mustContain(t, out, "'qualifications.QualificationsService',")
	mustContain(t, out, "'documents.DocumentsService/CreateDocument',")
	mustContain(t, out, "'documents.DocumentsService/UpdateDocument',")
	mustContain(t, out, "'qualifications.QualificationsService/CreateQualification',")
	mustContain(t, out, "'qualifications.QualificationsService/UpdateQualification',")
}

func TestClientsModule_Execute_RendersClientFactories(t *testing.T) {
	ast, targets := loadFrontendFixture(t)
	out := renderGeneratorOutput(t, Clients(), ast, targets)

	mustContain(t, out, "export async function getDocumentsDocumentsClient()")
	mustContain(t, out, "export async function getQualificationsQualificationsClient()")
	mustContain(t, out, "import('~~/gen/ts/services/documents/documents.client')")
	mustContain(t, out, "import('~~/gen/ts/services/qualifications/qualifications.client')")
}
