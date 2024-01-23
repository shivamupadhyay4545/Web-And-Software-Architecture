package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/models"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/database"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// to close the tcp connection : lsof -i :8080 -> kill -9 <PID>

func Dologin(w http.ResponseWriter, r *http.Request) {
	// Assuming you still want to check the authorization header
	// authHeader := r.Header.Get("Authorization")
	// expectedToken := "Bearer [wasaphoto security]"
	// if authHeader != expectedToken {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	w.Write([]byte(`{"error": "Unauthorized"}`))
	// 	return
	// }

	var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

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

	// Connect to the database
	db, _ := database.Usermain()
	defer db.Close()

	// Insert or ignore into users table
	_, err := db.Exec("INSERT OR IGNORE INTO users (username, name) VALUES (?, ?)", user.Id, user.Name)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal Server Error"}`))
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User Login Done"}`))
}

func GetMyStream(w http.ResponseWriter, r *http.Request) {
	// Removed unused context and cancellation lines
	username := r.URL.Query().Get("username")
	var Ignore struct {
		ignore int
	}
	// Handle errors when opening the database connection
	db, err := database.UpPhoto()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to connect to the database"}`))
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT photos.username,photos.photoNum,photos.photo, photos.date_time,photos.likes,photos.comments FROM photos INNER JOIN followers ON photos.username = followers.following WHERE followers.follower = ?", username)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to execute query"}`))
		return
	}
	defer rows.Close()
	var photos []models.Photo
	for rows.Next() {
		var photo models.Photo
		err := rows.Scan(&photo.Username, &Ignore.ignore, &photo.Photobytes, &photo.CreatedAt, &photo.Likes, &photo.NoComments)
		strphotonum := strconv.Itoa(Ignore.ignore)
		photo.PhotoId = photo.Username + "_" + strphotonum
		photo.Liked = CheckDislikeStatus(username, photo.PhotoId)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "Failed to scan row", "ERR": "` + err.Error() + `"}`))
			return
		}
		photos = append(photos, photo)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to scan row", "ERR": "` + err.Error() + `"}`))
		return
	}
	// Return the result to the client
	responseJSON, err := json.Marshal(map[string]interface{}{"photos": photos})
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

// UpdateUsername is a placeholder for the updateUsername handler
func UpdateUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	var changeName struct {
		Name    string `json:"name" validate:"required,min=3,max=16"`
		Newname string `json:"newname" validate:"required,min=3,max=16"`
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

	// Connect to the database
	db, err := database.Usermain()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to connect to the database"}`))
		return
	}
	defer db.Close()

	// Check if the record is present or not
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Database query error"}`))
		return
	}

	if count != 1 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error": "username conflict error"}`))
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE name = ?", changeName.Name).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Database query error"}`))
		return
	}

	if count != 1 {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"error": "username conflict error"}`))
		return
	}

	// Perform the username update
	_, err = db.Exec("UPDATE users SET name = ? WHERE name = ? and username = ?", changeName.Newname, changeName.Name, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Error while updating database"}`))
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Username updated successfully"}`))
}

// UploadPhoto is a placeholder for the uploadPhoto handler

func UploadPhoto(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

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

	// Connect to the database
	db, err := database.UpPhoto()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to connect to the database"}`))
		return
	}
	defer db.Close()

	// Check the count of photos for the given username
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM photos WHERE username = ?", username).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Database query error"}`))
		return
	}

	// Increment the count for the new photo
	t := time.Now()
	count = count + 1

	// Insert the new photo into the database
	_, err = db.Exec("INSERT INTO photos (username, photoNum, photo, date_time) VALUES (?,?,?,?)", username, count, fileBytes, t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to store the file in the database", "details": "` + err.Error() + `"}`))
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Successfully uploaded and stored the photo", "user": "` + username + `"}`))
}

// FollowUser is a placeholder for the followUser handler
func FollowUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

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

	db, err := database.Follower()
	if err != nil {
		http.Error(w, `{"error": "Failed to connect to the database"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM banlist WHERE who = ? AND whom = ?", username, follow.Following).Scan(&count)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if count != 0 {
		http.Error(w, `{"error": "User in ban list"}`, http.StatusConflict)
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", follow.Following).Scan(&count)
	if err != nil {
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}

	if count != 1 {
		http.Error(w, `{"error": "User does not exist"}`, http.StatusConflict)
		return
	}

	_, err = db.Exec("INSERT OR IGNORE INTO followers (follower, following) VALUES (?,?)", username, follow.Following)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User successfully followed"}`))
}

