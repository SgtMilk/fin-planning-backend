package tests

import (
	"net/http"
	"testing"

	"github.com/SgtMilk/fin-planning-backend/database"
)

func TestCreateUserSuccess(t *testing.T) {
	testUsername := "A%uHr%ksOUhs"
	testPassword:= "iUgh*h)6#kjdUFh"
	router := InitializeTestEnv(t)

	body := GenerateJSON(`{"username" : "` + testUsername + `", "password" : "` + testPassword + `"}`)
	req, _ := http.NewRequest("POST", "/auth/register", body)
	router.ServeRequest(req)

	testCreateSuccess(&router, testUsername, testPassword)
}

func TestCreateUserFail1(t *testing.T) {
	credentialsTest(t, "aaaaaaaa", "iUgh*h)6#kjdUFh", "insecure username, try using a longer username");
}

func TestCreateUserFail2(t *testing.T) {
	credentialsTest(t, "A%uHr%ksOUhs", "bbbbbbbb", "insecure password, try including more special characters, using uppercase letters, using numbers or using a longer password");
}

func TestCreateUserFail3(t *testing.T) {
	credentialsTest(t, "A%uH", "bbbbbbbb", "username not of right size");
}

func TestCreateUserFail4(t *testing.T) {
	credentialsTest(t, "aaaaaaaa", "b", "password not of right size");
}

func TestCreateUserFail5(t *testing.T) {
	testUsername := ""
	for i := 0 ; i < 257 ; i++ {
		testUsername += "a"
	} 

	testPassword := ""
	for i := 0 ; i < 73 ; i++ {
		testPassword += "b"
	} 
	credentialsTest(t, testUsername, testPassword, "username and password not of right size");
}

func TestCreateUserFail6(t *testing.T) {
	username := "A%uHr%ksOUhs"
	router := InitializeTestEnv(t)

	body := GenerateJSON(`{"username" : "` + username + `"}`)
	req, _ := http.NewRequest("POST", "/auth/register", body)
	router.ServeRequest(req)

	testCreateFail(&router, username, "Key: 'AuthenticateInput.Password' Error:Field validation for 'Password' failed on the 'required' tag")
}

func TestCreateUserFail7(t *testing.T) {
	username := "A%uHr%ksOUhs"
	router := InitializeTestEnv(t)

	body := GenerateJSON(`{"password" : "` + username + `"}`)
	req, _ := http.NewRequest("POST", "/auth/register", body)
	router.ServeRequest(req)

	testCreateFail(&router, username, "Key: 'AuthenticateInput.Username' Error:Field validation for 'Username' failed on the 'required' tag")
}

func testCreateSuccess(router *TestRouter, username string, password string){
	router.AssertCode(201)

	user, err := database.FindUserByUsername(username)
	router.FailOnError(err, "User wasn't created.")

	err = user.ValidatePassword(password)
	router.FailOnError(err, "Password is wrong.")

	options := user.Options
	router.AssertEqual(uint16(12), options.MonthInterval)
	router.AssertEqual(float32(4), options.Inflation)
	router.AssertEqual(float32(50), options.TaxRate)
	router.AssertEqual("2023-11", options.StartMonth)
	router.AssertEqual("2073-11", options.EndMonth)

	// cleanup
	err = user.Delete()
	if err != nil{
		router.T.Fatalf("User not properly deleted.")
	}
}

func testCreateFail(router *TestRouter, username string, errorMessage string){
	router.AssertCode(400)
	router.AssertBody(`{"error":"` + errorMessage + `"}`)

	_, err := database.FindUserByUsername(username)
	router.FailOnNoError(err, "User was created when it should not.")
}

func credentialsTest(t *testing.T, username string, password string, errorMessage string){
	router := InitializeTestEnv(t)

	body := GenerateJSON(`{"username" : "` + username + `", "password" : "` + password + `"}`)
	req, _ := http.NewRequest("POST", "/auth/register", body)
	router.ServeRequest(req)

	testCreateFail(&router, username, errorMessage)
}