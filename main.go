package main

import (
	"net/http"
	"encoding/json"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type User struct {
	ID 		string
	Name	string
	Password string
}

var users []User = []User{}

func main(){
	router := chi.NewRouter()

	router.Get("/user/{id}", getUser)
	router.Post("/user", createUser)

	http.ListenAndServe(":3000", router)
}

type GetUserResponse struct {
	ID   string    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID           string    `json:"id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
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
		ErrorMessage: "User not found",
	})
	return
}

func createUser(w http.ResponseWriter, r *http.Request){
	var userReq CreateUserRequest
	
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateUserResponse{
			ErrorMessage: "Please enter valid user details",
		})
		return
	}

	if userReq.Name == "" || userReq.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(CreateUserResponse{
			ErrorMessage: "Please enter valid user details",
		})
		return
	}

	userID := uuid.NewString()
	users = append(users, User{
		ID: userID,
		Name: userReq.Name,
		Password: userReq.Password,
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CreateUserResponse{
		ID: userID,
	})
}
