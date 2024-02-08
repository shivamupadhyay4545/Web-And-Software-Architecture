package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	rt.router.POST("/session", rt.wrap(rt.Dologin))

	rt.router.POST("/user/:username/follow_list", rt.wrap(rt.FollowUser))
	rt.router.DELETE("/user/:username/follow_list", rt.wrap(rt.UnfollowUser))

	rt.router.POST("/user/:username/ban_list", rt.wrap(rt.BanUser))
	rt.router.DELETE("/user/:username/ban_list", rt.wrap(rt.UnbanUser))

	rt.router.GET("/user/:username/photos/:photoId", rt.wrap(rt.GetPhoto))

	rt.router.POST("/user/:username/photos/likes", rt.wrap(rt.LikePhoto))
	rt.router.DELETE("/user/:username/photos/likes", rt.wrap(rt.UnlikePhoto))

	rt.router.POST("/user/:username/photos/comment", rt.wrap(rt.CommentPhoto))
	rt.router.DELETE("/user/:username/photos/comment", rt.wrap(rt.UncommentPhoto))

	rt.router.GET("/user/:username/profile", rt.wrap(rt.GetMyProfile))
	rt.router.DELETE("/user/:username/deleted_photos", rt.wrap(rt.RemovePhoto))

	rt.router.GET("/user/:username", rt.wrap(rt.GetMyStream))
	rt.router.PUT("/user/:username", rt.wrap(rt.UpdateUsername))
	rt.router.POST("/user/:username", rt.wrap(rt.UploadPhoto))

	return rt.router
}
