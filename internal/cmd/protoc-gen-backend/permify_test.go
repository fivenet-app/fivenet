package main

import (
	"strings"
	"testing"
)

const backendTestPermissionsProto = `syntax = "proto3";

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

const backendTestAttributesProto = `syntax = "proto3";

package resources.permissions.attributes;

enum AttributeType {
  ATTRIBUTE_TYPE_UNSPECIFIED = 0;
  ATTRIBUTE_TYPE_STRING_LIST = 1;
  ATTRIBUTE_TYPE_JOB_LIST = 2;
  ATTRIBUTE_TYPE_JOB_GRADE_LIST = 3;
}
`

const backendTestDocumentsProto = `syntax = "proto3";

package services.documents;

import "codegen/perms/perms.proto";

message Empty {}

service DocumentsService {
  option (codegen.perms.perms_svc) = {
    order: 5
    icon: "i-docs"
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
      },
      {
        name: "RestoreDocument"
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

  rpc UploadFile(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      namespace: "documents"
      service: "DocumentsService"
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

  rpc SearchDocument(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      names: "ArchiveDocument"
      names: "FindDocument"
    };
  }

  rpc DeleteDocument(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      names: "DeleteDocument"
      names: "RestoreDocument"
    };
  }

  rpc ViewAnyDocument(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      name: "Any"
    };
  }

  rpc ManageJobAdmins(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      name: "JobAdmin"
    };
  }

  rpc ManageConfigAdmins(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
      name: "ConfigAdmin"
    };
  }
}
`

const backendTestQualificationsProto = `syntax = "proto3";

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

const backendTestOverridesProto = `syntax = "proto3";

package services.settings;

import "codegen/perms/perms.proto";

message Empty {}

service OverridesService {
  option (codegen.perms.perms_svc) = {
    namespace: "override"
    service: "OverrideService"
    order: 42
    icon: "i-override"
    additional_perms: [
      {
        name: "ListThings"
      }
    ]
  };

  rpc UpdateSettings(Empty) returns (Empty) {
    option (codegen.perms.perms) = {
      enabled: true
    };
  }
}
`

func TestPermifyModule_Execute_MergesDocumentAttrs(t *testing.T) {
	ast, targets := loadBackendFixture(t)
	out := renderBackendGeneratorOutputs(t, Permify(), ast, targets)

	serviceOut, ok := out["services/documents/service_perms.go"]
	if !ok {
		t.Fatalf("missing service perms output: %v", mapsKeys(out))
	}
	constOut, ok := out["services/documents/perms/perms.go"]
	if !ok {
		t.Fatalf("missing const output: %v", mapsKeys(out))
	}

	mustContainBackend(t, serviceOut, "Name: permkeys.DocumentsServiceUpdateDocumentPerm,")
	mustContainBackend(t, serviceOut, "Key: permkeys.DocumentsServiceUpdateDocumentAccessPermField,")
	mustContainBackend(t, serviceOut, "ValidValues: []string{\"Own\", \"Lower_Rank\", \"Same_Rank\", \"Any\", },")
	mustContainBackend(t, constOut, "DocumentsServiceUpdateDocumentPerm perms.Name = \"UpdateDocument\"")
}

func TestPermifyModule_Execute_MergesQualificationAttrs(t *testing.T) {
	ast, targets := loadBackendFixture(t)
	out := renderBackendGeneratorOutputs(t, Permify(), ast, targets)

	qualServiceOut, ok := out["services/qualifications/service_perms.go"]
	if !ok {
		t.Fatalf("missing qualifications service perms output: %v", mapsKeys(out))
	}
	qualConstOut, ok := out["services/qualifications/perms/perms.go"]
	if !ok {
		t.Fatalf("missing qualifications const output: %v", mapsKeys(out))
	}

	mustContainBackend(t, qualServiceOut, "Name: permkeys.QualificationsServiceUpdateQualificationPerm,")
	mustContainBackend(t, qualServiceOut, "Key: permkeys.QualificationsServiceUpdateQualificationAccessPermField,")
	mustContainBackend(t, qualServiceOut, "Key: permkeys.QualificationsServiceUpdateQualificationFieldsPermField,")
	mustContainBackend(t, qualConstOut, "QualificationsServiceUpdateQualificationAccessPermField perms.Key = \"Access\"")
	mustContainBackend(t, qualConstOut, "QualificationsServiceUpdateQualificationFieldsPermField perms.Key = \"Fields\"")
	mustContainBackend(t, qualConstOut, "QualificationsServiceUpdateQualificationPerm perms.Name = \"UpdateQualification\"")
}

