package path

import (
	"github.com/gin-gonic/gin"

	"github.com/roguexray007/loan-app/internal/constants"
)

type pathAuth string

var (
	PathAuthAdmin   pathAuth = "admin"
	PathAuthPrivate pathAuth = "private"
	PathAuthPublic  pathAuth = "public"
)

type Endpoint struct {
	Method  string
	Path    string
	Auth    pathAuth
	Handler constants.HandlerFunc
}

func (ep *Endpoint) IsAdmin() bool {
	return ep.Auth == PathAuthAdmin
}

func (ep *Endpoint) IsPrivate() bool {
	return ep.Auth == PathAuthPrivate
}

func (ep *Endpoint) IsPublic() bool {
	return ep.Auth == PathAuthPublic
}

func (ep *Endpoint) Authenticate(gctx *gin.Context) error {
	if ep.IsAdmin() {
		//return (&Admin{Path: ep}).Request(gctx.Request)
	}
	return (&Authenticate{path: ep}).Request(gctx.Request)
}
