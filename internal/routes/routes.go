package routes

import (
	"github.com/gin-gonic/gin"
)

type Routes struct {
}

func NewRoutes() *Routes {
	return &Routes{}
}
func (r *Routes) SetupRoutes(router *gin.Engine) {

}
