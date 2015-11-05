package main

import (
	"bytes"
	"crypto/tls"
	_ "encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
)

var conn *tls.Conn

type Message struct {
	ClientId int    `json:"clientId"`
	Message  string `json:"message"`
}

func main() {
	service := "192.168.1.17:1201"
	port := "6015"
	fmt.Println("Service", service)
	fmt.Println("Port", port)

	router := httprouter.New()

	router.GET("/lights", listLights)
	// router.GET("/sensors", sc.SensorsList)

	http.ListenAndServe(":"+port, router)
}

func listLights(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	service := "192.168.1.17:1201"
	cert, err := tls.LoadX509KeyPair("client.pem", "client.key")
	if err != nil {
		panic("Error loading X509 key pair")
	}

	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}

	conn, err = tls.Dial("tcp", service, &config)
	defer conn.Close()
	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	_, err = conn.Write([]byte("test_auth_key\n"))
	if err != nil {
		fmt.Println(err.Error())
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
