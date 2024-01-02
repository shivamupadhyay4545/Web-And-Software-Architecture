package routes

import (
	controllers "wasaphoto/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/session", controllers.Dologin())
	incomingRoutes.GET("/user/:username", controllers.GetMyStream()) 
	incomingRoutes.PUT("/user/:username", controllers.UpdateUsername())
	incomingRoutes.POST("/user/:username", controllers.UploadPhoto()) 
	incomingRoutes.POST("/user/:username/follow_list", controllers.FollowUser())
	incomingRoutes.DELETE("/user/:username/follow_list", controllers.UnfollowUser())
	incomingRoutes.POST("/user/:username/ban_list", controllers.BanUser())
	incomingRoutes.DELETE("/user/:username/ban_list", controllers.UnbanUser())
	incomingRoutes.GET("/user/:username/profile", controllers.GetMyProfile()) 
	incomingRoutes.DELETE("/user/:username/deleted_photos", controllers.RemovePhoto()) 
}
