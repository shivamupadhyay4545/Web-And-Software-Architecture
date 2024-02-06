/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"crypto/sha256"
	"database/sql"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/reqcontext"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/models"
	"github.com/sirupsen/logrus"
)

//go:embed migration.sql
var migration string

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

	Ping() error

	CreateUser(userName string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	Authorize(username string, token string, w http.ResponseWriter, ctx reqcontext.RequestContext) (is_valid bool)

	Stream(username string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	ChangeUserName(Newname string, username string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	UpPhoto(username string, fileBytes []byte, w http.ResponseWriter, ctx reqcontext.RequestContext)

	Follow(username string, following string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	Unfollow(username string, following string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	Ban(username string, banned string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	UnBan(username string, banned string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	Profile(username string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	DelPhoto(username string, photocode int, Photoid string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	Dolike(username string, Photoid string, parts string, photocode int, w http.ResponseWriter, ctx reqcontext.RequestContext)

	DoUnlike(username string, Photoid string, parts string, photocode int, w http.ResponseWriter, ctx reqcontext.RequestContext)

	Getphoto(username string, Photoid string, photocode int, w http.ResponseWriter, ctx reqcontext.RequestContext)

	DoComment(username string, Photoid string, parts string, photocode int, comment string, w http.ResponseWriter, ctx reqcontext.RequestContext)

	DounComment(username string, Photoid string, parts string, photocode int, comment string, w http.ResponseWriter, ctx reqcontext.RequestContext)
}

var write_err = "error writing response"
var row_err = `{"error": "Failed to scan row", "ERR": "`

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err := db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='example_table';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {
		sqlStmt := `CREATE TABLE example_table (id INTEGER NOT NULL PRIMARY KEY, name TEXT);`
		_, err = db.Exec(sqlStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating database structure: %w", err)
		}
	}
	if migration != "" {

		_, err = db.Exec(migration)

		if err != nil {
			logrus.Error("error executing migration.sql: ", err)
			return nil, fmt.Errorf("error executing migration.sql: %w", err)
		}

	} else {

		logrus.Error("migration.sql not found")
		qury := `
		CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			id TEXT
		);
		CREATE TABLE IF NOT EXISTS photos (
			username TEXT NOT NULL,
			photoNum INTEGER NOT NULL,
			photo BLOB NOT NULL,
			date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			likes DEFAULT 0,
			comments DEFAULT 0,
			PRIMARY KEY (username, photoNum)
		);
		CREATE TABLE IF NOT EXISTS comments (
			photoid TEXT NOT NULL,
			commentuser TEXT NOT NULL,
			comment TEXT NOT NULL, 
			date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(commentuser,date_time)
		);
		CREATE TABLE IF NOT EXISTS like (
			photoid TEXT NOT NULL,
			likeuser TEXT NOT NULL,
			PRIMARY KEY(photoid,likeuser)
		);
		CREATE TABLE IF NOT EXISTS followers (
			follower TEXT NOT NULL,
			following TEXT NOT NULL,
			PRIMARY KEY (follower, following)
		);
		CREATE TABLE IF NOT EXISTS banlist (
			who TEXT NOT NULL,
			whom TEXT NOT NULL,
			PRIMARY KEY (who, whom)
			);`

		_, err = db.Exec(qury)

		if err != nil {
			return nil, fmt.Errorf("error executing migration: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}

var server_error = "internal server error:  %w"

func (db *appdbimpl) CreateUser(userName string, w http.ResponseWriter, ctx reqcontext.RequestContext) {

	var hashString string
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", userName).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Database query error"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	if count == 0 {

		hasher := sha256.New()

		// Write the username to the hash
		hasher.Write([]byte(userName))

		// Get the final hash as a byte slice
		hashBytes := hasher.Sum(nil)

		// Convert the hash to a hexadecimal string
		hashString = hex.EncodeToString(hashBytes)

		_, err := db.c.Exec("INSERT OR IGNORE INTO users (username, id) VALUES (?, ?)", userName, hashString)
		if err != nil {
			ctx.Logger.WithError(err).Error(
				fmt.Errorf(server_error, err).Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(`{"error": "Internal Server Error"}`))
			if err != nil {
				ctx.Logger.WithError(err).Error(write_err)
			}
			return
		}
	} else {
		err = db.c.QueryRow(`SELECT id FROM users WHERE username = ?`, userName).Scan(&hashString)
		if err != nil {
			ctx.Logger.WithError(err).Error(
				fmt.Errorf(server_error, err).Error())
			http.Error(w, `{"error": "Failed to scan hashstring"}`, http.StatusInternalServerError)
			return
		}

	}
	responseJSON := map[string]interface{}{
		"authtoken": hashString,
	}
	fmt.Println(hashString)
	responseBytes, err := json.Marshal(responseJSON)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("unmarshal error:  %w", err).Error())
		http.Error(w, `{"error": "Failed to marshal response JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}

}

func (db *appdbimpl) Authorize(username string, token string, w http.ResponseWriter, ctx reqcontext.RequestContext) (is_valid bool) {
	count := 0

	err := db.c.QueryRow(`SELECT COUNT(*) FROM users WHERE id = ? AND username = ?`, token, username).Scan(&count)
	is_valid = count == 1
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte("error validating user"))

		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}

		ctx.Logger.WithError(err).Error("error validating user")
		return false
	}

	if !is_valid {
		fmt.Println(token)
		w.WriteHeader(http.StatusUnauthorized)
		_, err := w.Write([]byte("invalid token"))

		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}

		ctx.Logger.WithError(err).Error(token)
		return is_valid
	}

	return is_valid

}

func (db *appdbimpl) Stream(userName string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	var Ignore struct {
		ignore int
	}
	rows, err := db.c.Query("SELECT photos.username,photos.photoNum,photos.photo, photos.date_time,photos.likes,photos.comments FROM photos INNER JOIN followers ON photos.username = followers.following WHERE followers.follower = ?", userName)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error converting photo code to integer:  %w", err).Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Failed to execute query"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}

		return
	}
	defer rows.Close()
	var photos []models.Photo
	for rows.Next() {
		var photo models.Photo
		err := rows.Scan(&photo.Username, &Ignore.ignore, &photo.Photobytes, &photo.CreatedAt, &photo.Likes, &photo.NoComments)
		if err != nil {
			ctx.Logger.WithError(err).Error(
				fmt.Errorf("failed to scan row:  %w", err).Error())
			http.Error(w, row_err+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		strphotonum := strconv.Itoa(Ignore.ignore)
		photo.PhotoId = photo.Username + "_" + strphotonum
		var likec int
		err = db.c.QueryRow("SELECT COUNT(*) FROM like WHERE likeuser = ? AND photoid = ?", userName, photo.PhotoId).Scan(&likec)
		if err != nil {
			ctx.Logger.WithError(err).Error(
				fmt.Errorf(server_error, err).Error())
		}
		photo.Liked = likec
		// if err != nil {
		// 	log.Fatal(err)
		// 	w.WriteHeader(http.StatusInternalServerError)
		// 	_, err := w.Write([]byte(row_err + err.Error() ))
		// 	if err != nil {
		// 		ctx.Logger.WithError(err).Error(write_err)
		// 	}
		// 	return
		//}
		photos = append(photos, photo)
	}
	if err := rows.Err(); err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(row_err + err.Error()))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}
	// Return the result to the client
	responseJSON, err := json.Marshal(map[string]interface{}{"photos": photos})
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("marshaling error:  %w", err).Error())
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Failed to marshal response"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseJSON)
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}

func (db *appdbimpl) ChangeUserName(Newname string, username string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Database query error"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	if count != 1 {
		w.WriteHeader(http.StatusConflict)
		_, err := w.Write([]byte(`{"error": "username conflict error"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	// Perform the username update
	_, err = db.c.Exec("UPDATE users SET username = ? WHERE  username = ?", Newname, username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Error while updating database"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "Username updated successfully"}`))
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}

func (db *appdbimpl) UpPhoto(username string, fileBytes []byte, w http.ResponseWriter, ctx reqcontext.RequestContext) {

	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM photos WHERE username = ?", username).Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Database query error"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	// Increment the count for the new photo
	t := time.Now()
	count = count + 1

	// Insert the new photo into the database
	_, err = db.c.Exec("INSERT INTO photos (username, photoNum, photo, date_time) VALUES (?,?,?,?)", username, count, fileBytes, t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := w.Write([]byte(`{"error": "Failed to store the file in the database", "details": "` + err.Error() + `"}`))
		if err != nil {
			ctx.Logger.WithError(err).Error(write_err)
		}
		return
	}

	// Respond with a success message
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(`{"message": "Successfully uploaded and stored the photo", "user": "` + username + `"}`))
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}

func (db *appdbimpl) Follow(username string, following string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM banlist WHERE who = ? AND whom = ?", username, following).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if count != 0 {
		http.Error(w, `{"error": "User in ban list"}`, http.StatusConflict)
		return
	}

	err = db.c.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", following).Scan(&count)
	if err != nil {
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}

	if count != 1 {
		http.Error(w, `{"error": "User does not exist"}`, http.StatusConflict)
		return
	}

	_, err = db.c.Exec("INSERT OR IGNORE INTO followers (follower, following) VALUES (?,?)", username, following)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "User successfully followed"}`))
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}

func (db *appdbimpl) Unfollow(username string, following string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	var count int
	err := db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE follower = ? AND following = ?", username, following).Scan(&count)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
	}

	if count != 1 {
		http.Error(w, `{"error": "You never followed that user"}`, http.StatusConflict)
		return
	}

	_, err = db.c.Exec("DELETE FROM followers WHERE follower = ? AND following = ?", username, following)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error unfollowing user:  %w", err).Error())
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "User successfully unfollowed"}`))
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}

func (db *appdbimpl) Ban(username string, banned string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	_, err := db.c.Exec("INSERT OR IGNORE INTO banlist (who, whom) VALUES (?, ?)", username, banned)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error while inserting into banlist  %w", err).Error())
	}

	_, err = db.c.Exec("DELETE FROM followers WHERE follower = ? AND following = ?", username, banned)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error while deleting from followlist:  %w", err).Error())
	}

	_, err = db.c.Exec("DELETE FROM followers WHERE follower = ? AND following = ?", banned, username)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error while deleting from followlist:  %w", err).Error())
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "BanUser successful"}`))
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}
func (db *appdbimpl) UnBan(username string, banned string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	_, err := db.c.Exec("DELETE FROM banlist WHERE who = ? AND whom = ?", username, banned)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("error while deleting from banlist:  %w", err).Error())
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "User unbanned"}`))
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}
func (db *appdbimpl) Profile(username string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	var profile models.Myprofile
	var Ignore struct {
		ignore int
	}
	var count int
	var followerNo int
	var followingNo int

	err := db.c.QueryRow("SELECT COUNT(*) FROM photos WHERE username = ? ORDER BY date_time DESC", username).Scan(&count)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}

	err = db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE following = ? ", username).Scan(&followerNo)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}

	err = db.c.QueryRow("SELECT COUNT(*) FROM followers WHERE follower = ? ", username).Scan(&followingNo)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}

	rows, err := db.c.Query("SELECT * FROM photos WHERE username = ? ORDER BY date_time DESC", username)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Database query error"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var photos []models.Photo

	for rows.Next() {
		var photo models.Photo
		err := rows.Scan(&photo.Username, &Ignore.ignore, &photo.Photobytes, &photo.CreatedAt, &photo.Likes, &photo.NoComments)
		if err != nil {
			ctx.Logger.WithError(err).Error(
				fmt.Errorf(server_error, err).Error())
			http.Error(w, row_err+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
		strphotonum := strconv.Itoa(Ignore.ignore)
		photo.PhotoId = photo.Username + "_" + strphotonum
		var likec int
		err = db.c.QueryRow("SELECT COUNT(*) FROM like WHERE likeuser = ? AND photoid = ?", username, photo.PhotoId).Scan(&likec)
		if err != nil {
			ctx.Logger.WithError(err).Error(
				fmt.Errorf(server_error, err).Error())
		}
		photo.Liked = likec

		if err != nil {
			ctx.Logger.WithError(err).Error(
				fmt.Errorf(server_error, err).Error())
			http.Error(w, row_err+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		photos = append(photos, photo)
	}

	if err := rows.Err(); err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, row_err+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	profile.PhotoNo = count
	profile.Followers = followerNo
	profile.Following = followingNo
	profile.Photos = photos

	response, err := json.Marshal(map[string]interface{}{"my profile": profile})
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("unmarshal error:  %w", err).Error())
		http.Error(w, `{"error": "Failed to marshal response", "ERR": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	log.Print("done", http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}

func (db *appdbimpl) DelPhoto(username string, photocode int, Photoid string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	_, err := db.c.Exec("DELETE FROM photos WHERE username = ? AND photoNum = ? ", username, photocode)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to delete photo from photos table"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.c.Exec("DELETE FROM like WHERE photoid = ? ", Photoid)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to delete photo from like table"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.c.Exec("DELETE FROM comments WHERE photoid = ?", Photoid)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to delete photo comments"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "Image removed successfully"}`))
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}

func (db *appdbimpl) Dolike(username string, Photoid string, parts string, photocode int, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	tx, err := db.c.Begin()
	if err != nil {
		http.Error(w, `{"error": "Failed to start transaction"}`, http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT OR IGNORE INTO like (photoid, likeuser) VALUES (?, ?)", Photoid, username)
	if err != nil {
		rollbackTransaction(tx, w, "Failed to execute INSERT statement")
		return
	}

	_, err = tx.Exec("UPDATE photos SET likes = likes+1 WHERE username = ? AND photoNum = ?", parts, photocode)
	if err != nil {
		rollbackTransaction(tx, w, "Failed to execute UPDATE statement")
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, `{"error": "Failed to commit transaction"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "` + strconv.Itoa(photocode) + `", "abc": "Photo liked successfully"}`))
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}
func rollbackTransaction(tx *sql.Tx, w http.ResponseWriter, errorMsg string) {
	err := tx.Rollback()
	if err != nil {
		http.Error(w, `{"error": "Error rolling back transaction: `+errorMsg+`"}`, http.StatusInternalServerError)
		return
	}
	http.Error(w, errorMsg, http.StatusInternalServerError)
}

func (db *appdbimpl) DoUnlike(username string, Photoid string, parts string, photocode int, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	_, err := db.c.Exec("DELETE FROM like WHERE photoid = ? AND likeuser = ?", Photoid, username)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete like entry"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.c.Exec("UPDATE photos SET likes = likes-1 WHERE username = ? AND photoNum = ?", parts, photocode)
	if err != nil {
		http.Error(w, `{"error": "Failed to update photo likes count"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`{"message": "Photo unliked successfully"}`))
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}

func (db *appdbimpl) Getphoto(username string, Photoid string, photocode int, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	stmt, err := db.c.Prepare("SELECT photos.photo FROM photos WHERE username = ? AND photoNum = ?")
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to prepare SQL statement"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(username, photocode)
	var photodata []byte
	err = row.Scan(&photodata)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to scan photo data"}`, http.StatusInternalServerError)
		return
	}

	photoid := Photoid
	stmt, err = db.c.Prepare("SELECT * FROM comments WHERE photoid = ? ORDER BY date_time")
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to prepare SQL statement for comments"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(photoid)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to query comments"}`, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Comment struct to represent a comment
	type Comment struct {
		Photoid     string
		CommentUser string
		Comment     string
		DateTime    time.Time
	}

	// Create a list to store comments
	var comments []Comment

	// Iterate over the rows and scan the comments
	for rows.Next() {
		var now Comment
		err := rows.Scan(&now.Photoid, &now.CommentUser, &now.Comment, &now.DateTime)
		if err != nil {
			ctx.Logger.WithError(err).Error(
				fmt.Errorf(server_error, err).Error())
			http.Error(w, `{"error": "Failed to scan comment row"}`, http.StatusInternalServerError)
			return
		}

		// Append the comment to the list
		comments = append(comments, now)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to iterate over comment rows"}`, http.StatusInternalServerError)
		return
	}

	responseJSON := map[string]interface{}{
		"photobytes": photodata,
		"comments":   comments,
	}

	responseBytes, err := json.Marshal(responseJSON)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("unmarshal error:  %w", err).Error())
		http.Error(w, `{"error": "Failed to marshal response JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}
func (db *appdbimpl) DoComment(username string, Photoid string, parts string, photocode int, comment string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	t := time.Now()
	_, err := db.c.Exec("INSERT INTO comments (photoid, commentuser, comment, date_time) VALUES (?, ?, ?, ?)", Photoid, username, comment, t)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to insert comment into database"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.c.Exec("UPDATE photos SET comments = comments+1 WHERE username = ? AND photoNum = ?", parts, photocode)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to update photo comments count"}`, http.StatusInternalServerError)
		return
	}

	responseJSON := map[string]interface{}{
		"message": "Comment added successfully",
	}

	responseBytes, err := json.Marshal(responseJSON)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("marshalling error:  %w", err).Error())
		http.Error(w, `{"error": "Failed to marshal response JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}

func (db *appdbimpl) DounComment(username string, Photoid string, parts string, photocode int, comment string, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	_, err := db.c.Exec("DELETE FROM comments WHERE photoid = ? AND commentuser = ? AND comment=?", Photoid, username, comment)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to delete comment from database"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.c.Exec("UPDATE photos SET comments = comments-1 WHERE username = ? AND photoNum =?", parts, photocode)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf(server_error, err).Error())
		http.Error(w, `{"error": "Failed to update photo comments count"}`, http.StatusInternalServerError)
		return
	}

	responseJSON := map[string]interface{}{
		"message": "Comment removed successfully",
	}

	responseBytes, err := json.Marshal(responseJSON)
	if err != nil {
		ctx.Logger.WithError(err).Error(
			fmt.Errorf("marshalling comment error:  %w", err).Error())
		http.Error(w, `{"error": "Failed to marshal response JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(responseBytes)
	if err != nil {
		ctx.Logger.WithError(err).Error(write_err)
	}
}
