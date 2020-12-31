package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lethekhoi/twofactor/api/user_api"
	"github.com/lethekhoi/twofactor/entities"
	"github.com/lethekhoi/twofactor/service"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

var tpl *template.Template

func Index(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "index.gohtml", nil)
}
func RegisterPage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "register.gohtml", nil)
}
func LoginPage(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}
func Success(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Success!\n")
}
func Fail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Fail!\n")
}
func SubmitCode(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "code.gohtml", nil)
}
func ScanQRCode(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "codeRegister.gohtml", nil)
}
func loginUserPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))

		url := "/api/user/login/" + username + "/" + password
		fmt.Println("url" + url)
		http.Redirect(w, r, url, http.StatusFound)
	}
}
func processor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing...")
	if r.Method == "POST" {

		code := r.FormValue("otpcode")
		fmt.Println("code :" + code)
		fmt.Println("SecretOTPKey :" + os.Getenv("SecretOTPKey"))
		valid := totp.Validate(code, os.Getenv("SecretOTPKey"))
		if valid {
			println("Valid passcode!")
			http.Redirect(w, r, "/api/user/success", http.StatusFound)
		} else {
			println("Invalid passcode!")
			http.Redirect(w, r, "/api/user/fail", http.StatusFound)
		}
	}

}

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "POST" {
		url := "http://localhost:5000/api/user/create"
		username := strings.TrimSpace(r.FormValue("username"))
		password := strings.TrimSpace(r.FormValue("password"))
		user := &entities.User{Username: username, Password: password}
		var inputBody bytes.Buffer
		if err := json.NewEncoder(&inputBody).Encode(user); err != nil {
			fmt.Println("q")
			return
		}

		client := http.Client{}

		req, err := http.NewRequest("POST", url, &inputBody)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Unable to reach the server.")
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Println("body=", string(body))
			var newUser entities.User
			err := json.Unmarshal(body, &newUser)
			if err != nil {

				http.Redirect(w, r, "/api/user/fail", http.StatusFound)
			}
			key, buf, _ := service.TOTPKey(user)

			display(key, buf.Bytes())

			http.Redirect(w, r, "/ScanQRCode", http.StatusFound)
		}

	}
}
func display(key *otp.Key, data []byte) {
	fmt.Printf("Issuer:       %s\n", key.Issuer())
	fmt.Printf("Account Name: %s\n", key.AccountName())
	fmt.Printf("Secret:       %s\n", key.Secret())
	os.Setenv("SecretOTPKey", key.Secret())
	fmt.Println("Writing PNG to qr-code.png....")
	ioutil.WriteFile("./static/qr-code.png", data, 0644)
	fmt.Println("")
	fmt.Println("Please add your TOTP to your OTP Application now!")
	fmt.Println("")
}
func main() {

	router := mux.NewRouter().StrictSlash(true)
	// Choose the folder to serve
	staticDir := "/static/"
	// Create the route
	router.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	router.HandleFunc("/", Index)
	router.HandleFunc("/SubmitCode", SubmitCode)
	router.HandleFunc("/ScanQRCode", ScanQRCode)
	router.HandleFunc("/loginuserpassword", loginUserPassword)
	router.HandleFunc("/register", RegisterPage)
	router.HandleFunc("/registerUser", RegisterUser)
	router.HandleFunc("/login", LoginPage)
	router.HandleFunc("/api/user/success", Success)
	router.HandleFunc("/api/user/fail", Fail)
	router.HandleFunc("/process", processor)
	router.HandleFunc("/api/user/create", user_api.Create).Methods("POST")
	router.HandleFunc("/api/user/login/{user}/{password}", user_api.Login).Methods("GET")

	fmt.Println("Listen port")
	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:5000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
