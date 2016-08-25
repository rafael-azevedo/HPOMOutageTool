package database

import (
	"fmt"
	"log"

	r "gopkg.in/dancannon/gorethink.v2"
)

func Example() {
	session, err := r.Connect(r.ConnectOpts{
		Address: "127.0.0.1:29015",
	})
	if err != nil {
		log.Fatalln(err)
	}

	res, err := r.Expr("Hello World").Run(session)
	if err != nil {
		log.Fatalln(err)
	}

	var response string
	err = res.One(&response)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(response)

	// Output:
	// Hello World
}
