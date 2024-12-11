package types

import (
	"time"
)

// Some of the structs that will be used to interface between
// frontend, api, and db. Go has good handling for this

type User struct {
	Username     string    `json:"username" binding:"required"`
	Bio          string    `json:"bio"`
	Links        []string  `json:"links"`
	CreationDate time.Time `json:"created_on"`
	Picture      string    `json:"picture"`
}

type Project struct {
	ID           int64     `json:"id"`
	Owner        int64     `json:"owner" binding:"required"`
	Name         string    `json:"name" binding:"required"`
	Description  string    `json:"description" binding:"required"`
	Status       int16     `json:"status"`
	Likes        int64     `json:"likes"`
	Tags         []string  `json:"tags"`
	Links        []string  `json:"links"`
	CreationDate time.Time `json:"created_on"`
}

type Post struct {
	ID           int64     `json:"id"`
	User         int64     `json:"user" binding:"required"`
	Project      int64     `json:"project" binding:"required"`
	Likes        int64     `json:"likes"`
	Content      string    `json:"content" binding:"required"`
	Comments     []int64   `json:"comments"`
	CreationDate time.Time `json:"created_on"`
}

type Comment struct {
	ID            int64     `json:"id"`
	User          int64     `json:"user" binding:"required"`
	Post          int64     `json:"post" binding:"required"`
	ParentComment int64     `json:"parent_comment" binding:"required"`
	CreationDate  time.Time `json:"created_on"`
	Content       string    `json:"content" binding:"required"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}
