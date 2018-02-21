package main

import (
	"log"

	"github.com/vuolen/reaktor-k18/server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	s, err := server.New()
	if err != nil {
		panic(err)
	}
	defer s.Close()
	panic(s.Run())
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
