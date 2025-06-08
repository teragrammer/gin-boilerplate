package authentication_test

import (
	"database/sql"
	"fmt"
	"gin-boilerplate/database/migration"
	"gin-boilerplate/internal/tester"
	"gin-boilerplate/internal/utilities"
	"gin-boilerplate/pkg"
	"gin-boilerplate/pkg/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestTwoFactoAuthenticationModel struct {
	user  migration.User
	token *migration.AuthenticationToken
}

var TestTwoFactoAuthentication = TestTwoFactoAuthenticationModel{}

func TestTwoFactoAuthenticationSetupHttp(t *testing.T) {
	var env = "test"
	var bootstrap = pkg.InitBoot("../../../env.json", &env)

	// enable tfa_req temporarily
	bootstrap.DB.Model(&migration.Setting{}).Where("slug = 'tfa_req'").
		Updates(migration.Setting{
			Value: &utilities.NullString{NullString: sql.NullString{Valid: true, String: "1"}},
		})

	// mock authentication
	var user, token = tester.GenerateAuthentication(bootstrap, "admin", "admin")
	TestTwoFactoAuthentication.user = user
	TestTwoFactoAuthentication.token = token
}

func TestTwoFactoAuthenticationSendHttp(t *testing.T) {
	var env = "test"
	var bootstrap = pkg.InitBoot("../../../env.json", &env)

	// Set Gin to Test mode
	gin.SetMode(gin.TestMode)

	// routes
	routes.V1Routes(bootstrap)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/tfa/send", nil)
	req.Header.Set("X-Secret-Key", bootstrap.Env.App.Key)
	req.Header.Set("Authorization", TestTwoFactoAuthentication.token.Token)
	bootstrap.Engine.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		fmt.Println("Err Body", w.Body.String())
	}
	assert.Equal(t, http.StatusOK, w.Code)
}
