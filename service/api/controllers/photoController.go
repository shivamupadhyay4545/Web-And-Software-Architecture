package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"wasaphoto/database"
	"wasaphoto/models"

	"github.com/gin-gonic/gin"
)

func LikePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var Photoid struct {
			Photoid string `json:"photoid" binding:"required"`
		}

		if err := c.ShouldBindQuery(&Photoid); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}
		parts := strings.Split(Photoid.Photoid, "_")
		photocode, err := strconv.Atoi(parts[1])
		if err != nil {
			// Handle the error if the conversion fails
			fmt.Println("Error1:", err)
			return
		}
		db, _ := database.Userlike()
		defer db.Close()
		// Begin a transaction
		tx, err := db.Begin()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction", "err": err.Error()})
			return
		}

		// Insert statement
		// Insert statement
		result, err := tx.Exec("INSERT OR IGNORE INTO like (photoid, likeuser) VALUES (?, ?)", Photoid.Photoid, username)
		if err != nil {
			// Rollback the transaction if there is an error
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute INSERT statement", "err": err.Error()})
			return
		}

		// Check if the INSERT statement affected any rows
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			// Rollback the transaction if there is an error
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get affected rows", "err": err.Error()})
			return
		}

		// Update statement (only if the INSERT statement affected a row)
		if rowsAffected > 0 {
			_, err = tx.Exec("UPDATE photos SET likes = likes+1 WHERE username = ? AND photoNum = ?", parts[0], photocode)
			if err != nil {
				// Rollback the transaction if there is an error
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute UPDATE statement", "err": err.Error()})
				return
			}
		}

		// Commit the transaction if both statements are successful
		err = tx.Commit()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction", "err": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": photocode, "abc": "Photo liked successfully"})
	}
}

func UnlikePhoto() gin.HandlerFunc {
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
			fmt.Println("Error:", err)
			return
		}
		db, _ := database.Userlike()
		defer db.Close()
		_, err = db.Exec("DELETE FROM  like WHERE photoid = ? AND likeuser = ?", Photoid.Photoid, username)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		_, err = db.Exec("UPDATE photos set likes = likes-1 WHERE username = ? AND photoNum =?", parts[0], photocode)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Photo unliked successfully"})
	}
}
func GetPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		PhotoId := c.Param("PhotoId")
		parts := strings.Split(PhotoId, "_")
		photocode, err := strconv.Atoi(parts[1])
		if err != nil {
			// Handle the error if the conversion fails
			fmt.Println("Error:", err.Error())
			return
		}
		db, _ := database.UpPhoto()
		defer db.Close()
		stmt, err := db.Prepare("SELECT photos.photo FROM photos WHERE username = ? AND photoNum = ?")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer stmt.Close()

		// Execute the SQL statement with the specified primary key values
		row := stmt.QueryRow(username, photocode)
		var photodata []byte
		row.Scan(&photodata)
		photoid := PhotoId
		stmt, err = db.Prepare("SELECT * FROM comments WHERE photoid = ? ORDER BY date_time")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer stmt.Close()

		// Execute the SQL statement with the specified photoid
		rows, err := stmt.Query(photoid)
		if err != nil {
			fmt.Println("Error:", err)
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

		// Iterate over the rows and scan the comments
		var comments []Comment
		for rows.Next() {
			var now Comment

			err := rows.Scan(&now.Photoid, &now.CommentUser, &now.Comment, &now.DateTime)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			// Append the comment to the list
			comments = append(comments, now)
		}

		// Check for errors from iterating over rows
		if err := rows.Err(); err != nil {
			fmt.Println("Error scanning row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row", "ERR": err})
			return
		}

		c.JSON(http.StatusOK, gin.H{"photobytes": photodata, "comments": comments})

	}
}

func CommentPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var Photoid struct {
			Photoid string `json:"photoid" binding:"required"`
		}
		var comment models.Comment
		if err := c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return

		}
		if err := c.ShouldBindQuery(&Photoid); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
			return
		}
		parts := strings.Split(Photoid.Photoid, "_")
		photocode, err := strconv.Atoi(parts[1])
		if err != nil {
			// Handle the error if the conversion fails
			fmt.Println("Error:", err)
			return
		}
		db, _ := database.UserComment()
		defer db.Close()
		t := time.Now()
		_, err = db.Exec("INSERT INTO  comments (photoid, commentuser, comment, date_time) VALUES (?,?,?,?)", Photoid.Photoid, username, comment.Content, t)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error3": err.Error()})
			return
		}

		_, err = db.Exec("UPDATE photos set comments = comments+1 WHERE username = ? AND photoNum =?", parts[0], photocode)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error4": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Comment added successfully"})
	}
}

func UncommentPhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		var comment models.Comment
		if err := c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
		username := c.Param("username")
		var Photoid struct {
			Photoid string `json:"photoid" binding:"required"`
		}
		if err := c.ShouldBindQuery(&Photoid); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
			return
		}
		parts := strings.Split(Photoid.Photoid, "_")
		photocode, err := strconv.Atoi(parts[1])
		if err != nil {
			// Handle the error if the conversion fails
			fmt.Println("Error:", err)
			return
		}
		db, _ := database.UserComment()
		defer db.Close()
		_, err = db.Exec("DELETE FROM comments WHERE photoid = ? AND commentuser = ? AND comment=?", Photoid.Photoid, username, comment.Content)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		_, err = db.Exec("UPDATE photos set comments = comments-1 WHERE username = ? AND photoNum =?", parts[0], photocode)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Comment removed successfully"})
	}
}
