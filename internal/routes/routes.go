package routes

import (
	"github.com/Prajna1999/atlas-be/internal/service"
	"github.com/gin-gonic/gin"
)

type Routes struct {
	hetznerRoutes *HetznerRoutes
	// add other route groups
}

func NewRoutes(
	hetznerService *service.HetznerService,
) *Routes {
	return &Routes{
		hetznerRoutes: NewHetznerRoutes(hetznerService),
		// initialise other routers here
	}
}
func (r *Routes) SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		r.hetznerRoutes.Setup(api)
		// setup other routes here
	}
}
