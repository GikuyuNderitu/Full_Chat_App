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

// NewUserController takes a Database connection to manage transactions
func NewUserController(db *mgo.Database) *UserController {
	return &UserController{db}
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
	_, ok := input.Validate()
	if !ok {
		// Implement handle error handling when Validate is complete
	}

	user := models.NewUser()
	err = uc.db.C("users").Find(bson.M{"email": input.Email}).One(&user)
	if err != nil {
		// Create User because it wasn't found
		input.ID = bson.NewObjectId()
		log.Printf("New to db %v", input)
		uc.db.C("users").Insert(input)
		utils.WriteJSON(true, input, w)
		return
	}

	log.Printf("Found from db %v", &user)

	log.Printf("Login By User: %v", input.Username)
	utils.WriteJSON(true, user, w)
}
