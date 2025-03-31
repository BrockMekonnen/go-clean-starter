// package middleware

// import (
// 	"net/http"
// )

// // AsyncHandler represents a function that handles HTTP requests and returns an error if something goes wrong.
// type AsyncHandler func(w http.ResponseWriter, r *http.Request) error

// // RunAsync wraps an AsyncHandler and ensures that errors are properly handled.
// func RunAsync(handler AsyncHandler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if err := handler(w, r); err != nil {
// 			// Handle error (this could be logging, setting a status code, etc.)
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 		}
// 	}
// }