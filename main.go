package main

import (
	"bytes"
	"crypto/tls"
	_ "encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sasbury/mini"
	"net/http"
	"os"
	"time"
)

var conn *tls.Conn

type Message struct {
	ClientId int    `json:"clientId"`
	Message  string `json:"message"`
}

var AUTH_KEY string
var CONNECTIN_STRING string

func main() {
	var port string

	config, err := mini.LoadConfiguration("settings.ini")
	if err != nil {
		fmt.Println("ERROR: Missing settings.ini")
		os.Exit(1)
	} else {
		AUTH_KEY = config.String("authkey", "TEST_AUTH_KEY_2015")
		CONNECTIN_STRING = config.String("connectionString", "")
		port = config.String("websiteport", "5000")
	}
	fmt.Println("AUTH_KEY=", AUTH_KEY)

	router := httprouter.New()

	router.GET("/lights", listLights)
	// router.GET("/sensors", sc.SensorsList)

	fmt.Println("Port", port)
	fmt.Println("Service", CONNECTIN_STRING)
	fmt.Println("Auth", AUTH_KEY)
	http.ListenAndServe(":"+port, router)
}

func listLights(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	authKey := r.FormValue("auth")
	if authKey != AUTH_KEY {
		fmt.Fprintf(w, "%s", time.Now())
		return
	}

	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		panic("Error loading X509 key pair")
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	conn, err = tls.Dial("tcp", CONNECTIN_STRING, &config)
	defer conn.Close()
	if err != nil {
		fmt.Fprintf(w, "Error connecting to backend %s", err.Error())
		return
	}

	_, err = conn.Write([]byte(AUTH_KEY))
	if err != nil {
		fmt.Fprintf(w, "Error writing auth key: %s", err.Error())
		return
	}

	command := "list-lights"
	cmdName := r.FormValue("cmd")
	fmt.Println(cmdName)
	switch cmdName {
	case "lights":
		command = "list-lights"
		break
	case "all":
		command = "list-all"
		break
	case "sensors":
		command = "list-sensors"
		break
	case "test":
		command = "test"
		break
	}

	fmt.Println("Command:", command)
	socketString := `{"message": "` + command + `"}`

	_, err = conn.Write([]byte(socketString))
	if err != nil {
		fmt.Println(err.Error())
	}
	buf := make([]byte, 1500)
	conn.Read(buf)
	if err != nil {
		conn.Close()
		os.Exit(1)
	}
	stringCleaned := bytes.Trim(buf, "\x00")
	// var str string = fmt.Sprintf("%s", stringCleaned)
	// var message Message
	// err = json.Unmarshal([]byte(stringCleaned), &message)
	// if err != nil {
	// 	fmt.Println("Error parsing JSON", err.Error)
	// }
	// fmt.Println(str)
	// fmt.Printf("%s", message.Message)
	fmt.Fprintf(w, "%s", stringCleaned)
}
