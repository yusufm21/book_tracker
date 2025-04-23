package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidBookInput(t *testing.T) {
	clear(books)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{"title": "Hamlet", "author": "Yusuf", "status": "reading"}`))
	w := httptest.NewRecorder()
	handleBooksPostRequest(w, req)
	res := w.Result()
	defer res.Body.Close()

	expected := http.StatusOK
	actual := res.StatusCode

	if expected != actual {
		t.Errorf("Expected %d, was %d", expected, actual)
	}

}

func TestInvalidBookInput(t *testing.T) {
	clear(books)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{"title": "Hamlet", "status": "reading"}`))
	w := httptest.NewRecorder()
	handleBooksPostRequest(w, req)
	res := w.Result()
	defer res.Body.Close()

	expected := http.StatusBadRequest
	actual := res.StatusCode

	if expected != actual {
		t.Errorf("Expected %d, was %d", expected, actual)
	}

}

func TestInvalidJson(t *testing.T) {
	clear(books)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{"title": "Incomplete", "author": "something", "status": "reading"`))
	w := httptest.NewRecorder()
	handleBooksPostRequest(w, req)
	res := w.Result()
	defer res.Body.Close()

	expected := http.StatusBadRequest
	actual := res.StatusCode

	if expected != actual {
		t.Errorf("Expected %d, was %d", expected, actual)
	}

}

func TestInvalidStatus(t *testing.T) {
	clear(books)
	fmt.Println(len(books))
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBufferString(`{"title": "Hamlet", "author": "Yusuf", "status": "borrowed"}`))
	w := httptest.NewRecorder()

	handleBooksPostRequest(w, req)
	res := w.Result()
	defer res.Body.Close()

	expected := http.StatusBadRequest
	actual := res.StatusCode

	if expected != actual {
		t.Errorf("Expected %d, was %d", expected, actual)
	}

	if len(books) > 0 {
		t.Fatalf("Expected 0 book in map, got %d", len(books))
	}

}

func TestPagination(t *testing.T) {
	// Clear books map and add test data
	clear(books)

	// Add multiple books for testing pagination
	books["1"] = Book{ID: "1", Title: "A Book", Author: "Author A", Status: "reading"}
	books["2"] = Book{ID: "2", Title: "B Book", Author: "Author B", Status: "unread"}
	books["3"] = Book{ID: "3", Title: "C Book", Author: "Author C", Status: "completed"}
	books["4"] = Book{ID: "4", Title: "D Book", Author: "Author D", Status: "reading"}
	books["5"] = Book{ID: "5", Title: "E Book", Author: "Author E", Status: "unread"}

	t.Run("No pagination parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		// Check that all books are returned (when no pagination is applied)
		var returnedBooks []Book
		if err := json.NewDecoder(res.Body).Decode(&returnedBooks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(returnedBooks) != 5 {
			t.Errorf("Expected 5 books, got %d", len(returnedBooks))
		}
	})

	t.Run("With limit parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books?limit=3", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		var returnedBooks []Book
		if err := json.NewDecoder(res.Body).Decode(&returnedBooks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(returnedBooks) != 3 {
			t.Errorf("Expected 3 books (limit), got %d", len(returnedBooks))
		}
	})

	t.Run("With offset parameter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books?offset=2", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		var returnedBooks []Book
		if err := json.NewDecoder(res.Body).Decode(&returnedBooks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(returnedBooks) != 3 {
			t.Errorf("Expected 3 books (after offset of 2), got %d", len(returnedBooks))
		}
	})

	t.Run("With both limit and offset parameters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books?limit=2&offset=1", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		var returnedBooks []Book
		if err := json.NewDecoder(res.Body).Decode(&returnedBooks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(returnedBooks) != 2 {
			t.Errorf("Expected 2 books (limit 2 after offset 1), got %d", len(returnedBooks))
		}
	})

	t.Run("Invalid limit - negative", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books?limit=-1", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d for negative limit, got %d", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("Invalid limit - too large", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books?limit=10", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d for too large limit, got %d", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("Invalid offset - negative", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books?offset=-1", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d for negative offset, got %d", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("Invalid offset - too large", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books?offset=10", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status code %d for too large offset, got %d", http.StatusBadRequest, res.StatusCode)
		}
	})

	t.Run("Pagination with status filter", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/books?status=reading&limit=1", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		var returnedBooks []Book
		if err := json.NewDecoder(res.Body).Decode(&returnedBooks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(returnedBooks) != 1 {
			t.Errorf("Expected 1 book (limit 1 with status filter), got %d", len(returnedBooks))
		}

		if returnedBooks[0].Status != "reading" {
			t.Errorf("Expected book with status 'reading', got '%s'", returnedBooks[0].Status)
		}
	})

	t.Run("Check correct books are returned with pagination", func(t *testing.T) {
		// Since books are sorted by title, we know the order should be A, B, C, D, E
		req := httptest.NewRequest(http.MethodGet, "/books?offset=1&limit=2", nil)
		w := httptest.NewRecorder()

		handleBooksGetRequest(w, req)

		res := w.Result()
		defer res.Body.Close()

		var returnedBooks []Book
		if err := json.NewDecoder(res.Body).Decode(&returnedBooks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if len(returnedBooks) != 2 {
			t.Fatalf("Expected 2 books, got %d", len(returnedBooks))
		}

		// Should be books B and C (second and third when sorted by title)
		if returnedBooks[0].Title != "B Book" || returnedBooks[1].Title != "C Book" {
			t.Errorf("Expected 'B Book' and 'C Book', got '%s' and '%s'",
				returnedBooks[0].Title, returnedBooks[1].Title)
		}
	})
}
