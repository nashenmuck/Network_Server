CREATE TABLE IF NOT EXISTS followers (
	followeeId int NOT NULL,
	followerId int NOT NULL, 
	followedwhen timestamp NOT NULL
);
CREATE TABLE IF NOT EXISTS groups (
	groupname varchar(32) NOT NULL UNIQUE,
	username varchar(32) NOT NULL
);
CREATE TABLE IF NOT EXISTS posts (
	username varchar(32) NOT NULL,
	body text NOT NULL,
	groupname varchar(32) NOT NULL,
	server varchar NOT NULL
);
CREATE TABLE IF NOT EXISTS users (
	username varchar(32) PRIMARY KEY NOT NULL UNIQUE,
	password bytea NOT NULL,
	salt bytea NOT NULL UNIQUE,
	isadmin boolean DEFAULT FALSE,
	parentAdmin varchar(32)
);
CREATE TABLE IF NOT EXISTS servers (
	server varchar NOT NULL UNIQUE,
	blacklisted boolean DEFAULT FALSE,
	whitelisted boolean DEFAULT TRUE
);
CREATE TABLE IF NOT EXISTS authtokens (
	username varchar(32) NOT NULL,
	token bytea NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS regtokens (
	issuer varchar(32) NOT NULL,
	token bytea NOT NULL UNIQUE
);


