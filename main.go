package main

import (
	"RestfullApi_Mqtt/msgmqtt"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	reqBody, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(reqBody))
	if err != nil {
		fmt.Fprintf(w, "Post form not true.! ")
	} else {
		// fmt.Fprintf(w, `{"error_code":10000}`)
	}
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
func I2SSample(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Post form not true.! Control device not support")
	} else {
		fmt.Fprintf(w, `{"error_code":10000}`)
	}
	//fmt.Println(string(reqBody))

	var _, err_file = os.Stat("i2s.raw")

	// create file if not exists
	if os.IsNotExist(err_file) {
		// open output file
		fo, err := os.Create("i2s.raw")
		if err != nil {
			panic(err)
		}
		fmt.Println("Tao file")
		// close fo on exit and check for its returned error
		defer func() {
			if err := fo.Close(); err != nil {
				panic(err)
			}
		}()
	}
	var file, err_b = os.OpenFile("i2s.raw", os.O_RDWR|os.O_APPEND, 0777)
	if err_b != nil {
		////return
		fmt.Printf(" Write_log %v\n", err_b)
	}
	defer file.Close()
	_, err = file.Write(reqBody)
	if err != nil {
		fmt.Printf(" Write_log %v\n", err)
		////return
	}
	// save changes
	err = file.Sync()
	if err != nil {
		fmt.Printf(" Write_log %v", err)
		////return
	}
	// err1 := ioutil.WriteFile("i2s.raw", reqBody, 0777) //[]byte(string(reqBody))
	// if err != nil {
	// 	log.Fatal(err1)
	// } else {
	// 	fmt.Println("Ghi thêm file")

	// }

}
func ADCSample(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Post form not true.! Control device not support")
	} else {
		fmt.Fprintf(w, `{"error_code":10000}`)
	}
	fmt.Println(string(reqBody))
	// open output file
	fo, err := os.Create("adc.raw")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	err1 := ioutil.WriteFile("adc.raw", []byte(string(reqBody)), 0644)
	if err != nil {
		log.Fatal(err1)
	}

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
	// msgmqtt.MqttBegin()
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
	router.HandleFunc("/i2s_samples", I2SSample).Methods("POST")
	router.HandleFunc("/adc_samples", ADCSample).Methods("POST")
	log.Fatal(http.ListenAndServe(":8888", router))

}
