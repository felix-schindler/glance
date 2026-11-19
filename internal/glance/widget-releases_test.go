package glance

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestReleaseRequestUnmarshalSource(t *testing.T) {
	tests := []struct {
		name       string
		yaml       string
		source     releaseSource
		repository string
		baseURL    string
		prerelease bool
	}{
		{"string github", `glanceapp/glance`, releaseSourceGithub, "glanceapp/glance", "", false},
		{"string gitlab", `gitlab:inkscape/inkscape`, releaseSourceGitlab, "inkscape/inkscape", "", false},
		{"string gitea", `gitea:gitea/act`, releaseSourceGitea, "gitea/act", "", false},
		{"gitlab custom url", `gitlab:https://gitlab.example.com/group/project`, releaseSourceGitlab, "group/project", "https://gitlab.example.com", false},
		{"gitea custom url with port", `gitea:https://gitea.example.com:3000/owner/repo`, releaseSourceGitea, "owner/repo", "https://gitea.example.com:3000", false},
		{"dockerhub tag", `dockerhub:nginx:stable-alpine`, releaseSourceDockerHub, "nginx:stable-alpine", "", false},
		{"object github", "repository: glanceapp/glance\ninclude-prereleases: true", releaseSourceGithub, "glanceapp/glance", "", true},
		{"object github prefix", "repository: github:glanceapp/glance\ninclude-prereleases: true", releaseSourceGithub, "glanceapp/glance", "", true},
		{"object gitlab", "repository: gitlab:inkscape/inkscape", releaseSourceGitlab, "inkscape/inkscape", "", false},
		{"object gitea custom url", "repository: gitea:https://gitea.example.com/owner/repo", releaseSourceGitea, "owner/repo", "https://gitea.example.com", false},
		{"object codeberg", "repository: codeberg:redict/redict", releaseSourceCodeberg, "redict/redict", "", false},
		{"object dockerhub tag", "repository: dockerhub:nginx:stable-alpine", releaseSourceDockerHub, "nginx:stable-alpine", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r releaseRequest
			if err := yaml.Unmarshal([]byte(tt.yaml), &r); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if r.source != tt.source {
				t.Errorf("source = %q, want %q", r.source, tt.source)
			}

			if r.Repository != tt.repository {
				t.Errorf("repository = %q, want %q", r.Repository, tt.repository)
			}

			if r.baseURL != tt.baseURL {
				t.Errorf("baseURL = %q, want %q", r.baseURL, tt.baseURL)
			}

			if r.IncludePreleases != tt.prerelease {
				t.Errorf("include-prereleases = %v, want %v", r.IncludePreleases, tt.prerelease)
			}
		})
	}
}

func TestReleaseRequestUnmarshalInvalidSource(t *testing.T) {
	inputs := []string{
		`bitbucket:owner/repo`,
		"repository: bitbucket:owner/repo",
		`github:https://github.com/owner/repo`,
		`https://gitea.com/owner/repo`,
		`gitea:`,
		`gitea:https://gitea.example.com/`,
	}
	for _, input := range inputs {
		var r releaseRequest
		if err := yaml.Unmarshal([]byte(input), &r); err == nil {
			t.Errorf("expected an error for %q", input)
		}
	}
}

func TestFetchReleaseUsesConfiguredBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		fetch    func(*releaseRequest) (*appRelease, error)
		wantPath string
		response string
	}{
		{
			name:     "gitea",
			fetch:    fetchLatestGiteaRelease,
			wantPath: "/api/v1/repos/owner/repo/releases/latest",
			response: `{"tag_name":"v1.0.0","published_at":"2024-01-02T03:04:05Z","html_url":"https://example.com/owner/repo/releases/tag/v1.0.0"}`,
		},
		{
			name:     "gitlab",
			fetch:    fetchLatestGitLabRelease,
			wantPath: "/api/v4/projects/owner%2Frepo/releases/permalink/latest",
			response: `{"tag_name":"v1.0.0","released_at":"2024-01-02T03:04:05Z","_links":{"self":"https://example.com/owner/repo/-/releases/v1.0.0"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.EscapedPath()
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()

			release, err := tt.fetch(&releaseRequest{Repository: "owner/repo", baseURL: server.URL})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotPath != tt.wantPath {
				t.Errorf("request path = %q, want %q", gotPath, tt.wantPath)
			}

			if release.Version != "v1.0.0" {
				t.Errorf("version = %q, want %q", release.Version, "v1.0.0")
			}
		})
	}
}
