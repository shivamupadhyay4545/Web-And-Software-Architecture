        CREATE TABLE IF NOT EXISTS users (
			username TEXT PRIMARY KEY,
			name TEXT
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
		);