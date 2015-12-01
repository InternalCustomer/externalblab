# Blab Application Overview

Blab is a simple, reusable Go web application.
Users can add and view posts on the main dashboard page.

It uses **Cloudant NoSQL db** from Bluemix Alchemy services.

The following code shows how to handle Cloudant service in Go:

```
	appEnv, err := cfenv.Current()
	  if appEnv != nil {
		log.Printf("ID %+v\n", appEnv.ID)
	  }
      if err != nil {
		log.Printf("err")
	  }

	cloudantServices, err := appEnv.Services.WithLabel("cloudantNoSQLDB")
      if err != nil || len(cloudantServices) == 0 {
       log.Printf("No Cloudant service info found\n")
       return
      }

    creds := cloudantServices[0].Credentials
	basicUrl = creds["url"].(string) // URL to use for making POST and GET request
```


## Run the app locally

1. [Install Go][]
2. Download and extract the code into your $GOPATH/src directory
3. cd into the app directory
4. Run `go run app.go`
5. Access the running app in a browser at http://localhost:8080

[Install Go]: https://golang.org/doc/install
