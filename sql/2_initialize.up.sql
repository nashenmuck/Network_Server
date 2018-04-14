CREATE TABLE IF NOT EXISTS followers (
	followeeId int NOT NULL,
	followerId int NOT NULL,
	followedwhen timestamp NOT NULL
);
CREATE TABLE IF NOT EXISTS groups (
	groupname varchar(32) NOT NULL,
	username varchar(32) NOT NULL
);
CREATE TABLE IF NOT EXISTS posts (
	username varchar(32) NOT NULL,
	body text NOT NULL,
	groupname varchar(32) NOT NULL
);
CREATE TABLE IF NOT EXISTS users (
	username varchar(32) PRIMARY KEY NOT NULL,
	password bytea NOT NULL,
	salt bytea NOT NULL,
	isadmin boolean DEFAULT FALSE,
	parentAdmin int
);
