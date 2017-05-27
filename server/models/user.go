package models

import (
	"fmt"
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"

	"gopkg.in/mgo.v2/bson"
)

var emailRegex = regexp.MustCompile("\\A[a-zA-z\\d\\.]+@[a-zA-Z]{2,}\\.[a-z]{2,}\\z")
var passwordRegex = regexp.MustCompile("\\A[a-zA-z\\d\\.\\@\\#\\!\\?]{8,32}\\z")

// User is a struct that defines a user to add to the chat server
type User struct {
	Username             string        `json:"username" bson:"username" binding:"required"`
	Email                string        `json:"email" bson:"email" binding:"required"`
	Password             string        `json:"password" bson:"password" binding:"required"`
	PasswordConfirmation string        `json:"passwordConfirmation"`
	ID                   bson.ObjectId `json:"_id" bson:"_id"`
}

// UserError contains the fields for any user errors the struct that
type UserError struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Err      error
}

func (err *UserError) Error() string {
	return fmt.Sprintf("A User interaction error occured: %v %v %v %v", err.Username, err.Email, err.Password, err.Err)
}

// NewUser Returns an empty User struct
func NewUser() *User {
	return &User{}
}

// Validate validates the email and password received from the required fields from the user to make sure they are in the correct format
func (u *User) Validate() error {
	// log.Printf("%v", u)

	errorOccured := false
	errorObject := UserError{}

	if !emailRegex.MatchString(u.Email) {
		// TODO concatenate error object
		log.Printf("Email did not match. Given email: %v", u.Email)
		errorOccured = true
		errorObject.Email = u.Email + " is not formatted correctly."
	}

	if !passwordRegex.MatchString(u.Password) {
		// TODO concatenate error object
		log.Printf("Password did not match. Given password: %v", u.Password)
		errorOccured = true
		errorObject.Password = "Improper Formatting"
	}

	if errorOccured {
		return &errorObject
	}

	return nil
}

// Login validates input that user makes when logging in username
func (u *User) Login(ip string) error {
	errorObject := UserError{}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(ip)); err != nil {
		errorObject.Password = "User did not supply a matching password."
		return &errorObject
	}
	return nil
}

// Register validates the user input when a user supplies a not before seen email address (will implement with a register route later)
func (u *User) Register() (*UserError, bool) {
	log.Printf("%v", u)

	errorOccured := false
	errorObject := UserError{}

	if len(u.Username) < 3 {
		// TODO: Concatenate error object
		errorOccured = true
		errorObject.Username = fmt.Sprintf("Username was not long enough. Given Username: %v", u.Username)
		log.Printf("Username was not long enough. Given Username: %v", u.Username)
	}

	if !emailRegex.MatchString(u.Email) {
		// TODO concatenate error object
		errorOccured = true
		errorObject.Email = fmt.Sprintf("Email did not match correct format. Given email: %v", u.Email)
		log.Printf("Email did not match. Given email: %v", u.Email)
	}

	if !passwordRegex.MatchString(u.Password) {
		// TODO concatenate error object
		errorOccured = true
		errorObject.Password = fmt.Sprintf("Password did not match correct format. Given password: %v", u.Password)
		log.Printf("Password did not match correct format. Given password: %v", u.Password)
	}

	if errorOccured {
		return &errorObject, false
	}

	if u.Password != u.PasswordConfirmation {
		errorOccured = true
		errorObject.Password = fmt.Sprintf("Password and Password Confirmation did not match. Given password: %v Given Confirmation: %v", u.Password, u.PasswordConfirmation)
		log.Printf("Password and Password Confirmation did not match. Given password: %v Given Confirmation: %v", u.Password, u.PasswordConfirmation)

	}

	if errorOccured {
		return &errorObject, false
	}

	// No errors, so just hash the password and return
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		// handle error
		panic("Error occured in hashing password. User Email: " + u.Email)
	}
	log.Println(string(hash))

	u.Password = string(hash)

	return nil, true
}
