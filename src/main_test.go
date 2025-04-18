package main

import (
	"bytes"
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
