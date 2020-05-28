package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type device struct {
	Key string `json:"key"`
	Mac string `json:"mac"`
}

var Device device

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome IOT!")
}

// AddDevice asds
func AddDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "AddDevice coming soon!")
}

// GetlistDevice asds
func GetlistDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "GetlistDevice coming soon!")
}
func DeleteDevice(w http.ResponseWriter, r *http.Request) {

}
func GetSttDevice(w http.ResponseWriter, r *http.Request) {

}

// ControDevice asds
func ControDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("API control device")
	stt := mux.Vars(r)["stt"]
	reqBody, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(reqBody))
	if err != nil {
		fmt.Fprintf(w, "Post form not true.! Control device not support")
	} else {
		fmt.Fprintf(w, `{"error_code":10000}`)
	}

	// json.Unmarshal(reqBody, &Device)
	err1 := json.Unmarshal(reqBody, &Device)
	if err1 != nil {
		log.Println(err)
	}
	fmt.Println(Device)
	mac := Device.Mac
	fmt.Printf("Mac device:%s\t-Trạng thái: %s \n", mac, stt)

}

func main() {
	fmt.Println("====> START MAIN <=====")
	//msgmqtt.Initmqtt()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/add-device", AddDevice).Methods("POST")
	router.HandleFunc("/control-device/{stt}", ControDevice).Methods("POST")
	router.HandleFunc("/get-list-device", GetlistDevice).Methods("GET")
	log.Fatal(http.ListenAndServe(":8888", router))
}
