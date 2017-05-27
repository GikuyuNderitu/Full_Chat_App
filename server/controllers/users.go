package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GikuyuNderitu/chat_application/server/models"
	"github.com/GikuyuNderitu/chat_application/server/utils"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// UserController is the struct that holds the context for the interaction of the User Controller
type UserController struct {
	db *mgo.Database
}

type returnUser struct {
	Username string        `json:"username"`
	ID       bson.ObjectId `json:"_id"`
}

// NewUserController takes a Database connection to manage transactions
func NewUserController(db *mgo.Database) *UserController {
	return &UserController{db}
}

// Register is a method that gets a request to register a user, validates the supplied data and returns an error or a sanitized user object accordingly
func (uc UserController) Register(w http.ResponseWriter, r *http.Request) {
	input := models.NewUser()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		log.Fatalf("Uh Oh. Error Decoding the JSON %v\n", err)
	}
	defer r.Body.Close()

	user := models.NewUser()
	err = uc.db.C("users").Find(bson.M{"email": input.Email}).One(&user)
	if err != nil {
		// Create User because it wasn't found
		input.ID = bson.NewObjectId()
		input.Register()
		log.Printf("New to db %v", input)
		uc.db.C("users").Insert(input)
		utils.WriteJSON(true, sanitizeUser(input), w)
		return
	}

	log.Printf("Can't register. User with email '%v' already found in db.", user.Email)

	utils.WriteJSON(false, models.UserError{Email: "Email already exists"}, w)
}

// Login is a method of UserController that is a Handler Func It logs the user in
func (uc UserController) Login(w http.ResponseWriter, r *http.Request) {
	input := models.NewUser()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		log.Fatalf("Uh Oh. Error Decoding the JSON %v\n", err)
	}
	defer r.Body.Close()

	// Validate input
	err = input.Validate()
	if err != nil {
		// Implement handle error handling when Validate is complete
		utils.WriteJSON(false, err, w)
		return
	}

	user := models.NewUser()
	err = uc.db.C("users").Find(bson.M{"email": input.Email}).One(&user)
	if err != nil {
		// Create User because it wasn't found
		input.ID = bson.NewObjectId()
		input.Register()
		log.Printf("New to db %v", input)
		uc.db.C("users").Insert(input)
		utils.WriteJSON(true, sanitizeUser(input), w)
		return
	}
	log.Printf("Logging in user %v", user)
	if err = user.Login(input.Password); err != nil {
		log.Printf("User provided incorrect password")
		utils.WriteJSON(false, err, w)
		return
	}
	utils.WriteJSON(true, sanitizeUser(user), w)
}

// TODO: Add Santize method to remove Password Field
func sanitizeUser(u *models.User) *returnUser {
	return &returnUser{u.Username, u.ID}
}
