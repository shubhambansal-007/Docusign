package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"project/handlers"

	"github.com/gorilla/mux"
)

//					do this on console
//		cd C:\Users\Admin\Desktop\project
//		set GOOGLE_APPLICATION_CREDENTIALS=./appEngineCreds.json
//		go run main.go

//					to work on gcloud
//		gcloud init				do this on Google Cloud SDK Shell console
//      to deploy yaml file:- gcloud app deploy app.yaml index.yaml

// load configuration file
func loadConfigurationFile(fileName string) handlers.Data {
	config_file, err := ioutil.ReadFile(fileName)
	var dataOfConfigFile handlers.Data

	if err != nil {
		dataOfConfigFile = defaultConfigurationFile(dataOfConfigFile)
		fmt.Println("error in opening config.json file:- ", err.Error())
	}

	error := json.Unmarshal(config_file, &dataOfConfigFile)

	if error != nil {
		dataOfConfigFile = defaultConfigurationFile(dataOfConfigFile)
		fmt.Println("error during unmarshal to config file:- ", error.Error())
	}
	return dataOfConfigFile
}

// default code if loading configuration file show some error
func defaultConfigurationFile(configFile handlers.Data) handlers.Data {

	if configFile.Host == "" {
		configFile.Host = "localhost"
	}
	if configFile.Port == "" {
		configFile.Port = ":1111"
	}
	if configFile.ProjectID == "" {
		configFile.ProjectID = "urlmanager-386505"
	}
	if configFile.Valid_kind == "" {
		configFile.Valid_kind = "valid_urls"
	}
	if configFile.Invalid_kind == "" {
		configFile.Invalid_kind = "inValid_urls"
	}

	return configFile
}

func init() {

	fmt.Println("------------------------we are in func init()==================")
	cmd := exec.Command("set GOOGLE_APPLICATION_CREDENTIALS=./appEngineCreds.json")
	err := cmd.Start()
	if err != nil {
		fmt.Println("error in staring cmd:------", err)
	}
}

func main() {

	//creating router
	getRouter := mux.NewRouter()
	fmt.Println("------------------------we are in main.go==================")
	//getting data from configuration file
	var configFile handlers.Data
	configFile = loadConfigurationFile("config.json")
	configFile = defaultConfigurationFile(configFile)
	fmt.Println("config file is: ", configFile)

	//getting data in form of request
	getRouter.HandleFunc("/write", configFile.WriteJSONHandleFunc).Methods("POST")

	getRouter.HandleFunc("/read", handlers.ReadJSONHandleFunc).Methods("GET")

	getRouter.HandleFunc("/delete", handlers.DeleteJSONHandleFunc).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "1001"
		log.Print("defaulting to port------ ", port)
	}
	server := http.Server{}
	server.Handler = getRouter
	server.Addr = fmt.Sprintf(":%s", port)
	server.ListenAndServe()

}
