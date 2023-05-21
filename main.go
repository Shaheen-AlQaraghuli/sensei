package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	router.Patch("/user/{id}", updateUser)
	router.Delete("/user/{id}", deleteUser)

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

type UpdateUserRequest struct {
	Password string `json:"password"`
}

type LogDetail struct {
	Key string
	Value any
}

func getUser(w http.ResponseWriter, r *http.Request){
	for _, user := range users {
		if user.ID == chi.URLParam(r, "id"){
			respond(
				w,
				http.StatusOK,
				GetUserResponse{
					ID: user.ID,
					Name: user.Name,
				},
			)
			return
		}
	}
	respond(
		w,
		http.StatusNotFound,
		GetUserResponse{
			Error: &ErrorResponse{
				Message: "User not found",
				Code: "user_not_found",
			},
		},
	)
	return
}

func createUser(w http.ResponseWriter, r *http.Request){
	var userReq CreateUserRequest
	
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		logError(err, LogDetail{Key: "message", Value: "decoder failed"})
		respond(
			w,
			http.StatusBadRequest,
			CreateUserResponse{
				Error: &ErrorResponse{
					Message: "Something unexpected happened. Please try again",
					Code: "unexpected_error",
				},
			},
		)
		return
	}

	if userReq.Name == "" || userReq.Password == "" {
		respond(
			w,
			http.StatusBadRequest,
			CreateUserResponse{
				Error: &ErrorResponse{
					Message: "Please enter valid user details",
					Code: "invalid_input",
				},
			},
		)
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

	respond(
		w, 
		http.StatusCreated,
		CreateUserResponse{
			ID: userID,
		}, 
	)
}

func deleteUser(w http.ResponseWriter, r *http.Request){
	userID := chi.URLParam(r, "id")
	for i, user := range users {
		if user.ID == userID {
			userMutex.Lock()
			users = append(users[:i], users[i+1:]...)
			userMutex.Unlock()

			respond(
				w,
				http.StatusOK,
				nil,
			)
			return
		}
	}
	respond(
		w,
		http.StatusNotFound,
		ErrorResponse{
			Message: "User not found",
			Code: "user_not_found",
		},
	)
	return
}

func updateUser(w http.ResponseWriter, r *http.Request){
	var userDetails UpdateUserRequest
	err := json.NewDecoder(r.Body).Decode(&userDetails)
	if err != nil {
		logError(err, LogDetail{Key: "message", Value: "decoder failed"})
		respond(
			w,
			http.StatusBadRequest,
			ErrorResponse{
				Message: "Something unexpected happened. Please try again",
				Code: "unexpected_error",
			},
		)
	}

	if userDetails.Password == "" {
		respond(
			w,
			http.StatusBadRequest,
			ErrorResponse{
				Message: "Please enter valid user details",
				Code: "invalid_input",
			},
		)
		return
	}

	userID := chi.URLParam(r, "id")
	for i, user := range users {
		if user.ID == userID {
			userMutex.Lock()
			users[i].Password = userDetails.Password
			userMutex.Unlock()

			respond(
				w,
				http.StatusOK,
				nil,
			)
			return
		}
	}

	respond(
		w,
		http.StatusNotFound,
		ErrorResponse{
			Message: "User not found",
			Code: "user_not_found",
		},
	)
	return
}

func respond(w http.ResponseWriter, status int, resp any){
	logResponse(status, resp)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

func logResponse(status int, resp any){
	//todo: add requestID
	respJSON, _ := json.Marshal(resp)
	log.Printf("Status: %d - Response: %+v", status, string(respJSON))
}

func logError(err error, lds ...LogDetail){
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Error: %+v", err))

	for _, ld := range lds {
		sb.WriteString(fmt.Sprintf("\n%s: %+v", ld.Key, ld.Value))
	}

	log.Print(sb.String())
}

//decode function concept
/* func decodeRequest(w http.ResponseWriter, r *http.Request, v any) bool{
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		respond(
			w,
			http.StatusBadRequest,
			CreateUserResponse{
				Error: &ErrorResponse{
					Message: "Something unexpected happened. Please try again",
					Code: "unexpected_error",
				},
			},
		)
		return false
	}
	return true
} */
