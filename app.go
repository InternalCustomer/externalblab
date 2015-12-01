package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"html/template"
	"bytes"
	"time"
	"strings"
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

func (self *post) SetDefaults() {
	if len(self.Author) == 0 { self.Author = "Anonymous" }
	if len(self.Title) == 0 { self.Title = "No Title" }
}

type toPost struct {

	Title string
	Author string
	Date string
	Text string
	//Category []string
}

type cloudant_data struct {
	Docs []post
	//Total_rows int //int
	//Offset int //int
	//Rows []struct {
    //    Id string
    //}
}

const (
	DEFAULT_PORT = "8080"
)

var index = template.Must(template.ParseFiles(
  "templates/index.html",
))

var blab = template.Must(template.ParseFiles(
	"templates/_base.html",
  "templates/blab.html",
))

var blabPost = template.Must(template.ParseFiles(
	"templates/_base.html",
  "templates/blabPost.html",
))

//To delete - for local host test
var basicUrl = "https://9bd28748-9a66-441b-8b98-48f993b17e8e-bluemix:483f8ba7b8507e15548befa5e1b9c53cd3535f2c7ec75504dec3369df37b058d@9bd28748-9a66-441b-8b98-48f993b17e8e-bluemix.cloudant.com"

func blabHandler(w http.ResponseWriter, req *http.Request) {

	//LIST ALL DOCS IN blab_data sorted by Date
	//POST https://$USERNAME:$PASSWORD@$USERNAME.cloudant.com/blab_data/_find
	//POST Body
	// 	{
	//   "selector": {
	//     "Date": {
	//       "$gt": 0
	//     }
	//   },
	//   "fields": [
	//     "_id",
	//     "Title",
	//     "Author",
	//     "Date"
	//   ],
	//   "sort": [
	//     {
	//       "Date": "desc"
	//     }
	//   ]
	// }
	//Query to sort posts by date
	var jsonString = []byte(`{ "selector": { "Date": { "$gt": 0 } }, "fields": [ "_id", "Title", "Author", "Date" ], "sort": [ { "Date": "desc" } ] }`)
	DocsSortedUrl := basicUrl + "/blab_data/_find"

	req, err := http.NewRequest("POST", DocsSortedUrl, bytes.NewBuffer(jsonString))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error scope: POST %s\nError message: %s\n", DocsSortedUrl, err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	//encoding JSON response
	var data cloudant_data
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Error encoding body response: %s\n", err)
		os.Exit(1)
	}

	docUrl := basicUrl + "/blab_data/"
	var postBody string

	for i := 0; i < len(data.Docs); i++ {

			//Getting info single post
			//GET https://$USERNAME:$PASSWORD@$USERNAME.cloudant.com/blab_data/<data_id>
			singleDocUrl := docUrl + data.Docs[i].Id
			resp, err := http.Get(singleDocUrl)
			if err != nil {
				log.Printf("Error scope: GET %s\nError message: %s\n", err)
				os.Exit(1)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			//encoding JSON response
			var p post
			err = json.Unmarshal(body, &p)
			if err != nil {
				log.Printf("Error encoding body response: %s\n", err)
				os.Exit(1)
			}
			p.SetDefaults()
			p.Date = strings.Replace(p.Date, "T", " ", 1)
			p.Date = strings.Replace(p.Date, "Z", "", 1)

			postBody += "<a href='../blabPost?id=" + p.Id + "' class ='list-group-item'><span class='list-group-item-heading' style='font-size: medium;font-weight: 600;'>"+ p.Title +"</span>&#09; - &#09;<span>"+ p.Author +"</span><br><h6>"+ p.Date +"</h6> </a>"
	}

	page := struct {
			Title	string
			Body interface{}
	}{"BLAB",template.HTML(postBody)}

	blab.Execute(w, page)
}

func blabPostHandler(w http.ResponseWriter, req *http.Request) {

	query := req.URL.Query()
	id := query.Get("id")

	singleDocUrl := basicUrl + "/blab_data/" + id
	var postBody string

	//Getting info single post
	//GET https://$USERNAME:$PASSWORD@$USERNAME.cloudant.com/blab_data/<data_id>
	resp, err := http.Get(singleDocUrl)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//encoding JSON response
	var p post
	p.SetDefaults()
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Printf("Err body: %s\n", err)
		os.Exit(1)
	}
	p.Date = strings.Replace(p.Date, "T", " ", 1)
	p.Date = strings.Replace(p.Date, "Z", "", 1)
	postBody += "<div class ='post'><h3>"+ p.Title +"</h3>&#9;<h5>"+ p.Author +"</h5><br><span class='data'>"+ p.Date +"</span>" + "</span><br><p class='text'>"+ p.Text +"</p>" + "<a href='../blab'> <<< </a>" + "</div><br>"

	page := struct {
			Title	string
			Body interface{}
	}{"BLAB POST",template.HTML(postBody)}

	blabPost.Execute(w, page)
}

//Index Page - about
func indexHandler(w http.ResponseWriter, req *http.Request) {
  index.Execute(w, nil)
}

func save(w http.ResponseWriter, r *http.Request) {

    name := r.FormValue("author")
    text := r.FormValue("text")
		title := r.FormValue("title")

		t := time.Now() //get current date
		//date := t[:strings.Index(t, ".")]
		date := t.Format(time.RFC3339)
		log.Printf("---------- date value----------")
		log.Printf("%s", date)

    data := &toPost{title,name,date, text}

    b, err := json.Marshal(data)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    log.Printf("---------- DATA RECEIVED----------")
		log.Printf("%s", b)

		log.Printf("------------URL_DB--------------")
		log.Printf("%s", basicUrl)

		log.Printf("-------- MAKING POST ON DB ------------")

		dbToPost := basicUrl +"/blab_data/"

		req, err := http.NewRequest("POST", dbToPost, bytes.NewBuffer(b))
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

	cloudantServices, err := appEnv.Services.WithLabel("cloudantNoSQLDB")
  if err != nil || len(cloudantServices) == 0 {
    log.Printf("No Cloudant service info found\n")
    return
  }

  creds := cloudantServices[0].Credentials
	basicUrl = creds["url"].(string)

	var port string
		if port = os.Getenv("PORT"); len(port) == 0 {
			port = DEFAULT_PORT
		}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/blab", blabHandler)
	http.HandleFunc("/blabPost", blabPostHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Starting app on port %+v\n", port)
	http.ListenAndServe(":"+port, nil)
}
