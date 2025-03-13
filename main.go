package main

import (
	"expense-tracker/database"
	"expense-tracker/handlers"
	"expense-tracker/middleware"
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	database.ConnectDB()
	defer database.CloseDB()

	router := mux.NewRouter()
	router.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleWare)

	api.HandleFunc("/expenses", handlers.AddExpense).Methods("POST")
	api.HandleFunc("/expenses", handlers.GetExpenses).Methods("GET")
	api.HandleFunc("/expenses/category", handlers.GetExpensesByCategory).Methods("GET")
	api.HandleFunc("/expenses/{id:[0-9]+}", handlers.UpdateExpense).Methods("PUT")
	api.HandleFunc("/expenses/{id:[0-9]+}", handlers.DeleteExpense).Methods("DELETE")

	fmt.Println("Server started on :8080")
	log.Println(http.ListenAndServe(":8080", router))
}