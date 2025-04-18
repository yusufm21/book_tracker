# Book Tracker

Book Tracker is small HTTP API in Go to manage a book tracker application.

Here are the available enpoints and their descriptions:

### HTTP Server & API Endpoints


| Method | Endpoint       | Description                             |
|--------|----------------|-----------------------------------------|
| POST   | `/books`       | Add a new book                          |
| GET    | `/books`       | Retrieve a list of all books            |
| PUT    | `/books/{id}`  | Update the status of a specific book    |
| DELETE | `/books/{id}`  | Delete a book by its ID                 |


## 🚀 Setup instructions

### 1. Clone the Repository

```bash
git clone https://github.com/your-username/your-repo-name.git
cd your-repo-name
```

### 2. Install Dependencies
Make sure to be in the root directory of the project
```bash
go get github.com/google/uuid
```
This installs the UUID package used by the application.

## Usage

Navigate to the src directory:
```bash
cd src
```
Then start the server with
```bash
go run main.go
```

## Example

Use curl to add a new book
```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title": "kebenkaise", "author": "karl gustav", "status": "unread"}'
```
## Testing

Run the tests from the src directory
```bash
go test
```

## Key design decisions

### 1. Route handling

Originally, the API had only two handler functions for `/books` and `/books/`. The functions started to get big and clunky To solve this, I refactored the code to use four separate handler functions for each HTTP method and path. 

Luckily I found out in the documentation that `ListenAndServe()` has a default httphandler `deafultservemux`that has in-built pattern matching reducing the need to do manual router handling.

## 2. JSON validation

The validation logic for the JSON was seperated into its own functions to keep each function small, focused and easier to read. This also makes testing bug finding easier.

## 3. UUID as the key for the map

I used the `UUID`as the keys for the map in order to mimic ISBN's of real books. This ensures that two books with the same name and same author can exist. I later realized from a user perspective this probably wasn't the best decision. It's probably more likely for a user to add the same book twice to their tracker, rather than the likelihood of them wanting to add two books that has the same author and title hahaha.
