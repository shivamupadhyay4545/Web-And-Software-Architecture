package api

import (
	"encoding/json"
	"fmt"
	"io"
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
var write_err = "error writing response"

func (rt *_router) Dologin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user models.Username

	// Decode JSON data coming from the request
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	// Validate the data based on the user struct
	validationErr := validate.Struct(user)
	if validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(validationErr.Error()))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}

		return
	}

	rt.db.CreateUser(user.Username, w, ctx)

	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	return
	// }

	// Respond with a success message
	// w.WriteHeader(http.StatusOK)
	// _, err := w.Write([]byte(`{"message": "User Login Done"}`))
	// if err != nil {
	// 	ctx.Logger.WithError(err).Error(write_err)
	// }
}

func (rt *_router) GetMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")

	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.Stream(username, w, ctx)
	}
}
func (rt *_router) GetMyProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")

	rt.db.Profile(username, w, ctx)

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
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error decoding username:  %w", err).Error())
		return
	}

	// Validate the data based on the 'follow' struct
	validationErr := validate.Struct(follow)
	if validationErr != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error validating following name:  %w", err).Error())
		return
	}

	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.Follow(username, follow.Following, w, ctx)
	}

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the data based on the 'follow' struct
	validationErr := validate.Struct(follow)
	if validationErr != nil {
		http.Error(w, validationErr.Error(), http.StatusBadRequest)
		return
	}
	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.Unfollow(username, follow.Following, w, ctx)
	}

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the data based on the 'ban' struct
	validationErr := validate.Struct(ban)
	if validationErr != nil {
		http.Error(w, validationErr.Error(), http.StatusBadRequest)
		return
	}
	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.Ban(username, ban.Banned, w, ctx)
	}

	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	_, err := w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	if err != nil {
	// 		ctx.Logger.WithError(err).Error(write_err)
	// 	}
	// 	return
	// }

}
func (rt *_router) UnbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	username := ps.ByName("username")

	// Read the JSON data from the request body
	var ban struct {
		Banned string `json:"banned" validate:"required,min=3,max=16"`
	}
	err := json.NewDecoder(r.Body).Decode(&ban)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the data based on the 'ban' struct
	validationErr := validate.Struct(ban)
	if validationErr != nil {
		http.Error(w, validationErr.Error(), http.StatusBadRequest)
		return
	}

	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.UnBan(username, ban.Banned, w, ctx)
	}

	// if err != nil {
	// 	log.Fatal(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	_, err := w.Write([]byte(`{"error": "Internal Server Error"}`))
	// 	if err != nil {
	// 		ctx.Logger.WithError(err).Error(write_err)
	// 	}
	// 	return
	// }
}
func (rt *_router) RemovePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
		rt.db.DelPhoto(username, photocode, Photoid, w, ctx)
	}

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
		Newname string `json:"Newname" validate:"required,min=3,max=16"`
	}

	// Decode JSON data coming from the request

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&changeName); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(err.Error()))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	// Validate the data based on the 'user' struct
	// Note: replace validate with your validation logic
	validationErr := validate.Struct(changeName)
	if validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(validationErr.Error()))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}
	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.ChangeUserName(changeName.Newname, username, w, ctx)
	}

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
		_, err := w.Write([]byte(`{"error": "Unable to parse form data"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("image")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "Missing or invalid image file in the request"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}
	defer file.Close()

	// Check if the file extension is valid
	allowedExtensions := []string{"png", "jpg", "jpeg"}
	ext := filepath.Ext(handler.Filename)
	if ext == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte(`{"error": "Invalid file name"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
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
		_, err := w.Write([]byte(`{"error": "Only PNG, JPG, and JPEG images are allowed"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	// Read the binary data from the file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Failed to read the file", "details": "` + err.Error() + `"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}
	token := r.Header.Get("Authorization")

	is_valid := rt.db.Authorize(username, token, w, ctx)

	if is_valid {
		rt.db.UpPhoto(username, fileBytes, w, ctx)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Failed to connect to the database"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}
}
