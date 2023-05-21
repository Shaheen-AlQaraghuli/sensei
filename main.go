package main

import (
	"net/http"
	"encoding/json"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type User struct {
	ID 		string
	Name	string
	Password string
}

var users []User = []User{}
var userMutex sync.Mutex

func main(){
	router := chi.NewRouter()

	router.Get("/user/{id}", getUser)
	router.Post("/user", createUser)

	http.ListenAndServe(":3000", router)
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code string `json:"code"`
}

type GetUserResponse struct {
	ID   string    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Error *ErrorResponse `json:"error"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID           string    `json:"id,omitempty"`
	Error *ErrorResponse `json:"error"`
}

func getUser(w http.ResponseWriter, r *http.Request){
	for _, user := range users {
		if user.ID == chi.URLParam(r, "id"){
			json.NewEncoder(w).Encode(GetUserResponse{
				ID: user.ID,
				Name: user.Name,
			})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(GetUserResponse{
		Error: &ErrorResponse{
			Message: "User not found",
			Code: "user_not_found",
		},
	})
	return
}

func createUser(w http.ResponseWriter, r *http.Request){
	var userReq CreateUserRequest
	
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateUserResponse{
			Error: &ErrorResponse{
				Message: "Something unexpected happened. Please try again",
				Code: "unexpected_error",
			},
		})
		return
	}

	if userReq.Name == "" || userReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateUserResponse{
			Error: &ErrorResponse{
				Message: "Please enter valid user details",
				Code: "invalid_input",
			},
		})
		return
	}

	userID := uuid.NewString()
	userMutex.Lock()
	users = append(users, User{
		ID: userID,
		Name: userReq.Name,
		Password: userReq.Password,
	})
	userMutex.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateUserResponse{
		ID: userID,
	})
}
