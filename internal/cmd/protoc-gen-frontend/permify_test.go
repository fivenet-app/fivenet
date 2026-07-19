package main

import (
	"strings"
	"testing"
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

const testConductProto = `syntax = "proto3";

package services.jobs;

import "codegen/perms/perms.proto";

message Empty {}

service ConductService {
  option (codegen.perms.perms_svc) = {
    additional_perms: [
      {
        name: "RestoreConductEntry"
      }
    ]
  };
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

func TestPermifyModule_Execute_RendersDocumentAttrs(t *testing.T) {
	ast, targets := loadFrontendFixture(t)
	out := renderGeneratorOutput(t, Permify(), ast, targets)

	mustContain(t, out, "'documents.DocumentsService/ListDocuments': {")
	mustContain(t, out, "'documents.DocumentsService/UpdateDocument': {")
	mustContain(t, out, "'Access': {")
	mustContain(t, out, "values: ['Own','Lower_Rank','Same_Rank','Any',] as const,")

	if got := strings.Count(out, "'documents.DocumentsService/ListDocuments': {"); got != 1 {
		t.Fatalf("expected one ListDocuments block, got %d", got)
	}
	if got := strings.Count(out, "'documents.DocumentsService/UpdateDocument': {"); got != 1 {
		t.Fatalf("expected one UpdateDocument block, got %d", got)
	}
}

func TestPermifyModule_Execute_RendersQualificationAttrs(t *testing.T) {
	ast, targets := loadFrontendFixture(t)
	out := renderGeneratorOutput(t, Permify(), ast, targets)

	if got := strings.Count(out, "'qualifications.QualificationsService/UpdateQualification': {"); got != 1 {
		t.Fatalf("expected one UpdateQualification block, got %d", got)
	}
	mustContain(t, out, "'qualifications.QualificationsService/UpdateQualification': {")
	mustContain(t, out, "'Access': {")
	mustContain(t, out, "'Fields': {")
	mustContain(t, out, "values: ['Own','Lower_Rank','Same_Rank','Any',] as const,")
	mustContain(t, out, "values: ['Public',] as const,")
}

func TestPermifyModule_Execute_RendersServiceAdditionalPerms(t *testing.T) {
	ast, targets := loadFrontendFixture(t)
	out := renderGeneratorOutput(t, Permify(), ast, targets)

	mustContain(t, out, "'jobs.ConductService/RestoreConductEntry': {")

	if got := strings.Count(out, "'jobs.ConductService/RestoreConductEntry': {"); got != 1 {
		t.Fatalf("expected one RestoreConductEntry block, got %d", got)
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
