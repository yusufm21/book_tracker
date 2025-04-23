package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var books = make(map[string]Book)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Status string `json:"status"` // allowed: "unread", "reading", "completed"
}

func handleRootrequest(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to my website"))
}

func handleBooksPostRequest(w http.ResponseWriter, r *http.Request) {
	var book Book

	//Prevents incorrect json from user
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !validatePostRequest(book, w) {
		return
	}
	book.ID = uuid.New().String() //Generate a unique key for every book.
	books[book.ID] = book         //The keys for books in the map will be their unique uuid

}

func validatePostRequest(b Book, w http.ResponseWriter) bool {
	// Only allow certain statuses
	if !(b.Status == "unread" || b.Status == "reading" || b.Status == "completed") {
		http.Error(w, "Only the condtions of unread, reading and completed is allowed.", http.StatusBadRequest)
		return false
	} else if b.Author == "" {
		http.Error(w, "You have to provide the name of an author.", http.StatusBadRequest)
		return false
	} else if b.Title == "" {
		http.Error(w, "You have to provide a title.", http.StatusBadRequest)
		return false
	} else if b.ID != "" {
		// Prevents the user from adding an ID
		http.Error(w, "Only provide the title, author and status", http.StatusBadRequest)
		return false
	}

	return true

}

func paginationHandler(w http.ResponseWriter, r *http.Request, bookList []Book) []Book {
	limitString := r.URL.Query().Get("limit")
	offsetString := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitString)
	offset, _ := strconv.Atoi(offsetString)

	// Apply offest filter
	if offsetString != "" {
		if offset < len(bookList) && offset > 0 {
			bookList = bookList[offset:]
		} else {
			http.Error(w, "The offset has to be less than the amount of books and/or a postive integer", http.StatusBadRequest)
		}
	}

	// Apply limit filter
	if limitString != "" {
		if limit < len(bookList) && limit > 0 {
			bookList = bookList[:limit]
		} else {
			http.Error(w, "The limit has to be less than the amount of books and/or a postive integer", http.StatusBadRequest)
		}
	}

	return bookList

}

func handleBooksGetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	filter := r.URL.Query().Get("status")

	// Collect all books from map into a slice
	bookList := slices.Collect(maps.Values(books))

	// Apply status filter if specified
	if filter != "" {
		var listCopy []Book
		for _, value := range bookList {
			if value.Status == filter {
				listCopy = append(listCopy, value)
			}
		}
		bookList = listCopy
	}

	// Sort books alphabetically by title
	sort.Slice(bookList, func(i, j int) bool {
		return bookList[i].Title < bookList[j].Title
	})

	// paginationHandler will return a filtered bookList
	bookList = paginationHandler(w, r, bookList)
	json.NewEncoder(w).Encode(bookList)
}

func handleBooksIdDeleteRequest(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/books/")

	// Check if the book exists before deleting
	_, ok := books[id]
	if !ok {
		http.Error(w, "The book does not exist", http.StatusBadRequest)
		return
	}
	delete(books, id)
}

func handleBooksIdPutRequest(w http.ResponseWriter, r *http.Request) {
	var bookInput Book
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/books/")

	// Check if the book exists before updating
	book, ok := books[id]
	if !ok {
		http.Error(w, "The book does not exist", http.StatusBadRequest)
		return
	}

	//Prevents incorrect json from user
	err := json.NewDecoder(r.Body).Decode(&bookInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !validatePutRequest(bookInput, w) {
		return
	}

	// Only update the book's status
	book.Status = bookInput.Status
	books[id] = book
}

func validatePutRequest(b Book, w http.ResponseWriter) bool {
	// Only status is allowed for PUT updates — all other fields should be empty
	if !(b.Status == "unread" || b.Status == "reading" || b.Status == "completed") {
		http.Error(w, "Only the condtions of unread, reading and completed is allowed.", http.StatusBadRequest)
		return false
	} else if b.Author != "" {
		http.Error(w, "Only provide the status of the book", http.StatusBadRequest)
		return false
	} else if b.Title != "" {
		http.Error(w, "Only provide the status of the book", http.StatusBadRequest)
		return false
	} else if b.ID != "" {
		http.Error(w, "Only provide the status of the book", http.StatusBadRequest)
		return false
	}
	return true

}

func main() {
	http.HandleFunc("/", handleRootrequest)
	http.HandleFunc("POST /books", handleBooksPostRequest)
	http.HandleFunc("GET /books", handleBooksGetRequest)
	http.HandleFunc("PUT /books/", handleBooksIdPutRequest)
	http.HandleFunc("DELETE /books/", handleBooksIdDeleteRequest)
	fmt.Println("Listening to port 8080")
	http.ListenAndServe("localhost:8080", nil)
}
