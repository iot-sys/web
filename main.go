package main

import (
	"log"
	"net/http"

	"github.com/iot-sys/web/app"
	"github.com/urfave/negroni"
)

func main() {
	//m := app.MakeHandler("./test.db")
	m := app.MakeHandler("iotsys:111111@tcp(192.168.0.212:3306)/iotdb")
	defer m.Close()
	n := negroni.Classic()
	n.UseHandler(m)

	log.Println("Started App")
	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
