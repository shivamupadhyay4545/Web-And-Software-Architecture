package controllers

import (
	"context"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/models"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/database"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// to close the tcp connection : lsof -i :8080 -> kill -9 <PID>

func Dologin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// authHeader := c.GetHeader("Authorization")
		// expectedToken := "Bearer [wasaphoto security]"
		// if authHeader != expectedToken {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		// 	return
		// }
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.Details

		// convert the JSON data coming from postman to something that golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// validate the data based on user struct

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		// user is not assigned the value
		db, _ := database.Usermain()
		_, err := db.Exec("INSERT OR IGNORE INTO users (username,name) VALUES (?,?)", user.Id, user.Name)
		defer db.Close()

		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"message": "user Login Done"})

	}

}

func GetMyStream() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Removed unused context and cancellation lines
		username := c.Param("username")
		var Ignore struct {
			ignore int
		}
		// Handle errors when opening the database connection
		db, err := database.UpPhoto()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to the database"})
			return
		}
		defer db.Close()

		rows, err := db.Query("SELECT photos.username,photos.photoNum,photos.photo, photos.date_time,photos.likes,photos.comments FROM photos INNER JOIN followers ON photos.username = followers.following WHERE followers.follower = ?", username)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query"})
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "ERR": err})
				return
			}
			photos = append(photos, photo)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "ERR": err})
			return
		}
		// Return the result to the client
		c.JSON(http.StatusOK, gin.H{"photos": photos})
	}
}

// UpdateUsername is a placeholder for the updateUsername handler
func UpdateUsername() gin.HandlerFunc {
	return func(c *gin.Context) {

		username := c.Param("username")
		var changeName struct {
			Name    string `json:"name" validate:"required,min=3,max=16"`
			Newname string `json:"newname" validate:"required,min=3,max=16"`
		}

		// Convert the JSON data from the request to the 'user' struct
		if err := c.ShouldBindJSON(&changeName); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		// Validate the data based on the 'user' struct
		validationErr := validate.Struct(changeName)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": validationErr.Error()})
			return
		}

		// // Validate the data based on the 'user' struct
		// validationErr = validate.Struct(changeName.Newname)
		// if validationErr != nil {
		// 	c.JSON(http.StatusBadRequest, gin.H{"error2": validationErr.Error()})
		// 	return
		// }

		db, _ := database.Usermain()
		defer db.Close()

		// Check if the record is present or not
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}

		if count != 1 {
			c.JSON(http.StatusConflict, gin.H{"error": "username conflict error"})
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE name = ?", changeName.Name).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}

		if count != 1 {
			c.JSON(http.StatusConflict, gin.H{"error": "username conflict error"})
			return
		}

		// Perform the username update
		_, err = db.Exec("UPDATE users SET name = ? WHERE name = ? and username = ?", changeName.Newname, changeName.Name, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating database"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Username updated successfully"})
	}
}

// UploadPhoto is a placeholder for the uploadPhoto handler
func UploadPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		// Parse the form data, including file uploads
		err := c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form data"})
			return
		}

		// Get the file from the form data
		file, handler, err := c.Request.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid image file in the request"})
			return
		}
		defer file.Close()

		// Check if the file extension is valid
		allowedExtensions := []string{"png", "jpg", "jpeg"}
		ext := filepath.Ext(handler.Filename)
		if ext == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file name"})
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
			c.JSON(http.StatusUnsupportedMediaType, gin.H{"error": "Only PNG, JPG, and JPEG images are allowed"})
			return
		}

		// Read the binary data from the file
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the file", "details": err.Error()})
			return
		}
		db, _ := database.UpPhoto()
		defer db.Close()

		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM photos WHERE username = ?", username).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}

		t := time.Now()
		count = count + 1
		_, err = db.Exec("INSERT INTO photos (username, photoNum, photo, date_time) VALUES (?,?,?,?)", username, count, fileBytes, t)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store the file in the database", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Successfully uploaded and stored the photo", "file": fileBytes, "user": username})

	}
}

// FollowUser is a placeholder for the followUser handler
func FollowUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var follow struct {
			Following string `json:"following" validate:"required,min=3,max=16"`
		}

		// Convert the JSON data from the request to the 'user' struct
		if err := c.ShouldBindJSON(&follow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		// Validate the data based on the 'user' struct
		validationErr := validate.Struct(follow)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		db, _ := database.Follower()
		defer db.Close()
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM banlist WHERE who = ? AND whom = ?", username, follow.Following).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count != 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "User in ban list"})
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", follow.Following).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}

		if count != 1 {
			c.JSON(http.StatusConflict, gin.H{"error": "user does not exist"})
			return
		}
		_, err = db.Exec("INSERT OR IGNORE INTO followers (follower,following) VALUES (?,?)", username, follow.Following)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User successfully followed"})

	}
}

