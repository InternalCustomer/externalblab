package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"html/template"
	//for extracting service credentials from VCAP_SERVICES
	"github.com/cloudfoundry-community/go-cfenv"
)

type post struct {

	Id string `json:"_id"`
	Rev string `json:"_rev"`
	Title string
	Author string
	Data string
	Text string
	Category []string
}

type cloudant_data struct {
	
	Total_rows int //int
	Offset int //int
	Rows []struct {
        Id string
    }
}

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

	cloudantServices, err := appEnv.Services.WithLabel("cloudantNoSQLDB")
  if err != nil || len(cloudantServices) == 0 {
    log.Printf("No Cloudant service info found\n")
    return
  }

  creds := cloudantServices[0].Credentials
  username := creds["username"]
	password := creds["password"]

	//ACCESS TO CLOUDANT
	//GET https://$USERNAME:$PASSWORD@$USERNAME.cloudant.com
	basicUrl := "https://" + username.(string) + ":" + password.(string) + "@" + username.(string) + ".cloudant.com"

	resp, err := http.Get(basicUrl)

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		//fmt.Printf("%s", err)
		log.Printf("Err body: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Body: %s\n", string(body))
	//fmt.Printf("%s\n", string(contents))


	//LIST ALL DOCS IN blab_data
	//GET https://$USERNAME:$PASSWORD@$USERNAME.cloudant.com/blab_data/_all_docs
	allDocsUrl := basicUrl + "/blab_data/_all_docs"

	resp, err = http.Get(allDocsUrl)

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	//encoding JSON response - post info
	var p post
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Printf("Err body: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Body: %s\n", string(body))
	log.Printf("Title: %+v\n",p)
	log.Printf("Title: %s\n",p.Title)
	log.Printf("author: %s\n",p.Author)

	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}

	http.HandleFunc("/", helloworld)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Starting app on port %+v\n", port)
	http.ListenAndServe(":"+port, nil)
}
