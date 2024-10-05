package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Amangupta20000/mongoapi/router"
)

func main() {
	fmt.Println("MongoDb")
	r := router.Router()
	fmt.Println("Server is getting Started...")
	log.Fatal(http.ListenAndServe(":5000", r))
	fmt.Println("Listening on PORT 5000 ...")

}
