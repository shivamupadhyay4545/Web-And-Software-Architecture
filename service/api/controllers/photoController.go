package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/api/models"
	"github.com/shivamupadhyay4545/Web-And-Software-Architecture/service/database"
)

func LikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	Photoid := r.URL.Query().Get("Photoid")

	parts := strings.Split(Photoid, "_")
	photocode, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}

	db, err := database.Userlike()
	if err != nil {
		http.Error(w, `{"error": "Failed to connect to the database"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, `{"error": "Failed to start transaction"}`, http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec("INSERT OR IGNORE INTO like (photoid, likeuser) VALUES (?, ?)", Photoid, username)
	if err != nil {
		rollbackTransaction(tx, w, "Failed to execute INSERT statement")
		return
	}

	_, err = tx.Exec("UPDATE photos SET likes = likes+1 WHERE username = ? AND photoNum = ?", parts[0], photocode)
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
	w.Write([]byte(`{"message": "` + strconv.Itoa(photocode) + `", "abc": "Photo liked successfully"}`))
}

func rollbackTransaction(tx *sql.Tx, w http.ResponseWriter, errorMsg string) {
	err := tx.Rollback()
	if err != nil {
		http.Error(w, `{"error": "Error rolling back transaction: `+errorMsg+`"}`, http.StatusInternalServerError)
		return
	}
	http.Error(w, `{"error": "`+errorMsg+`"}`, http.StatusInternalServerError)
}

func UnlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	Photoid := r.URL.Query().Get("Photoid")

	parts := strings.Split(Photoid, "_")
	photocode, err := strconv.Atoi(parts[1])
	if err != nil {
		http.Error(w, `{"error": "Failed to convert photo code to integer"}`, http.StatusInternalServerError)
		return
	}

	db, err := database.Userlike()
	if err != nil {
		http.Error(w, `{"error": "Failed to connect to the database"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM like WHERE photoid = ? AND likeuser = ?", Photoid, username)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete like entry"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE photos SET likes = likes-1 WHERE username = ? AND photoNum = ?", parts[0], photocode)
	if err != nil {
		http.Error(w, `{"error": "Failed to update photo likes count"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Photo unliked successfully"}`))
}
func GetPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username := ps.ByName("username")
	PhotoId := ps.ByName("photoId")

	parts := strings.Split(PhotoId, "_")
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

	stmt, err := db.Prepare("SELECT photos.photo FROM photos WHERE username = ? AND photoNum = ?")
	if err != nil {
		log.Fatal("Error preparing SQL statement: ", err)
		http.Error(w, `{"error": "Failed to prepare SQL statement"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	row := stmt.QueryRow(username, photocode)
	var photodata []byte
	err = row.Scan(&photodata)
	if err != nil {
		log.Fatal("Error scanning photo data: ", err)
		http.Error(w, `{"error": "Failed to scan photo data"}`, http.StatusInternalServerError)
		return
	}

	photoid := PhotoId
	stmt, err = db.Prepare("SELECT * FROM comments WHERE photoid = ? ORDER BY date_time")
	if err != nil {
		log.Fatal("Error preparing SQL statement for comments: ", err)
		http.Error(w, `{"error": "Failed to prepare SQL statement for comments"}`, http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(photoid)
	if err != nil {
		log.Fatal("Error querying comments: ", err)
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
			log.Fatal("Error scanning comment row: ", err)
			http.Error(w, `{"error": "Failed to scan comment row"}`, http.StatusInternalServerError)
			return
		}

		// Append the comment to the list
		comments = append(comments, now)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating over comment rows: ", err)
		http.Error(w, `{"error": "Failed to iterate over comment rows"}`, http.StatusInternalServerError)
		return
	}

	responseJSON := map[string]interface{}{
		"photobytes": photodata,
		"comments":   comments,
	}

	responseBytes, err := json.Marshal(responseJSON)
	if err != nil {
		log.Fatal("Error marshalling response JSON: ", err)
		http.Error(w, `{"error": "Failed to marshal response JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

func CommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	db, err := database.UserComment()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
		http.Error(w, `{"error": "Failed to connect to the database"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	t := time.Now()
	_, err = db.Exec("INSERT INTO comments (photoid, commentuser, comment, date_time) VALUES (?, ?, ?, ?)", Photoid, username, comment.Content, t)
	if err != nil {
		log.Fatal("Error inserting comment into database: ", err)
		http.Error(w, `{"error": "Failed to insert comment into database"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE photos SET comments = comments+1 WHERE username = ? AND photoNum = ?", parts[0], photocode)
	if err != nil {
		log.Fatal("Error updating photo comments count: ", err)
		http.Error(w, `{"error": "Failed to update photo comments count"}`, http.StatusInternalServerError)
		return
	}

	responseJSON := map[string]interface{}{
		"message": "Comment added successfully",
	}

	responseBytes, err := json.Marshal(responseJSON)
	if err != nil {
		log.Fatal("Error marshalling response JSON: ", err)
		http.Error(w, `{"error": "Failed to marshal response JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
func UncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	db, err := database.UserComment()
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
		http.Error(w, `{"error": "Failed to connect to the database"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM comments WHERE photoid = ? AND commentuser = ? AND comment=?", Photoid, username, comment.Content)
	if err != nil {
		log.Fatal("Error deleting comment from database: ", err)
		http.Error(w, `{"error": "Failed to delete comment from database"}`, http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE photos SET comments = comments-1 WHERE username = ? AND photoNum =?", parts[0], photocode)
	if err != nil {
		log.Fatal("Error updating photo comments count: ", err)
		http.Error(w, `{"error": "Failed to update photo comments count"}`, http.StatusInternalServerError)
		return
	}

	responseJSON := map[string]interface{}{
		"message": "Comment removed successfully",
	}

	responseBytes, err := json.Marshal(responseJSON)
	if err != nil {
		log.Fatal("Error marshalling response JSON: ", err)
		http.Error(w, `{"error": "Failed to marshal response JSON"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}
