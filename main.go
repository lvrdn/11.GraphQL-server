package main

import "net/http"

func main() {
	port := "8080"

	handler := GetApp()

	http.ListenAndServe(":"+port, handler)
}
