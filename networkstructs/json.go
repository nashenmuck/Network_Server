package networkstructs

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Servers struct {
	Server      string `json:"server"`
	Blacklisted bool   `json:"blacklisted"`
	Whitelisted bool   `json:"whitelisted"`
}

type Users struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
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

func (p *Posts) Decode(w http.ResponseWriter, r *http.Request) error {
	var err error
	if r.Body == nil {
		http.Error(w, "No body found", 400)
		err = fmt.Errorf("No body found")
		return err
	}
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	return err
}

func (u *Users) Decode(w http.ResponseWriter, r *http.Request) error {
	var err error
	if r.Body == nil {
		http.Error(w, "No body found", 400)
		err = fmt.Errorf("No body found")
		return err
	}
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	return err
}

func (f *Followers) Decode(w http.ResponseWriter, r *http.Request) error {
	var err error
	if r.Body == nil {
		http.Error(w, "No body found", 400)
		err = fmt.Errorf("No body found")
		return err
	}
	err = json.NewDecoder(r.Body).Decode(&f)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return err
	}
	return err
}
