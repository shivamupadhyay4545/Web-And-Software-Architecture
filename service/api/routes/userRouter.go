package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/controllers"
)

func UserRoutes(router *httprouter.Router) {
	router.POST("/session", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		controllers.Dologin(w, r)
	})

	router.POST("/user/:username/follow_list", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.FollowUser(w, r, ps)
	})

	router.DELETE("/user/:username/follow_list", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.UnfollowUser(w, r, ps)
	})

	router.POST("/user/:username/ban_list", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.BanUser(w, r, ps)
	})

	router.DELETE("/user/:username/ban_list", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.UnbanUser(w, r, ps)
	})

	router.GET("/user/:username/profile", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.GetMyProfile(w, r, ps)
	})

	router.DELETE("/user/:username/deleted_photos", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.RemovePhoto(w, r, ps)
	})

	router.GET("/user/:username", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.GetMyStream(w, r, ps)
	})

	router.PUT("/user/:username", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.UpdateUsername(w, r, ps)
	})

	router.POST("/user/:username", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.UploadPhoto(w, r, ps)
	})
}
