package authentication_test

import (
	"bytes"
	"fmt"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/tester"
	"gin-boilerplate/pkg"
	"gin-boilerplate/pkg/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestAccountModel struct {
	user  migration.User
	token *migration.AuthenticationToken
}

var TestAccount = TestAccountModel{}

func TestAccountSetup(t *testing.T) {
	var env = "test"
	var bootstrap = pkg.InitBoot("../../../env.json", &env)

	// Set Gin to Test mode
	gin.SetMode(gin.TestMode)

	// routes
	routes.V1Routes(bootstrap)

	// mock authentication
	var user, token = tester.GenerateAuthentication(bootstrap, "customer", "admin")
	TestAccount.user = user
	TestAccount.token = token
}

func TestAccountInformationHttp(t *testing.T) {
	var env = "test"
	var bootstrap = pkg.InitBoot("../../../env.json", &env)

	// Set Gin to Test mode
	gin.SetMode(gin.TestMode)

	// routes
	routes.V1Routes(bootstrap)

	// Create a new multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add other fields
	_ = writer.WriteField("first_name", "User")
	_ = writer.WriteField("last_name", "One")

	// Close multipart writer
	err := writer.Close()
	if err != nil {
		t.Error("Error closing writer:", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/account/information", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Secret-Key", bootstrap.Env.App.Key)
	req.Header.Set("Authorization", TestAccount.token.Token)
	bootstrap.Engine.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		fmt.Println("Err Body", w.Body.String())
	}

	var user = migration.User{}
	bootstrap.DB.Where("id = ?", TestAccount.user.Id).First(&user)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "User", user.FirstName)
	assert.Equal(t, "One", user.LastName.String)
}
