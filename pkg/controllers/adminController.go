package controllers

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/blackpanther26/mvc/pkg/models"
	"github.com/blackpanther26/mvc/pkg/views"
	"github.com/blackpanther26/mvc/pkg/types"
)

func AdminListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := models.GetAllBooks()
	if err != nil {
		http.Error(w, "Failed to fetch books", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Books": books,
	}

	views.RenderTemplate(w, "adminListBooks", data)
}

func AdminAddBook(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        views.RenderTemplate(w, "adminAddBook", nil)
        return
    }

    title := r.FormValue("title")
    author := r.FormValue("author")
    isbn := r.FormValue("isbn")
    totalCopiesStr := r.FormValue("total_copies")

    totalCopies, err := strconv.Atoi(totalCopiesStr)
    if err != nil {
        views.RenderTemplateWithMessage(w, "adminAddBook", "Invalid total copies", "error")
        return
    }

    if len(isbn) != 13 {
        views.RenderTemplateWithMessage(w, "adminAddBook", "ISBN must be 13 characters long", "error")
        return
    }

    if isbn[0] == '-' {
        views.RenderTemplateWithMessage(w, "adminAddBook", "ISBN cannot be negative", "error")
        return
    }

    book := types.Book{
        Title:       title,
        Author:      author,
        ISBN:        isbn,
        TotalCopies: totalCopies,
    }

    err = models.AddBook(&book)
    if err != nil {
        views.RenderTemplateWithMessage(w, "adminAddBook", "Failed to add book", "error")
        return
    }

    views.RenderTemplateWithMessage(w, "adminAddBook", "Book added successfully", "success")
}

func AdminEditBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		book, err := models.GetBookByID(bookID)
		if err != nil {
			http.Error(w, "Failed to fetch book", http.StatusInternalServerError)
			return
		}

		data := map[string]interface{}{
			"Book": book,
		}

		views.RenderTemplate(w, "adminEditBook", data)
		return
	}

	totalCopies, err := strconv.Atoi(r.FormValue("totalCopies"))
	if err != nil {
		http.Error(w, "Invalid total copies value", http.StatusBadRequest)
		return
	}

	book := types.Book{
		ID:          uint(bookID),
		Title:       r.FormValue("title"),
		Author:      r.FormValue("author"),
		ISBN:        r.FormValue("isbn"),
		TotalCopies: totalCopies,
	}

	err = models.UpdateBook(&book)
	if err != nil {
		views.RenderTemplateWithMessage(w, "adminEditBook", "Failed to update book", "error")
		return
	}

	views.RenderTemplateWithMessage(w, "adminEditBook", "Book updated successfully", "success")
}

func AdminDeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	err = models.DeleteBook(uint(bookID))
	if err != nil {
		views.RenderTemplateWithMessage(w, "adminListBooks", "Failed to delete book", "error")
		return
	}

	views.RenderTemplateWithMessage(w, "adminListBooks", "Book deleted successfully", "success")
}

func AdminListTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := models.GetAllTransactions()
	if err != nil {
		http.Error(w, "Failed to fetch transactions", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"Transactions": transactions,
	}
	views.RenderTemplate(w, "adminListTransactions", data)
}

func AdminApproveTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	err = models.UpdateTransactionStatus(uint(transactionID), "approved")
	if err != nil {
		views.RenderTemplateWithMessage(w, "adminListTransactions", "Failed to approve transaction", "error")
		return
	}

	views.RenderTemplateWithMessage(w, "adminListTransactions", "Transaction approved successfully", "success")
}

func AdminDenyTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid transaction ID", http.StatusBadRequest)
		return
	}

	err = models.UpdateTransactionStatus(uint(transactionID), "denied")
	if err != nil {
		views.RenderTemplateWithMessage(w, "adminListTransactions", "Failed to deny transaction", "error")
		return
	}

	views.RenderTemplateWithMessage(w, "adminListTransactions", "Transaction denied successfully", "success")
}

func AdminListAdminRequests(w http.ResponseWriter, r *http.Request) {
	adminRequests, err := models.GetAllAdminRequests()
	if err != nil {
		http.Error(w, "Failed to fetch admin requests", http.StatusInternalServerError)
		return
	}
	data := map[string]interface{}{
		"AdminRequests": adminRequests,
	}
	views.RenderTemplate(w, "adminListAdminRequests", data)
}

func AdminApproveAdminRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	err = models.UpdateAdminRequestStatus(uint(requestID), "approved")
	if err != nil {
		views.RenderTemplateWithMessage(w, "adminListAdminRequests", "Failed to approve admin request", "error")
		return
	}

	views.RenderTemplateWithMessage(w, "adminListAdminRequests", "Admin request approved successfully", "success")
}

func AdminDenyAdminRequest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid request ID", http.StatusBadRequest)
		return
	}

	err = models.UpdateAdminRequestStatus(uint(requestID), "denied")
	if err != nil {
		views.RenderTemplateWithMessage(w, "adminListAdminRequests", "Failed to deny admin request", "error")
		return
	}

	views.RenderTemplateWithMessage(w, "adminListAdminRequests", "Admin request denied successfully", "success")
}