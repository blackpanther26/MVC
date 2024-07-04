package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/blackpanther26/mvc/pkg/utils"
	"github.com/blackpanther26/mvc/pkg/types"
	"github.com/gorilla/mux"
	"github.com/blackpanther26/mvc/pkg/models"
)

func ListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetAllBooks()
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Books": books,
	}

	utils.RenderTemplate(w, "clientPortal", data)
}

func CheckoutBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RenderTemplateWithMessage(w, "clientPortal", "Invalid book ID", "error")
		return
	}

	err = models.CheckoutBook(getUserIDFromContext(r.Context()), bookID)
	if err != nil {
		utils.RenderTemplateWithMessage(w, "clientPortal", err.Error(), "error")
		return
	}

	utils.RenderTemplateWithMessage(w, "clientPortal", "Book checkout request sent successfully.", "success")
}

func CheckinBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	err = models.CheckinBook(bookID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Book checked in successfully"}`))
}

func UserHistory(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())

	transactions, err := models.GetUserTransactions(userID)
	if err != nil {
		http.Error(w, "Failed to fetch user history", http.StatusInternalServerError)
		return
	}

	jsonTransactions, err := json.Marshal(transactions)
	if err != nil {
		http.Error(w, "Failed to marshal user history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonTransactions)
}

func getUserIDFromContext(ctx context.Context) uint {
	user := ctx.Value("user").(types.User)
	return user.ID
}