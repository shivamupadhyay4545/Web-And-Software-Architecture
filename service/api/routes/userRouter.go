package routes

import (
	"net/http"

	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/controllers"
)

func UserRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			controllers.Dologin(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			controllers.GetMyStream(w, r)
		case "PUT":
			controllers.UpdateUsername(w, r)
		case "POST":
			controllers.UploadPhoto(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user/follow_list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			controllers.FollowUser(w, r)
		case "DELETE":
			controllers.UnfollowUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user/ban_list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			controllers.BanUser(w, r)
		case "DELETE":
			controllers.UnbanUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user/profile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			controllers.GetMyProfile(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/user/deleted_photos", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			controllers.RemovePhoto(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
