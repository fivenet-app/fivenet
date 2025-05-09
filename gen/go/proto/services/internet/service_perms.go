// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/internet/ads.proto
// source: services/internet/domain.proto
// source: services/internet/internet.proto

package internet

import (
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
)

var PermsRemap = map[string]string{
	// Service: AdsService
	"AdsService/GetAds": "Any",

	// Service: DomainService
	"DomainService/CheckDomainAvailability": "Any",
	"DomainService/ListDomains":             "Any",
	"DomainService/ListTLDs":                "Any",
	"DomainService/RegisterDomain":          "Any",
	"DomainService/UpdateDomain":            "Any",

	// Service: InternetService
	"InternetService/GetPage": "Any",
	"InternetService/Search":  "Any",
}

func init() {
	perms.AddPermsToList([]*perms.Perm{})
}
