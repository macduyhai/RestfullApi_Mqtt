package main

import (
	"RestfullApi_Mqtt/msgmqtt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

type device struct {
	Key string `json:"key"`
	Mac string `json:"mac"`
	Id  int    `json:"id"`
}
type customClaims struct {
	Payload string `json:"payload"`
	jwt.StandardClaims
}

var Device device

// JWT
var jwtSecretKey = []byte("eNhomKou0CMJ694nK281vghbb6UtIQB2")
var msg = "{\"sensor\":\"gps\",\"time\":1351824120}"

// CreateJWT func will used to create the JWT while signing in and signing out
func CreateJWT(payload string) (response string, err error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := customClaims{
		Payload: payload,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "nameOfWebsiteHere",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err == nil {
		return tokenString, nil
	}
	return "", err
}

// VerifyToken func will used to Verify the JWT Token while using APIS
func VerifyToken(tokenString string) (tokenstr string, err error) {
	var claims customClaims

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if token != nil {
		return claims.Payload, nil
	}
	return "", err
}

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
		// fmt.Fprintf(w, `{"error_code":10000}`)
	}

	// json.Unmarshal(reqBody, &Device)
	err1 := json.Unmarshal(reqBody, &Device)
	if err1 != nil {
		log.Println(err)
	}
	fmt.Println(Device)
	mac := Device.Mac
	key := Device.Key
	id := Device.Id
	if key == "lvJvDWKiv0" {
		fmt.Printf("Mac device:%s\t-Trạng thái: %s \n", mac, stt)
		if stt == "1" {
			fmt.Println("Bật đèn")
			s := "{\"id\":" + strconv.Itoa(id) + "," + "\"value\":\"1\"}"
			fmt.Println(s)
			msgmqtt.PublishData(mac, s)
		} else if stt == "0" {
			fmt.Println("Tắt đèn")
			s := "{\"id\":" + strconv.Itoa(id) + "," + "\"value\":\"0\"}"
			fmt.Println(s)
			msgmqtt.PublishData(mac, s)
		}
		fmt.Println("------------------------------------------")
		fmt.Fprintf(w, `{"error_code":10000}`)
	} else {
		fmt.Println("Sai Key .! Vui long check lai API")
		fmt.Fprintf(w, `{"error_code":10002,"alert":"Key not true"}`)
	}

}

func main() {
	fmt.Println("====> START MAIN <=====")
	msgmqtt.MqttBegin()
	// time.Sleep(10)
	// s, _ := CreateJWT(msg)
	// fmt.Println(s)
	// fmt.Println("======")
	// payload, err := VerifyToken(s)
	// if err != nil {
	// 	fmt.Print("Paload:")
	// 	fmt.Println(payload)
	// } else {
	// 	fmt.Println("Error decode")
	// 	fmt.Println(err)
	// }
	fmt.Println("======")
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/add-device", AddDevice).Methods("POST")
	router.HandleFunc("/control-device/{stt}", ControDevice).Methods("POST")
	router.HandleFunc("/get-list-device", GetlistDevice).Methods("GET")
	log.Fatal(http.ListenAndServe(":8888", router))
}
