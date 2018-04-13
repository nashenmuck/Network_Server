CREATE TABLE IF NOT EXISTS followers (
	followeeId int NOT NULL,
	followerId int NOT NULL,
	followedwhen timestamp NOT NULL
);
