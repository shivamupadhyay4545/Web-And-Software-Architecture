package routes

import (
	"net/http"

	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/controllers"
)

func PhotoRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/user/:username/photos/likes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			controllers.LikePhoto(w, r)
		case "DELETE":
			controllers.UnlikePhoto(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user/:username/photos/comment", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			controllers.CommentPhoto(w, r)
		case "DELETE":
			controllers.UncommentPhoto(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user/:username/photos/:PhotoId", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			controllers.GetPhoto(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
