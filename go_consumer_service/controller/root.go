package controller

import "github.com/gin-gonic/gin"

type RootController struct {
}

func NewRootController() *RootController {
	return &RootController{}
}

func (uc *RootController) RegisterRoutes(rg *gin.Engine) {
	rg.GET("/health", uc.health)
}

func (uc *RootController) health(c *gin.Context) {
	response := map[string]string{
		"status": "ok",
	}
	c.JSON(200, response)
}
