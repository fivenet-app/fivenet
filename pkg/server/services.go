package server

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Service interface {
	RegisterHTTP(e *gin.Engine)
}

// AsService annotates the given constructor to state that
// it provides a HTTP service to the "httpservices" group.
func AsService(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Service)),
		fx.ResultTags(`group:"httpservices"`),
	)
}
