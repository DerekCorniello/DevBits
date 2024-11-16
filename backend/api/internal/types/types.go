package types

import "time"

// Some of the structs that will be used to interface between
// frontend, api, and db. Go has good handling for this

type User struct {
	Username     string    `json:"username"`
	Bio          string    `json:"bio"`
	Links        []string  `json:"links"`
	CreationDate time.Time `json:"created_on"`
	// TODO: handle profile pic, followers, following
}

type Project struct {
	ID           int64     `json:"id"`
	Owner        int64     `json:"owner"` // linked with a User
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Status       int16     `json:"status"` // TODO: Enum? Does Go have Enums?
	Likes        int64     `json:"likes"`
	Tags         []string  `json:"tags"`
	Links        []string  `json:"links"`
	CreationDate time.Time `json:"created_on"`
	// TODO: handle followers, posts
}

type Post struct {
	ID           int64     `json:"id"`
	User         int64     `json:"user"`    // Linked to User
	Project      int64     `json:"project"` // Linked to Project
	Likes        int64     `json:"likes"`
	Content      int64     `json:"content"`
	Comments     int64     `json:"comment"`
	CreationDate time.Time `json:"created_on"`
}

type Comment struct {
	ID            int64     `json:"id"`
	User          int64     `json:"user"`           // Linked to User
	Post          int64     `json:"post"`           // linked to Post
	ParentComment int64     `json:"parent_comment"` // TODO: does this need a Nil?
	CreationDate  time.Time `json:"created_on"`
}
