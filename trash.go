package main

import (
  "log"

  "github.com/ziutek/mymysql/mysql"
  _ "github.com/ziutek/mymysql/native"

  "github.com/cloudfoundry-community/go-cfenv"
)

var (
  db mysql.Conn
)

func init() {
  appEnv, err := cfenv.Current()
  if err != nil {
    log.Printf("Database information not available continuing without database support.\n")
    // TODO use an environment variable to get the local development environment mysql database
    return
  }

  dbServices, err := appEnv.Services.WithLabel("cleardb")
  if err != nil || len(dbServices) == 0 {
    log.Printf("No cleardb database info found\n")
    return
  }

  creds := dbServices[0].Credentials

  db = mysql.New("tcp", "", creds["hostname"]+":"+creds["port"], creds["username"], creds["password"], creds["name"])

  err = db.Connect()
  if err != nil {
    db = nil
    log.Printf("Error connecting to database: %v\n", err.Error())
    return
  }

  // Check for a database table first, create it if necessary
  _, _, err = db.QueryFirst("CREATE TABLE IF NOT EXISTS SIGNUPS (ID int AUTO_INCREMENT PRIMARY KEY, NAME VARCHAR(50), COMING BIT(1), COMMENT VARCHAR(100))")
  if err != nil {
    db = nil
    log.Printf("Error creating signup table: %v\n", err.Error())
    return
  }

  //_, _, err = db.Query("select * from X where id > %d", 20)
  //if err != nil {
  //panic(err)
  //}
}
