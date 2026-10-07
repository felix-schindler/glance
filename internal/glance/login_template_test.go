package glance

import (
	"bytes"
	"strings"
	"testing"
)

func loginTemplateData(app *application) *templateData {
	return &templateData{
		App:            app,
		ShowLocalLogin: app.showLocalLoginForm(),
		ShowOIDCLogin:  app.oidcEnabled,
		Request: templateRequestData{
			Theme: &app.Config.Theme.themeProperties,
		},
	}
}

func newLoginTemplateTestApp(users map[string]*user, disableLocalLogin bool) *application {
	return &application{
		RequiresAuth: true,
		oidcEnabled:  true,
		Config: config{
			Auth: struct {
				SecretKey string           `yaml:"secret-key"`
				Users     map[string]*user `yaml:"users"`
				OIDC      oidcConfig       `yaml:"oidc"`
			}{
				Users: users,
				OIDC: oidcConfig{
					Issuer:            "https://example.com",
					DisableLocalLogin: disableLocalLogin,
				},
			},
		},
	}
}

func TestLoginTemplate(t *testing.T) {
	admin := map[string]*user{"admin": {PasswordHash: []byte("x")}}

	tests := []struct {
		name         string
		users        map[string]*user
		disableLocal bool
		wantUsername bool
		wantSSO      bool
	}{
		{"LocalAndOIDC", admin, false, true, true},
		{"OIDCOnly", map[string]*user{}, false, false, true},
		{"DisableLocalLogin", admin, true, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := loginPageTemplate.Execute(&buf, loginTemplateData(newLoginTemplateTestApp(tt.users, tt.disableLocal))); err != nil {
				t.Fatal(err)
			}

			out := buf.String()
			if got := strings.Contains(out, `id="username"`); got != tt.wantUsername {
				t.Fatalf("username field = %v, want %v", got, tt.wantUsername)
			}
			if got := strings.Contains(out, "SIGN IN WITH SSO"); got != tt.wantSSO {
				t.Fatalf("SSO button = %v, want %v", got, tt.wantSSO)
			}
		})
	}
}
