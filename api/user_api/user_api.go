package user_api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lethekhoi/twofactor/config"
	"github.com/lethekhoi/twofactor/entities"
	"github.com/lethekhoi/twofactor/models"
	"github.com/lethekhoi/twofactor/service"
	"github.com/pquerna/otp"
)

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

//Create ...
func Create(response http.ResponseWriter, request *http.Request) {
	fmt.Println("c")
	var user entities.User
	err := json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		fmt.Println("json error")
		respondWithError(response, http.StatusBadRequest, err.Error())
	}
	db, err := config.GetDB()

	defer db.Close()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	} else {
		userModel := models.UserModel{
			Db: db,
		}
		err := userModel.Create(&user)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, err.Error())
		} else {
			respondWithJSON(response, http.StatusOK, user)
		}
	}
}

//Login ...
func Login(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	username := vars["user"]
	password := vars["password"]
	db, err := config.GetDB()

	defer db.Close()
	if err != nil {
		respondWithError(response, http.StatusBadRequest, err.Error())
	}
	{
		userModel := models.UserModel{
			Db: db,
		}
		user, err := userModel.Login(username, password)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, err.Error())
			return
		}
		{
			key, buf, _ := service.TOTPKey(&user)

			display(key, buf.Bytes())

			http.Redirect(response, request, "/SubmitCode", http.StatusFound)

		}
	}
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
func promptForPasscode() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Passcode: ")
	text, _ := reader.ReadString('\n')
	fmt.Print("text:", text)
	return text
}
