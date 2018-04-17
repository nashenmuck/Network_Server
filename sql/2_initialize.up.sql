CREATE TABLE IF NOT EXISTS followers (
	followeeId varchar(32) NOT NULL REFERENCES users(username),
	followee_server varchar NOT NULL REFERENCES servers(server),
	followerId varchar(32) NOT NULL REFERENCES users(username), 
	follower_server varchar NOT NULL REFERENCES servers(servers),
	followedwhen timestamp NOT NULL
);
CREATE TABLE IF NOT EXISTS group_followers (
	group_id NOT NULL REFERENCES groups(group_id),
	follower varchar(32) NOT NULL REFERENCES user(username),
	follower_server varchar references servers(server),
	UNIQUE (group_id, follower)
);
CREATE TABLE IF NOT EXISTS groups (
	group_id int PRIMARY KEY,
	owner varchar(32) NOT NULL REFERENCES users(username),
	PRIMARY KEY (group_id, owner),
);
CREATE TABLE IF NOT EXISTS posts (
	username varchar(32) NOT NULL REFERENCES users(username),
	body text NOT NULL,
	groupid int NOT NULL REFERENCES groups(group_id),
	is_special_group boolean NOT NULL,
	origin_server varchar NOT NULL REFERENCES servers(server),
);
CREATE TABLE IF NOT EXISTS users (
	username varchar(32) NOT NULL UNIQUE PRIMARY KEY,
	password bytea NOT NULL,
	salt bytea NOT NULL UNIQUE,
	isadmin boolean DEFAULT FALSE,
	parentAdmin varchar(32),
	canonical_user varchar(32) NOT NULL,
	canonical_server varchar NOT NULL REFERENCES servers(server)
);
CREATE TABLE IF NOT EXISTS servers (
	server varchar NOT NULL UNIQUE,
	blacklisted boolean DEFAULT FALSE,
	whitelisted boolean DEFAULT FALSE
);
CREATE TABLE IF NOT EXISTS authtokens (
	username varchar(32) NOT NULL REFERENCES users(username), 
	token bytea NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS regtokens (
	issuer varchar(32) NOT NULL REFERENCES users(username),
	token bytea NOT NULL UNIQUE
);
