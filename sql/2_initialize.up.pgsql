CREATE TABLE IF NOT EXISTS servers (
	server varchar NOT NULL UNIQUE,
	blacklisted boolean DEFAULT FALSE,
	whitelisted boolean DEFAULT FALSE
);
CREATE TABLE IF NOT EXISTS users (
	username varchar(32) NOT NULL UNIQUE PRIMARY KEY,
	password bytea NOT NULL,
	isadmin boolean DEFAULT FALSE,
	parentAdmin varchar(32),
	canonical_user varchar(32) NOT NULL,
	canonical_server varchar NOT NULL REFERENCES servers(server)
);
CREATE SEQUENCE group_id_seq;
CREATE TABLE IF NOT EXISTS groups (
	group_id int NOT NULL UNIQUE DEFAULT nextval('group_id_seq'),
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
	followedwhen timestamp NOT NULL DEFAULT NOW(),
    UNIQUE(followeeId, followee_server, followerId, follower_server)
);
CREATE TABLE IF NOT EXISTS group_followers (
	group_id int NOT NULL REFERENCES groups(group_id),
	follower varchar(32) NOT NULL REFERENCES users(username),
	follower_server varchar references servers(server),
	UNIQUE (group_id, follower)
);
CREATE SEQUENCE post_id_seq;
CREATE TABLE IF NOT EXISTS posts (
	id int primary key default nextval('post_id_seq'),
	username varchar(32) NOT NULL REFERENCES users(username),
	body text NOT NULL,
	groupid int REFERENCES groups(group_id),
	special_groupid int,
	is_special_group boolean NOT NULL,
	origin_server varchar NOT NULL REFERENCES servers(server),
	date timestamp NOT NULL DEFAULT NOW()
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
CREATE INDEX IF NOT EXISTS user_index ON users (username);
CREATE INDEX IF NOT EXISTS groups_index ON groups (group_id,owner);
CREATE INDEX IF NOT EXISTS posts_user_index ON posts (username);
CREATE INDEX IF NOT EXISTS posts_time_index ON posts (date);
