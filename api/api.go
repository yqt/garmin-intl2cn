package api

import "github.com/gin-gonic/gin"

func InitRoute(r *gin.Engine) error {
	g := r.Group("/api")

	g.GET("/sync", genSyncHandler)

	return nil
}
