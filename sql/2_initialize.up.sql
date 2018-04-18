CREATE TABLE IF NOT EXISTS servers (
	server varchar NOT NULL UNIQUE,
	blacklisted boolean DEFAULT FALSE,
	whitelisted boolean DEFAULT FALSE
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
CREATE TABLE IF NOT EXISTS groups (
	group_id int NOT NULL UNIQUE,
	owner varchar(32) NOT NULL REFERENCES users(username),
	group_name varchar(32) NOT NULL,
	PRIMARY KEY (group_id, owner),
	UNIQUE (owner,group_name)
);
CREATE TABLE IF NOT EXISTS followers (
	followeeId varchar(32) NOT NULL REFERENCES users(username),
	followee_server varchar NOT NULL REFERENCES servers(server),
	followerId varchar(32) NOT NULL REFERENCES users(username), 
	follower_server varchar NOT NULL REFERENCES servers(server),
	followedwhen timestamp NOT NULL
);
CREATE TABLE IF NOT EXISTS group_followers (
	group_id int NOT NULL REFERENCES groups(group_id),
	follower varchar(32) NOT NULL REFERENCES users(username),
	follower_server varchar references servers(server),
	UNIQUE (group_id, follower)
);
CREATE TABLE IF NOT EXISTS posts (
	username varchar(32) NOT NULL REFERENCES users(username),
	body text NOT NULL,
	groupid int NOT NULL REFERENCES groups(group_id),
	is_special_group boolean NOT NULL,
	origin_server varchar NOT NULL REFERENCES servers(server)
);
CREATE TABLE IF NOT EXISTS authtokens (
	username varchar(32) NOT NULL REFERENCES users(username), 
	token bytea NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS regtokens (
	issuer varchar(32) NOT NULL REFERENCES users(username),
	token bytea NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS audit (
	data timestamp NOT NULL,
	event_type varchar,
	event_message varchar
);
