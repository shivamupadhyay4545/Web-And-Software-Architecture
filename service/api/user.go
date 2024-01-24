package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/reqcontext"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/models"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func (rt *_router) Dologin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user models.Details

	// Decode JSON data coming from the request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	// Validate the data based on the user struct
	// Note: replace validate with your validation logic
	validationErr := validate.Struct(user)
	if validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "` + validationErr.Error() + `"}`))
		return
	}
	rt.db.CreateUser(user.Name, user.Id, w)

	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User Login Done"}`))
}

func (rt *_router) GetMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")
	rt.db.Stream(username, w)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }
}
func (rt *_router) GetMyProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")
	rt.db.Profile(username, w)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

}
func (rt *_router) FollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")

	// Read the JSON data from the request body
	var follow struct {
		Following string `json:"following" validate:"required,min=3,max=16"`
	}
	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Validate the data based on the 'follow' struct
	validationErr := validate.Struct(follow)
	if validationErr != nil {
		log.Fatal(validationErr)
		return
	}
	rt.db.Follow(username, follow.Following, w)
	// if err != nil {
	// 	http.Error(w, `{"error": "Failed to connect to the database"}`, http.StatusInternalServerError)
	// 	return
	// }

}
func (rt *_router) UnfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")

	// Read the JSON data from the request body
	var follow struct {
		Following string `json:"following" validate:"required,min=3,max=16"`
	}
	err := json.NewDecoder(r.Body).Decode(&follow)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Validate the data based on the 'follow' struct
	validationErr := validate.Struct(follow)
	if validationErr != nil {
		http.Error(w, `{"error": "`+validationErr.Error()+`"}`, http.StatusBadRequest)
		return
	}
	rt.db.Unfollow(username, follow.Following, w)
	// if err != nil {
	// 	http.Error(w, `{"error": "Failed to connect to the database"}`, http.StatusInternalServerError)
	// 	return
	// }

}
func (rt *_router) BanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")

	// Read the JSON data from the request body
	var ban struct {
		Banned string `json:"banned" validate:"required,min=3,max=16"`
	}
	err := json.NewDecoder(r.Body).Decode(&ban)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Validate the data based on the 'ban' struct
	validationErr := validate.Struct(ban)
	if validationErr != nil {
		http.Error(w, `{"error": "`+validationErr.Error()+`"}`, http.StatusBadRequest)
		return
	}
	rt.db.Ban(username, ban.Banned, w)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}

}
func (rt *_router) UnbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")

	// Read the JSON data from the request body
	var ban struct {
		Banned string `json:"banned" validate:"required,min=3,max=16"`
	}
	err := json.NewDecoder(r.Body).Decode(&ban)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Validate the data based on the 'ban' struct
	validationErr := validate.Struct(ban)
	if validationErr != nil {
		http.Error(w, `{"error": "`+validationErr.Error()+`"}`, http.StatusBadRequest)
		return
	}
	rt.db.UnBan(username, ban.Banned, w)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}
}
func (rt *_router) RemovePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")

	Photoid := r.URL.Query().Get("Photoid")

	parts := strings.Split(Photoid, "_")
	photocode, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal("Error converting photo code to integer: ", err)
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}
	rt.db.DelPhoto(username, photocode, Photoid, w)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

}
func (rt *_router) UpdateUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")

	var changeName struct {
		Name    string `json:"Name" validate:"required,min=3,max=16"`
		Newname string `json:"Newname" validate:"required,min=3,max=16"`
	}

	// Decode JSON data coming from the request

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&changeName); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error1": "` + err.Error() + `"}`))
		return
	}

	// Validate the data based on the 'user' struct
	// Note: replace validate with your validation logic
	validationErr := validate.Struct(changeName)
	if validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error1": "` + validationErr.Error() + `"}`))
		return
	}
	rt.db.ChangeUserName(changeName.Name, changeName.Newname, username, w)
	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

}
func (rt *_router) UploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")
	// Parse the form data, including file uploads
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Unable to parse form data"}`))
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Missing or invalid image file in the request"}`))
		return
	}
	defer file.Close()

	// Check if the file extension is valid
	allowedExtensions := []string{"png", "jpg", "jpeg"}
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid file name"}`))
		return
	}
	ext = strings.ToLower(filepath.Ext(handler.Filename)[1:])
	isValidExtension := false
	for _, allowedExt := range allowedExtensions {
		if ext == allowedExt {
			isValidExtension = true
			break
		}
	}
	if !isValidExtension {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(`{"error": "Only PNG, JPG, and JPEG images are allowed"}`))
		return
	}

	// Read the binary data from the file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to read the file", "details": "` + err.Error() + `"}`))
		return
	}
	rt.db.UpPhoto(username, fileBytes, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to connect to the database"}`))
		return
	}
}
