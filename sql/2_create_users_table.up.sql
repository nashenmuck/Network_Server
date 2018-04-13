CREATE TABLE IF NOT EXISTS users (
	user varchar(32) NOT NULL,
	password bytea(32) NOT NULL,
	salt bytea(32) NOT NULL,
	isadmin boolean DEFAULT FALSE,
	parentAdmin int
);