// UnfollowUser is a placeholder for the unfollowUser handler
func UnfollowUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var follow struct {
			Following string `json:"following" validate:"required,min=3,max=16"`
		}

		// Convert the JSON data from the request to the 'user' struct
		if err := c.ShouldBindJSON(&follow); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		// Validate the data based on the 'user' struct
		validationErr := validate.Struct(follow)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		db, _ := database.Follower()
		defer db.Close()
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower = ? AND following = ?", username, follow.Following).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
			return
		}

		if count != 1 {
			c.JSON(http.StatusConflict, gin.H{"error": "you never followed that user"})
			return
		}

		_, err = db.Exec("DELETE FROM followers WHERE follower = ? and following = ?", username, follow.Following)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User successfully unfollowed"})
	}
}

// BanUser is a placeholder for the banUser handler
func BanUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		// manage input request key name first make it consistent like
		username := c.Param("username")
		var ban struct {
			Banned string `json:"banned" validate:"required,min=3,max=16"`
		}
		// Convert the JSON data from the request to the 'user' struct
		if err := c.ShouldBindJSON(&ban); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		// Validate the data based on the 'user' struct
		validationErr := validate.Struct(ban)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		db, _ := database.Ban()
		defer db.Close()

		_, err := db.Exec("INSERT OR IGNORE INTO banlist (who,whom) VALUES (?,?)", username, ban.Banned)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error while inserting": err})
			return
		}
		// Check if the new username is already taken

		_, err = db.Exec("DELETE FROM followers WHERE follower = ? and following = ?", username, ban.Banned)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error while deleting": err})
			return
		}
		_, err = db.Exec("DELETE FROM followers WHERE follower = ? and following = ?", ban.Banned, username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error while deleting": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "BanUser successful"})
	}
}

// UnbanUser is a placeholder for the unbanUser handler
func UnbanUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		var ban struct {
			Banned string `json:"banned" validate:"required,min=3,max=16"`
		}

		// Convert the JSON data from the request to the 'user' struct
		if err := c.ShouldBindJSON(&ban); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		// Validate the data based on the 'user' struct

		// Validate the data based on the 'user' struct
		validationErr := validate.Struct(ban)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		db, _ := database.Ban()
		defer db.Close()

		_, err := db.Exec("DELETE FROM banlist WHERE who = ? AND whom = ? ", username, ban.Banned)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "user unbanned"})
	}
}

// GetMyProfile is a placeholder for the getMyProfile handler
func GetMyProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		// convert the JSON data coming from postman to something that golang understands
		db, _ := database.UpPhoto()
		defer db.Close()
		var profile models.Myprofile
		var Ignore struct {
			ignore int
		}
		var count int
		var followerNo int
		var followingNo int
		err := db.QueryRow("SELECT COUNT(*) FROM photos WHERE username = ? ORDER BY date_time DESC", username).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error1": "Database query error"})
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE following = ? ", username).Scan(&followerNo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error2": "Database query error"})
			return
		}
		err = db.QueryRow("SELECT COUNT(*) FROM followers WHERE follower = ? ", username).Scan(&followingNo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error3": "Database query error"})
			return
		}
		rows, err := db.Query("SELECT * FROM photos WHERE username = ? ORDER BY date_time DESC", username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error4": "Database query error"})
			return
		}
		var photos []models.Photo
		for rows.Next() {
			var photo models.Photo
			err := rows.Scan(&photo.Username, &Ignore.ignore, &photo.Photobytes, &photo.CreatedAt, &photo.Likes, &photo.NoComments)
			strphotonum := strconv.Itoa(Ignore.ignore)
			photo.PhotoId = photo.Username + "_" + strphotonum
			photo.Liked = CheckDislikeStatus(username, photo.PhotoId)
			if err != nil {
				log.Fatal(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "ERR": err})
				return
			}
			photos = append(photos, photo)
			// profile.PhotoNo = count
			// profile.Followers = followerNo
			// profile.Following = followingNo
			// profile.Photos = photos
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "ERR": err})
			return
		}
		profile.PhotoNo = count
		profile.Followers = followerNo
		profile.Following = followingNo
		profile.Photos = photos
		c.JSON(http.StatusOK, gin.H{"my profile": profile})
	}
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
func RemovePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var Photoid struct {
			Photoid string `json:"photoid" binding:"required"`
		}

		// Convert the JSON data from the request to the 'user' struct
		if err := c.ShouldBindQuery(&Photoid); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
			return
		}
		parts := strings.Split(Photoid.Photoid, "_")
		photocode, err := strconv.Atoi(parts[1])
		if err != nil {
			// Handle the error if the conversion fails
			log.Fatal(err)
			return
		}
		db, _ := database.UpPhoto()
		defer db.Close()
		_, err = db.Exec("DELETE FROM photos WHERE username = ? AND photoNum = ? ", username, photocode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		_, err = db.Exec("DELETE FROM like WHERE photoid = ? ", Photoid.Photoid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		_, err = db.Exec("DELETE FROM comments WHERE photoid = ?", Photoid.Photoid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "image removed successfully"})
	}
}
