package main

import (
	"net/http"

	"github.com/rikughi/technical-test/controller"
	"github.com/rikughi/technical-test/db"
)

// 17:40
func main() {
	db.Init()

	controller := controller.NewController()

	router := http.NewServeMux()

	router.HandleFunc("/data", controller.SetupA)
	router.HandleFunc("/data/", controller.SetupB)

	http.ListenAndServe(":3000", router)
}

// func extractID(path string) string {
// 	// Split the path by "/" and extract the last segment (id)
// 	parts := strings.Split(path, "/")
// 	if len(parts) > 1 {
// 		return parts[len(parts)-1]
// 	}
// 	return ""
// }
