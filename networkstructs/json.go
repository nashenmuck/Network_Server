package networkstructs

import (
	"time"
)

type Servers struct {
	Server      string `json:"server"`
	Blacklisted bool   `json:"blacklisted"`
	Whitelisted bool   `json:"whitelisted"`
}

type Users struct {
	Username        string `json:"username"`
	Password        []byte `json:"password"`
	IsAdmin         bool   `json:"isadmin"`
	ParentAdmin     string `json:"parentadmin"`
	CanonicalUser   string `json:"canonical_user"`
	CanonicalServer string `json:"canonical_server"`
}

type Groups struct {
	GroupId   int    `json:"group_id"`
	Owner     string `json:"owner"`
	GroupName string `json:"group_name"`
}

type Followers struct {
	FolloweeId     string    `json:"followeeid"`
	FolloweeServer string    `json:"followee_server"`
	FollowerId     string    `json:"followerid"`
	FollowerServer string    `json:"follower_server"`
	FollowedWhen   time.Time `json:"followedwhen"`
}

type GroupFollowers struct {
	GroupId        int    `json:"group_id"`
	Follower       string `json:"follower"`
	FollowerServer string `json:"follower_server"`
}

type Posts struct {
	Username       string    `json:"username"`
	Body           string    `json:"body"`
	GroupId        int       `json:"group_id"`
	IsSpecialGroup bool      `json:"is_special_group"`
	OriginServer   string    `json:"origin_server"`
	Date           time.Time `json:"date"`
}
