package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/controllers"
)

func PhotoRoutes(router *httprouter.Router) {
	router.POST("/user/:username/photos/likes", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.LikePhoto(w, r, ps)
	})

	router.DELETE("/user/:username/photos/likes", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.UnlikePhoto(w, r, ps)
	})

	router.POST("/user/:username/photos/comment", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.CommentPhoto(w, r, ps)
	})

	router.DELETE("/user/:username/photos/comment", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.UncommentPhoto(w, r, ps)
	})

	router.GET("/user/:username/photos/:photoId", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		controllers.GetPhoto(w, r, ps)
	})
}
