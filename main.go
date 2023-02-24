package main

import (
	"embed"

	"github.com/galexrt/arpanet/cmd"
)

//go:embed assets/*
var assets embed.FS

//	@title			arpanet
//	@version		0.0.1
//	@description	This is a example server celler server.
//	@termsOfService	https://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	cmd.SetAssets(assets)
	cmd.Execute()
}
