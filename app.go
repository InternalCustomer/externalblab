package main

import (
	"log"
	"net/http"
	"os"
	"html/template"
	//for extracting service credentials from VCAP_SERVICES
	"github.com/cloudfoundry-community/go-cfenv"
)

const (
	DEFAULT_PORT = "8080"
)

var index = template.Must(template.ParseFiles(
  "templates/_base.html",
  "templates/index.html",
))

func helloworld(w http.ResponseWriter, req *http.Request) {
  index.Execute(w, nil)
}

func main() {
	appEnv, err := cfenv.Current()
	if appEnv != nil {

		log.Printf("ID %+v\n", appEnv.ID)
	}
  if err != nil {

		log.Printf("err")
	}
	log.Printf("appEnv.Services: \n%+v\n", appEnv.Services)
	//log.Printf("Cloudant credentials: \n%+v\n", appEnv.Services)

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}

	http.HandleFunc("/", helloworld)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Starting app on port %+v\n", port)
	http.ListenAndServe(":"+port, nil)
}
