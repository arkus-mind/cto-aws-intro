package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/JorgeAdd/intro5/homeworks-jorge/homeworksAPI/internal/service"
	"github.com/gorilla/mux"
)

func HomeworksController() *mux.Router {

	fmt.Println("GO running")
	router := mux.NewRouter()

	// Get homework
	router.HandleFunc("/homework/{homeworkId}", service.GetHomework).Methods("GET")

	// Get homeworks
	router.HandleFunc("/homeworks", service.GetHomeworks).Methods("GET")

	// Create homework
	router.HandleFunc("/homework", service.CreateHomework).Methods("POST")

	// Update homework
	router.HandleFunc("/homework/{homeworkId}", service.UpdateHomework).Methods("PUT")

	// Delete homework
	router.HandleFunc("/homework/{homeworkId}", service.DeleteHomework).Methods("DELETE")

	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))

	return router
}
