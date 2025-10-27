package handlers

import "github.com/gin-gonic/gin"

func (r *appRoute) TestRoute(path string) {
	api := r.Route.Group(path)
	api.GET("/test", r.testHandler)
}

func (r *appRoute) testHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Test successful",
	})
}
