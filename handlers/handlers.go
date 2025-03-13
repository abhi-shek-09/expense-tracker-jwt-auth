package handlers

import (
	"encoding/json"
	"expense-tracker/database"
	"expense-tracker/middleware"
	"expense-tracker/models"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

func AddExpense(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uID, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid User ID", http.StatusUnauthorized)
		return
	}

	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	var expenseID int
	query := "INSERT INTO expenses (user_id, amount, category, description, date) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	err := database.DB.QueryRow(query, uID, expense.Amount, expense.Category, expense.Description, expense.Date).Scan(&expenseID)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int {"id" : expenseID})
}

func GetExpenses(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uID, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid User ID", http.StatusUnauthorized)
		return
	}

	query := "SELECT id, amount, category, description, date FROM expenses WHERE user_id=$1 ORDER BY date DESC;"
	rows, err := database.DB.Query(query, uID)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next(){
		var exp models.Expense
		if err := rows.Scan(&exp.ID, &exp.Amount, &exp.Category, &exp.Description, &exp.Date); err != nil {
			http.Error(w, "Error scanning expenses", http.StatusInternalServerError)
            return
		}
		expenses = append(expenses, exp)
	}

	if err := rows.Err(); err != nil {
        http.Error(w, "Error reading expenses", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func GetExpensesByCategory(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uID, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid User ID", http.StatusUnauthorized)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
        http.Error(w, "Category is required", http.StatusBadRequest)
        return
    }

	query := "SELECT id, amount, category, description, date FROM expenses WHERE user_id=$1 AND category=$2 ORDER BY date DESC;"
	rows, err := database.DB.Query(query, uID, category)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next(){
		var exp models.Expense
		if err := rows.Scan(&exp.ID, &exp.Amount, &exp.Category, &exp.Description, &exp.Date); err != nil {
			http.Error(w, "Error scanning expenses", http.StatusInternalServerError)
            return
		}
		expenses = append(expenses, exp)
	}

	if err := rows.Err(); err != nil {
        http.Error(w, "Error reading expenses", http.StatusInternalServerError)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)
}

func UpdateExpense(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uID, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid User ID", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	expenseID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid expense ID", http.StatusBadRequest)
        return
    }

	var updatedExpense models.Expense
    if err := json.NewDecoder(r.Body).Decode(&updatedExpense); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

	query := "UPDATE expenses SET amount=$1, category=$2, description=$3, date=$4 WHERE id=$5 AND user_id=$6 RETURNING id;"
    err = database.DB.QueryRow(query, updatedExpense.Amount, updatedExpense.Category, updatedExpense.Description, updatedExpense.Date, expenseID, uID).Scan(&expenseID)
    if err != nil {
        http.Error(w, "Expense not found or unauthorized", http.StatusForbidden)
        return
    }
	updatedExpense.ID = expenseID
    updatedExpense.UserID = uID
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedExpense)
}

func DeleteExpense(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value(middleware.UserIDKey)
	if userID == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	uID, ok := userID.(int)
	if !ok {
		http.Error(w, "Invalid User ID", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	expenseID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid expense ID", http.StatusBadRequest)
        return
    }

	query := "DELETE FROM expenses WHERE id=$1 AND user_id=$2"
    result, err := database.DB.Exec(query, expenseID, uID)
    if err != nil {
        http.Error(w, "Failed to delete expense", http.StatusInternalServerError)
        return
    }

	rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "Expense not found or unauthorized", http.StatusForbidden)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}