func TestPermifyModule_Execute_RendersAliasAndSpecialRemaps(t *testing.T) {
	ast, targets := loadBackendFixture(t)
	out := renderBackendGeneratorOutputs(t, Permify(), ast, targets)

	remapOut, ok := out["perms_remap.go"]
	if !ok {
		t.Fatalf("missing remap output: %v", mapsKeys(out))
	}

	mustContainBackend(t, remapOut, "\"documents.DocumentsService/UploadFile\": {")
	mustContainBackend(t, remapOut, "permsdocuments.DocumentsService.UpdateDocument.Perm")
	mustContainBackend(t, remapOut, "\"documents.DocumentsService/SearchDocument\": {")
	mustContainBackend(t, remapOut, "permsdocuments.DocumentsService.ArchiveDocument.Perm")
	mustContainBackend(t, remapOut, "permsdocuments.DocumentsService.FindDocument.Perm")
	mustContainBackend(t, remapOut, "\"documents.DocumentsService/DeleteDocument\": {")
	mustContainBackend(t, remapOut, "permsdocuments.DocumentsService.DeleteDocument.Perm")
	mustContainBackend(t, remapOut, "permsdocuments.DocumentsService.RestoreDocument.Perm")
	mustContainBackend(t, remapOut, "\"documents.DocumentsService/ViewAnyDocument\": {")
	mustContainBackend(t, remapOut, "perms.PermAnyRef")
	mustContainBackend(t, remapOut, "\"documents.DocumentsService/ManageJobAdmins\": {")
	mustContainBackend(t, remapOut, "perms.PermJobAdminRef")
	mustContainBackend(t, remapOut, "\"documents.DocumentsService/ManageConfigAdmins\": {")
	mustContainBackend(t, remapOut, "perms.PermConfigAdminRef")

	if got := strings.Count(remapOut, "\"documents.DocumentsService/SearchDocument\": {"); got != 1 {
		t.Fatalf("expected SearchDocument remap entry once, got %d", got)
	}
	if got := strings.Count(remapOut, "perms.PermAnyRef"); got < 1 {
		t.Fatalf("expected Any special remap to be rendered")
	}
	if got := strings.Count(remapOut, "perms.PermJobAdminRef"); got < 1 {
		t.Fatalf("expected JobAdmin special remap to be rendered")
	}
	if got := strings.Count(remapOut, "perms.PermConfigAdminRef"); got < 1 {
		t.Fatalf("expected ConfigAdmin special remap to be rendered")
	}
}

func TestPermifyModule_Execute_RendersNamespaceServiceOverrideAndServiceMetadata(t *testing.T) {
	ast, targets := loadBackendOverrideFixture(t)
	out := renderBackendGeneratorOutputs(t, Permify(), ast, targets)

	serviceOut, ok := out["services/settings/service_perms.go"]
	if !ok {
		t.Fatalf("missing service perms output: %v", mapsKeys(out))
	}
	constOut, ok := out["services/settings/perms/perms.go"]
	if !ok {
		t.Fatalf("missing const output: %v", mapsKeys(out))
	}
	remapOut, ok := out["perms_remap.go"]
	if !ok {
		t.Fatalf("missing remap output: %v", mapsKeys(out))
	}

	mustContainBackend(t, serviceOut, "// Service: override.OverrideService")
	mustContainBackend(t, serviceOut, "Name: permkeys.OverrideServiceListThingsPerm,")
	mustContainBackend(t, serviceOut, "Order: 4200,")
	mustContainBackend(t, serviceOut, "Icon: \"i-override\",")

	mustContainBackend(t, constOut, "OverrideServiceListThingsPerm perms.Name = \"ListThings\"")
	mustContainBackend(t, remapOut, "\"settings.OverridesService/UpdateSettings\": {")
	mustContainBackend(t, remapOut, "permsoverride.OverrideService.UpdateSettings.Perm")
	mustContainBackend(t, remapOut, "permsoverride \"github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/override/perms\"")
}
