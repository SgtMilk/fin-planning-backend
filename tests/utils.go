package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/SgtMilk/fin-planning-backend/controller"
	"github.com/SgtMilk/fin-planning-backend/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestRouter struct{
	Recorder *httptest.ResponseRecorder
	Router *gin.Engine
	T *testing.T
}

func InitializeTestEnv(t *testing.T) TestRouter{
	database.LoadDB()
	gin.SetMode(gin.ReleaseMode)
	router := controller.CreateRouter(true)
	recorder := httptest.NewRecorder()
	return TestRouter{
		Recorder: recorder,
		Router: router,
		T: t,
	}
}

func (testRouter *TestRouter) ServeRequest(req *http.Request){
	testRouter.Router.ServeHTTP(testRouter.Recorder, req)
}

func (testRouter *TestRouter) AssertCode(code int){
	assert.Equal(testRouter.T, code, testRouter.Recorder.Code)
}

func (testRouter *TestRouter) AssertBody(body string){
	assert.Equal(testRouter.T, body, testRouter.Recorder.Body.String())
}

func (testRouter *TestRouter) AssertEqual(expected any, actual any){
	assert.Equal(testRouter.T, expected, actual)
}

func (testRouter *TestRouter) FailOnError(err error, message string){
	if err != nil{
		testRouter.T.Fatalf(message)
	}
}

func (testRouter *TestRouter) FailOnNoError(err error, message string){
	if err == nil{
		testRouter.T.Fatalf(message)
	}
}

// JSON STUFF

func GenerateJSON(input string) *bytes.Reader {
	jsonBody := []byte(input)
 	bodyReader := bytes.NewReader(jsonBody)
	return bodyReader
}