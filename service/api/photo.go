package api

import (
	"encoding/json"
	"fmt"
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
	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.Dolike(username, Photoid, parts[0], photocode, w, ctx)
	}

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
	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.DoUnlike(username, Photoid, parts[0], photocode, w, ctx)
	}

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
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error converting photo code to integer:  %w", err).Error())
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}
	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.DoComment(username, Photoid, parts[0], photocode, comment.Content, w, ctx)
	}

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
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error converting photo code to integer:  %w", err).Error())
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}
	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.DounComment(username, Photoid, parts[0], photocode, comment.Content, w, ctx)
	}

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
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error converting photo code to integer:  %w", err).Error())
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}

	rt.db.Getphoto(username, Photoid, photocode, w, ctx)
	// token := r.Header.Get("Authorization")

	// is_valid := rt.db.Authorize(username, token, w, ctx)

	// if is_valid {
	// 	rt.db.Getphoto(username, Photoid, photocode, w, ctx)
	// }

	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

}
