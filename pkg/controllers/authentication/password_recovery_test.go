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

type TestPasswordRecoverModel struct {
	user  migration.User
	token *migration.AuthenticationToken
}

var TestPasswordRecover = TestPasswordRecoverModel{}

func TestPasswordRecoverySetup(t *testing.T) {
	var env = "test"
	var bootstrap = pkg.InitBoot("../../../env.json", &env)

	// mock authentication
	var user, token = tester.GenerateAuthentication(bootstrap, "admin", "admin")
	TestPasswordRecover.user = user
	TestPasswordRecover.token = token
}

func TestTPasswordRecoverySendHttp(t *testing.T) {
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
	_ = writer.WriteField("to", "email")
	_ = writer.WriteField("email", TestPasswordRecover.user.Email.String)

	// Close multipart writer
	err := writer.Close()
	if err != nil {
		t.Error("Error closing writer:", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/password-recovery/send", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Secret-Key", bootstrap.Env.App.Key)
	bootstrap.Engine.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		fmt.Println("Err Body", w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}