// UnfollowUser is a placeholder for the unfollowUser handler

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

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

	db, err := database.Follower()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower = ? AND following = ?", username, follow.Following).Scan(&count)
	if err != nil {
		log.Fatal("Database query error: ", err)
	}

	if count != 1 {
		http.Error(w, `{"error": "You never followed that user"}`, http.StatusConflict)
		return
	}

	_, err = db.Exec("DELETE FROM followers WHERE follower = ? AND following = ?", username, follow.Following)
	if err != nil {
		log.Fatal("Error unfollowing user: ", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User successfully unfollowed"}`))
}

// BanUser is a placeholder for the banUser handler
func BanUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

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

	db, err := database.Ban()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT OR IGNORE INTO banlist (who, whom) VALUES (?, ?)", username, ban.Banned)
	if err != nil {
		log.Fatal("Error while inserting into banlist: ", err)
	}

	_, err = db.Exec("DELETE FROM followers WHERE follower = ? AND following = ?", username, ban.Banned)
	if err != nil {
		log.Fatal("Error while deleting from followers: ", err)
	}

	_, err = db.Exec("DELETE FROM followers WHERE follower = ? AND following = ?", ban.Banned, username)
	if err != nil {
		log.Fatal("Error while deleting from followers: ", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "BanUser successful"}`))
}

// UnbanUser is a placeholder for the unbanUser handler

func UnbanUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

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

	db, err := database.Ban()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM banlist WHERE who = ? AND whom = ?", username, ban.Banned)
	if err != nil {
		log.Fatal("Error while deleting from banlist: ", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User unbanned"}`))
}

// GetMyProfile is a placeholder for the getMyProfile handler
func GetMyProfile(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	db, err := database.UpPhoto()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}
	defer db.Close()

	var profile models.Myprofile
	var Ignore struct {
		ignore int
	}
	var count int
	var followerNo int
	var followingNo int

	err = db.QueryRow("SELECT COUNT(*) FROM photos WHERE username = ? ORDER BY date_time DESC", username).Scan(&count)
	if err != nil {
		log.Fatal("Database query error: ", err)
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE following = ? ", username).Scan(&followerNo)
	if err != nil {
		log.Fatal("Database query error: ", err)
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower = ? ", username).Scan(&followingNo)
	if err != nil {
		log.Fatal("Database query error: ", err)
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT * FROM photos WHERE username = ? ORDER BY date_time DESC", username)
	if err != nil {
		log.Fatal("Database query error: ", err)
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var photos []models.Photo

	for rows.Next() {
		var photo models.Photo
		err := rows.Scan(&photo.Username, &Ignore.ignore, &photo.Photobytes, &photo.CreatedAt, &photo.Likes, &photo.NoComments)
		strphotonum := strconv.Itoa(Ignore.ignore)
		photo.PhotoId = photo.Username + "_" + strphotonum
		photo.Liked = CheckDislikeStatus(username, photo.PhotoId)

		if err != nil {
			log.Fatal("Failed to scan row: ", err)
			http.Error(w, `{"error": "Failed to scan row", "ERR": "`+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		log.Fatal("Failed to scan row: ", err)
		http.Error(w, `{"error": "Failed to scan row", "ERR": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	profile.PhotoNo = count
	profile.Followers = followerNo
	profile.Following = followingNo
	profile.Photos = photos

	response, err := json.Marshal(map[string]interface{}{"my profile": profile})
	fmt.Println(response, profile)
	if err != nil {
		log.Fatal("Failed to marshal response: ", err)
		http.Error(w, `{"error": "Failed to marshal response", "ERR": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	log.Print("done", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	fmt.Print(response)
}

func CheckDislikeStatus(username string, photoid string) int {

	db, err := database.Userlike()

	if err != nil {
		log.Fatal(err)
		return 0
	}
	defer db.Close()

	// Query to count the number of records in the like table
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM like WHERE likeuser = ? AND photoid = ?", username, photoid).Scan(&count)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	return count

}

// RemovePhoto is a placeholder for the removePhoto handler
func RemovePhoto(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	var Photoid struct {
		Photoid string `json:"photoid" binding:"required"`
	}

	// Read the JSON data from the query parameters
	err := json.NewDecoder(r.Body).Decode(&Photoid)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	parts := strings.Split(Photoid.Photoid, "_")
	photocode, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Fatal("Error converting photo code to integer: ", err)
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}

	db, err := database.UpPhoto()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
		http.Error(w, `{"error": "Failed to connect to the database"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM photos WHERE username = ? AND photoNum = ? ", username, photocode)
	if err != nil {
		log.Fatal("Error while deleting photo from photos table: ", err)
		http.Error(w, `{"error": "Failed to delete photo from photos table"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM like WHERE photoid = ? ", Photoid.Photoid)
	if err != nil {
		log.Fatal("Error while deleting photo from like table: ", err)
		http.Error(w, `{"error": "Failed to delete photo from like table"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM comments WHERE photoid = ?", Photoid.Photoid)
	if err != nil {
		log.Fatal("Error while deleting photo comments: ", err)
		http.Error(w, `{"error": "Failed to delete photo comments"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Image removed successfully"}`))
}
