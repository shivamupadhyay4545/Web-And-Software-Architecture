package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func Usermain() (*sql.DB, error) {
	// Open SQLite database (creates the file if not exists)

	db, err := sql.Open("sqlite3", "wasaphoto.db")
	if err != nil {
		log.Fatal(err)
	}

	// Create a table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			name TEXT
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil

}

func UpPhoto() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "wasaphoto.db")
	if err != nil {
		log.Fatal(err)
	}
	// Create a table
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS photos (
		username TEXT NOT NULL,
		photoNum INTEGER NOT NULL,
		photo BLOB NOT NULL,
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		likes DEFAULT 0,
		comments DEFAULT 0,
		PRIMARY KEY (username, photoNum)
	);
	
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func UserComment() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "wasaphoto.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS comments (
		photoid TEXT NOT NULL,
		commentuser TEXT NOT NULL,
		comment TEXT NOT NULL, 
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY(commentuser,date_time)
	)
	
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil

}

func InsertComment(db *sql.DB, photoid string, commentuser string, comment string) {
	t := time.Now()
	_, err := db.Exec("INSERT INTO comments (photoid,commentuser,comment,date_time) VALUES (?,?,?,?)", photoid, commentuser, comment, t.Format("2017.09.07 17:06:06"))
	if err != nil {
		log.Fatal(err)
	}

}

func Userlike() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "wasaphoto.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS like (
		photoid TEXT NOT NULL,
		likeuser TEXT NOT NULL,
		PRIMARY KEY(photoid,likeuser)
	)
	
	`)
	if err != nil {
		log.Fatal()
	}
	return db, nil
}

func Follower() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "wasaphoto.db")

	if err != nil {
		log.Fatal()
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS followers (
		follower TEXT NOT NULL,
		following TEXT NOT NULL,
		PRIMARY KEY (follower, following)
	)
	`)
	if err != nil {
		log.Fatal()
	}

	return db, nil
}
func Ban() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", "wasaphoto.db")

	if err != nil {
		log.Fatal()
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS banlist (
		who TEXT NOT NULL,
		whom TEXT NOT NULL,
		PRIMARY KEY (who, whom)
	)
	`)
	if err != nil {
		log.Fatal()
	}

	return db, nil
}
