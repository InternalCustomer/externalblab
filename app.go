package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"html/template"
	"bytes"
	//for extracting service credentials from VCAP_SERVICES
	"github.com/cloudfoundry-community/go-cfenv"
)

type post struct {

	Id string `json:"_id"`
	Rev string `json:"_rev"`
	Title string
	Author string
	Date string
	Text string
	Category []string
}

type toPost struct {

	Title string
	Author string
	Date string
	Text string
	//Category []string
}

type cloudant_data struct {

	Total_rows int //int
	Offset int //int
	Rows []struct {
        Id string
    }
}

var urlDb string = ""

const (
	DEFAULT_PORT = "8080"
)

var index = template.Must(template.ParseFiles(
  "templates/_base.html",
  "templates/index.html",
))

var blab = template.Must(template.ParseFiles(
	"templates/_base.html",
  "templates/blab.html",
))

func blabHandler(w http.ResponseWriter, req *http.Request) {

  postBody := "<p></p><hr><p>In progress</p>"
	page := struct {
				Title	string
			  Body interface{}
		}{"BLAB",template.HTML(postBody)}

		page.Title = "BLAB"

    blab.Execute(w, page)
}

//Index Page - about
func indexHandler(w http.ResponseWriter, req *http.Request) {
  index.Execute(w, nil)
}

func save(w http.ResponseWriter, r *http.Request) {
    name := r.FormValue("author")
    text := r.FormValue("text")
		title := r.FormValue("title")

    data := &toPost{title,name,"", text}

    b, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    log.Printf("---------- DATA RECEIVED----------")
		log.Printf("%s", b)

		log.Printf("------------URL_DB--------------")
		log.Printf("%s", urlDb)

		log.Printf("-------- MAKING POST ON DB ------------")

		dbToPost := urlDb+"/blab_data/"

		req, err := http.NewRequest("POST", dbToPost, bytes.NewBuffer(b))
 	 //req.Header.Set("X-Custom-Header", "myvalue")
 	 req.Header.Set("Content-Type", "application/json")

 	 client := &http.Client{}
 	 resp, err := client.Do(req)
 	 if err != nil {
 			 panic(err)
 	 }
 	 defer resp.Body.Close()

 	 log.Println("response Status:", resp.Status)
 	 log.Println("response Headers:", resp.Header)
 	 body, _ := ioutil.ReadAll(resp.Body)
 	 log.Println("response Body:", string(body))

	 http.Redirect(w, r, "/blab", 301)

}

func init() {
    http.HandleFunc("/save", save)
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
	urlDb = creds["url"].(string)

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

	http.HandleFunc("/", indexHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/blab", blabHandler)

	log.Printf("Starting app on port %+v\n", port)
	http.ListenAndServe(":"+port, nil)
}
