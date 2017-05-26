package models

import "gopkg.in/mgo.v2/bson"

// User is a struct that defines a user to add to the chat server
type User struct {
	Username string        `json:"username" bson:"username" binding:"required"`
	Email    string        `json:"email" bson:"email" binding:"required"`
	Password string        `json:"password" bson:"password" binding:"required"`
	ID       bson.ObjectId `json:"_id" bson:"_id"`
}

// UserError contains the fields for any user errors the struct that
type UserError struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// func (err *UserError) Error() string {

// }

// NewUser Returns an empty User struct
func NewUser() *User {
	return &User{}
}

// Validate validates input that user makes when entering username
func (u *User) Validate() (*User, bool) {

	return nil, true
}
