package authentication_test

import (
	"bytes"
	"fmt"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/tester"
	"gin-boilerplate/internal/utilities"
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
	bootstrap.DB.Where("username = ?", "test_123_abc").Delete(&migration.User{})
	var user, token = tester.GenerateAuthentication(bootstrap, "customer", "customer")
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

func TestAccountPasswordHttp(t *testing.T) {
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
	_ = writer.WriteField("current_password", "123456")
	_ = writer.WriteField("username", "test_123_abc")
	_ = writer.WriteField("email", "test_123_abc@gmail.com")

	// Close multipart writer
	err := writer.Close()
	if err != nil {
		t.Error("Error closing writer:", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/account/password", body)
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
	assert.Equal(t, "test_123_abc", user.Username)
	assert.Equal(t, "test_123_abc@gmail.com", user.Email.String)
}

func TestAccountNewPasswordHttp(t *testing.T) {
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
	_ = writer.WriteField("current_password", "123456")
	_ = writer.WriteField("new_password", "abc123")

	// Close multipart writer
	err := writer.Close()
	if err != nil {
		t.Error("Error closing writer:", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/v1/account/password", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Secret-Key", bootstrap.Env.App.Key)
	req.Header.Set("Authorization", TestAccount.token.Token)
	bootstrap.Engine.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		fmt.Println("Err Body", w.Body.String())
	}

	var user = migration.User{}
	bootstrap.DB.Where("id = ?", TestAccount.user.Id).First(&user)

	// check if the new password is correctly set
	_, err = utilities.VerifyHash("abc123", user.Password)
	if err != nil {
		t.Fatal("Password do not match:", err)
	}

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test_123_abc", user.Username)
	assert.Equal(t, "test_123_abc@gmail.com", user.Email.String)
}
