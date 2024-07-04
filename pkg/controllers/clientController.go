package controllers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/blackpanther26/mvc/pkg/models"
	"github.com/blackpanther26/mvc/pkg/types"
	"github.com/blackpanther26/mvc/pkg/views"
	"github.com/gorilla/mux"
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

	views.RenderTemplate(w, "clientPortal", data)
}

func CheckoutBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		views.RenderTemplateWithMessage(w, "clientPortal", "Invalid book ID", "error")
		return
	}

	err = models.CheckoutBook(getUserIDFromContext(r.Context()), bookID)
	if err != nil {
		views.RenderTemplateWithMessage(w, "clientPortal", err.Error(), "error")
		return
	}

	views.RenderTemplateWithMessage(w, "clientPortal", "Book checkout request sent successfully.", "success")
}

func CheckinBook(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    bookID, err := strconv.Atoi(vars["id"])
    if err != nil {
        views.RenderTemplateWithMessage(w, "clientPortal", "Invalid book ID", "error")
        return
    }

    err = models.CheckinBook(bookID)
    if err != nil {
        views.RenderTemplateWithMessage(w, "clientPortal", err.Error(), "error")
        return
    }

    views.RenderTemplateWithMessage(w, "clientPortal", "Book checked in successfully", "success")
}

func UserHistory(w http.ResponseWriter, r *http.Request) {
    userID := getUserIDFromContext(r.Context())

    transactions, err := models.GetUserTransactions(userID)
    if err != nil {
        http.Error(w, "Failed to fetch user history", http.StatusInternalServerError)
        return
    }

    data := map[string]interface{}{
        "History": transactions,
    }
	fmt.Println(data)
    views.RenderTemplate(w, "userHistory", data)
}

func getUserIDFromContext(ctx context.Context) uint {
	user := ctx.Value("user").(types.User)
	return user.ID
}

func RequestAdmin(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r.Context())

	err := models.SendAdminRequest(userID)
	if err != nil {
		views.RenderTemplateWithMessage(w, "clientPortal", err.Error(), "error")
		return
	}

	views.RenderTemplateWithMessage(w, "clientPortal", "Admin request sent successfully.", "success")
}

func SearchBooks(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("search")
	if query == "" {
		views.RenderTemplateWithMessage(w, "clientPortal", "Please enter a search term", "error")
		return
	}

	books, err := models.SearchBooks(query)
	if err != nil {
		views.RenderTemplateWithMessage(w, "clientPortal", "Failed to search books", "error")
		return
	}

	data := map[string]interface{}{
		"Books": books,
	}

	views.RenderTemplate(w, "clientPortal", data)
}