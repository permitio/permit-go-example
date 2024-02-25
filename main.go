package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/permitio/permit-golang/pkg/config"
	"github.com/permitio/permit-golang/pkg/enforcement"
	"github.com/permitio/permit-golang/pkg/models"
	"github.com/permitio/permit-golang/pkg/permit"
)

func main() {
	// Initialize the Permit client
	permitConfig := config.NewConfigBuilder(
		"<API_KEY>"). //Please insert your API KEY
		WithPdpUrl("http://localhost:7766"). // change if needed according to PDP external port
		Build()
	permitClient := permit.NewPermit(permitConfig)

	// Set up routes
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to the Go application!")
		fmt.Fprintln(w, "This is a simple welcome page.")
	}).Methods("GET")

	r.HandleFunc("/signup", func(w http.ResponseWriter, r *http.Request) {
		// Perform user signup
		ctx := r.Context()
		var userData struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		newUser := models.NewUserCreate(userData.Name)
		user, err := permitClient.SyncUser(ctx, *newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "User synced successfully. Key: %s", user.Key)
	}).Methods("POST")

	r.HandleFunc("/transfer_payment_for_blog", func(w http.ResponseWriter, r *http.Request) {
		// Extract user name from request body
		var userData struct {
			UserName string `json:"user_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Check if user exists and create a user object
		user := enforcement.UserBuilder(userData.UserName).Build()

		// Mocked resource for demonstration
		resource := enforcement.ResourceBuilder("blog").Build()

		// Check if the user is permitted to perform the action
		permitted, err := permitClient.Check(user, "read", resource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if permitted {
			// Transfer payment logic goes here
			fmt.Fprintf(w, "Payment transferred successfully for blog")
		} else {
			http.Error(w, "Access denied", http.StatusForbidden)
		}
	}).Methods("POST")

	r.HandleFunc("/blogs", func(w http.ResponseWriter, r *http.Request) {
		// Extract user name from request body
		var userData struct {
			UserName string `json:"user_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Check if user exists and create a user object
		user := enforcement.UserBuilder(userData.UserName).Build()

		// Mocked resource for demonstration
		resource := enforcement.ResourceBuilder("blog").Build()

		// Check if the user is permitted to perform the action
		permitted, err := permitClient.Check(user, "read", resource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !permitted {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		// Logic to list blogs
		// (You would implement this logic according to your application's requirements)
		blogs := []string{"blog1", "blog2", "blog3"}
		blogsJSON, _ := json.Marshal(blogs)
		w.Header().Set("Content-Type", "application/json")
		w.Write(blogsJSON)
	}).Methods("GET")

	r.HandleFunc("/blogs", func(w http.ResponseWriter, r *http.Request) {
		// Extract user name from request body
		var userData struct {
			UserName string `json:"user_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		// Check if user exists and create a user object
		user := enforcement.UserBuilder(userData.UserName).Build()

		// Mocked resource for demonstration
		resource := enforcement.ResourceBuilder("blog").Build()

		// Check if the user is permitted to perform the action
		permitted, err := permitClient.Check(user, "write", resource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !permitted {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		// Logic to create a new blog
		// (You would implement this logic according to your application's requirements)
		fmt.Fprintf(w, "New blog created successfully")
	}).Methods("POST")

	r.HandleFunc("/blogs/{blog_id}", func(w http.ResponseWriter, r *http.Request) {
		// Extract user name from request body
		var userData struct {
			UserName string `json:"user_name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
	
		// Check if user exists and create a user object
		user := enforcement.UserBuilder(userData.UserName).Build()
	
		// Mocked resource for demonstration
		resource := enforcement.ResourceBuilder("blog").Build()
	
		// Check if the user is permitted to perform the action
		permitted, err := permitClient.Check(user, "delete", resource)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if !permitted {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}
	
		// Get the blog ID from the request path parameters
		vars := mux.Vars(r)
		blogID := vars["blog_id"]
	
		// Check if the blog ID exists (you would need to implement this logic based on your application)
		if !blogExists(blogID) {
			http.Error(w, "Blog not found", http.StatusNotFound)
			return
		}
	
		// Delete the blog (you would need to implement this logic based on your application)
		err = deleteBlog(blogID)
		if err != nil {
			http.Error(w, "Failed to delete blog", http.StatusInternalServerError)
			return
		}
	
		// Return a success message
		fmt.Fprintf(w, "Blog with ID %s deleted successfully", blogID)
	}).Methods("DELETE")	

	// Start the server
	port := ":8080"
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}

// Function to check if a blog with the specified ID exists (example implementation)
func blogExists(blogID string) bool {
	// Implement logic to check if the blog exists in your data store
	// For simplicity, this is just a placeholder implementation
	return true
}

// Function to delete a blog with the specified ID (example implementation)
func deleteBlog(blogID string) error {
	// Implement logic to delete the blog from your data store
	// For simplicity, this is just a placeholder implementation
	return nil
}
