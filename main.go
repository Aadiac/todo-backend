package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aadiac/todo/router"
	"github.com/rs/cors"
)

func main() {
	r := router.Router()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},//only alow this frondend 
		// AllowOriginFunc: func(orgin string) bool{
		// 	return true
		// },
		AllowedMethods: []string{"POST","PUT","OPTIONS","GET","DELETE"},
		AllowedHeaders: []string{"Content-Type","Authorization"},
		AllowCredentials: true,
		// Debug: true, //enable this for logging
	})

	handler := c.Handler(r)
	
	fmt.Println("Listening To port 4001...")
	log.Fatal(http.ListenAndServe(":4001", handler))
	
}
