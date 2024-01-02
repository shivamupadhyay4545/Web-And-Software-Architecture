package routes

import (
	"wasaphoto/controllers"

	"github.com/gin-gonic/gin"
)

func PhotoRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/user/:username/photos/likes", controllers.LikePhoto())
	incomingRoutes.DELETE("/user/:username/photos/likes", controllers.UnlikePhoto())
	incomingRoutes.POST("/user/:username/photos/comment", controllers.CommentPhoto())
	incomingRoutes.GET("/user/:username/photos/:PhotoId", controllers.GetPhoto())
	incomingRoutes.DELETE("/user/:username/photos/comment", controllers.UncommentPhoto())
}
