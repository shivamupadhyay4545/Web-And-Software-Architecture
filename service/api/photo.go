package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/reqcontext"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/models"
)

func (rt *_router) LikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")
	Photoid := r.URL.Query().Get("Photoid")

	parts := strings.Split(Photoid, "_")
	photocode, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}
	rt.db.Dolike(username, Photoid, parts[0], photocode, w)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

}
func (rt *_router) UnlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	username := ps.ByName("username")
	Photoid := r.URL.Query().Get("Photoid")

	parts := strings.Split(Photoid, "_")
	photocode, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}
	rt.db.DoUnlike(username, Photoid, parts[0], photocode, w)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

}
func (rt *_router) CommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")
	Photoid := r.URL.Query().Get("Photoid")

	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	parts := strings.Split(Photoid, "_")
	photocode, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal("Error converting photo code to integer: ", err)
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}
	rt.db.DoComment(username, Photoid, parts[0], photocode, comment.Content, w)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

}
func (rt *_router) UncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	username := ps.ByName("username")
	Photoid := r.URL.Query().Get("Photoid")

	parts := strings.Split(Photoid, "_")
	photocode, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal("Error converting photo code to integer: ", err)
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}
	rt.db.DounComment(username, Photoid, parts[0], photocode, comment.Content, w)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

}
func (rt *_router) GetPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")
	Photoid := ps.ByName("photoId")

	parts := strings.Split(Photoid, "_")
	photocode, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal("Error converting photo code to integer: ", err)
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}
	rt.db.Getphoto(username, Photoid, photocode, w)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

}
