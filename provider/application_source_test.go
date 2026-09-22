package dokploy

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	p "github.com/pulumi/pulumi-go-provider"
	"github.com/pulumi/pulumi-go-provider/infer"
	"github.com/pulumi/pulumi/sdk/v3/go/property"
	"github.com/stretchr/testify/require"
)

func TestApplicationGitSourceSavesAndClearsSSHKey(t *testing.T) {
	key := "key-1"
	for _, tc := range []struct {
		name   string
		source GitApplicationSource
		body   string
	}{
		{"set", GitApplicationSource{URL: "https://git.test/repo", Branch: "main", SSHKeyID: &key, Build: ApplicationBuild{Type: BuildNixpacks}}, `{"applicationId":"a1","customGitBranch":"main","customGitBuildPath":"","customGitSSHKeyId":"key-1","customGitUrl":"https://git.test/repo","enableSubmodules":false,"watchPaths":null}`},
		{"cleared", GitApplicationSource{URL: "https://git.test/repo", Branch: "main", Build: ApplicationBuild{Type: BuildNixpacks}}, `{"applicationId":"a1","customGitBranch":"main","customGitBuildPath":"","customGitSSHKeyId":null,"customGitUrl":"https://git.test/repo","enableSubmodules":false,"watchPaths":null}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newScriptedServer(t, expectPOST("/api/application.saveGitProvider", tc.body, `true`))
			err := configureApplicationSource(context.Background(), fixedClient(s.API())(context.Background()), "a1", ApplicationSource{Type: SourceGit, Git: &tc.source})
			require.NoError(t, err)
		})
	}
}

func TestApplicationDockerRegistrySourceRequestShape(t *testing.T) {
	password := "docker-password-sentinel"
	s := newScriptedServer(t, expectPOST("/api/application.saveDockerProvider", `{"applicationId":"application-sentinel","dockerImage":"registry.example/team/api:1","password":"docker-password-sentinel","registryUrl":"https://registry.example","username":"registry-user-sentinel"}`, `true`))
	err := configureApplicationSource(context.Background(), fixedClient(s.API())(context.Background()), "application-sentinel", ApplicationSource{Type: SourceDocker, Docker: &DockerSource{Image: "registry.example/team/api:1", RegistryURL: stringPtr("https://registry.example"), Username: stringPtr("registry-user-sentinel"), Password: &password}})
	require.NoError(t, err)
}

func TestApplicationGitLabSourceRequestShape(t *testing.T) {
	s := newScriptedServer(t, expectPOST("/api/application.saveGitlabProvider", `{"applicationId":"application-sentinel","enableSubmodules":true,"gitlabBranch":"release","gitlabBuildPath":"services/api","gitlabId":"integration-sentinel","gitlabOwner":"owner-sentinel","gitlabPathNamespace":"platform/api","gitlabProjectId":42,"gitlabRepository":"service-sentinel","watchPaths":["services/**"]}`, `true`))
	err := configureApplicationSource(context.Background(), fixedClient(s.API())(context.Background()), "application-sentinel", ApplicationSource{Type: SourceGitLab, GitLab: &GitLabAppSource{IntegrationID: "integration-sentinel", ProjectID: 42, Owner: "owner-sentinel", Namespace: "platform/api", Repository: "service-sentinel", Branch: "release", BuildPath: stringPtr("services/api"), WatchPaths: []string{"services/**"}, EnableSubmodules: true}})
	require.NoError(t, err)
}

func TestApplicationGitHubSourceRequestShape(t *testing.T) {
	s := newScriptedServer(t, expectPOST("/api/application.saveGithubProvider", `{"applicationId":"application-sentinel","branch":"release","buildPath":"services/api","enableSubmodules":true,"githubId":"integration-sentinel","owner":"owner-sentinel","repository":"service-sentinel","triggerType":"tag","watchPaths":["services/**"]}`, `true`))
	err := configureApplicationSource(context.Background(), fixedClient(s.API())(context.Background()), "application-sentinel", ApplicationSource{Type: SourceGitHub, GitHub: &GitHubAppSource{IntegrationID: "integration-sentinel", Owner: "owner-sentinel", Repository: "service-sentinel", Branch: "release", BuildPath: stringPtr("services/api"), WatchPaths: []string{"services/**"}, TriggerType: ApplicationTriggerTag, EnableSubmodules: true}})
	require.NoError(t, err)
}

// TestApplicationGitHubSourceDecodeToEncodeRoundTrip decodes a GitHub application
// exactly as application.one reports it and then re-encodes the decoded value
// through the same save call Create and Update use, so a read followed by a write
// cannot silently drop or rename a GitHub field.
func TestApplicationGitHubSourceDecodeToEncodeRoundTrip(t *testing.T) {
	const response = `{"type":"github","githubId":"integration-sentinel","owner":"owner-sentinel","repository":"service-sentinel","branch":"release","buildPath":"services/api","watchPaths":["services/**","README.md"],"triggerType":"tag","enableSubmodules":true,"buildType":"nixpacks"}`
	var raw map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(response), &raw))

	decoded, err := decodeApplicationSource(raw, ApplicationSource{})
	require.NoError(t, err)
	require.Equal(t, ApplicationSource{Type: SourceGitHub, GitHub: &GitHubAppSource{
		IntegrationID: "integration-sentinel", Owner: "owner-sentinel", Repository: "service-sentinel",
		Branch: "release", BuildPath: stringPtr("services/api"), WatchPaths: []string{"services/**", "README.md"},
		TriggerType: ApplicationTriggerTag, EnableSubmodules: true, Build: ApplicationBuild{Type: BuildNixpacks},
	}}, decoded)

	s := newScriptedServer(t, expectPOST("/api/application.saveGithubProvider", `{"applicationId":"application-sentinel","branch":"release","buildPath":"services/api","enableSubmodules":true,"githubId":"integration-sentinel","owner":"owner-sentinel","repository":"service-sentinel","triggerType":"tag","watchPaths":["services/**","README.md"]}`, `true`))
	require.NoError(t, configureApplicationSource(context.Background(), fixedClient(s.API())(context.Background()), "application-sentinel", decoded))
}

func TestApplicationGitHubSourceDecodeDefaultsTriggerTypeToPush(t *testing.T) {
	decoded, err := decodeApplicationSource(map[string]interface{}{
		"type": "github", "githubId": "i1", "owner": "owner", "repository": "repo", "branch": "main",
	}, ApplicationSource{})
	require.NoError(t, err)
	require.Equal(t, ApplicationTriggerPush, decoded.GitHub.TriggerType)
}

func TestApplicationGitHubSourceCheckDefaultsTriggerTypeToPush(t *testing.T) {
	source := property.New(map[string]property.Value{
		"type": property.New("github"),
		"github": property.New(map[string]property.Value{
			"integrationId": property.New("i1"), "owner": property.New("owner"),
			"repository": property.New("repo"), "branch": property.New("main"),
			"build": property.New(map[string]property.Value{"type": property.New("nixpacks")}),
		}),
	})
	got, err := (Application{}).Check(t.Context(), infer.CheckRequest{NewInputs: property.NewMap(map[string]property.Value{
		"name": property.New("demo"), "environmentId": property.New("e1"), "source": source,
	})})
	require.NoError(t, err)
	require.Empty(t, got.Failures)
	require.Equal(t, ApplicationTriggerPush, got.Inputs.Source.GitHub.TriggerType)
}

// TestApplicationGitHubSourceIdenticalReadProducesNoDiff guards the import path:
// feeding a read straight back into Diff must not report a change.
func TestApplicationGitHubSourceIdenticalReadProducesNoDiff(t *testing.T) {
	const response = `{"applicationId":"a1","name":"demo","environmentId":"e1","applicationStatus":"done","type":"github","githubId":"i1","owner":"owner","repository":"repo","branch":"main","buildPath":"services/api","watchPaths":["services/**"],"triggerType":"push","enableSubmodules":true,"buildType":"nixpacks"}`
	s := newScriptedServer(t, expectGET("/api/application.one", map[string][]string{"applicationId": {"a1"}}, http.StatusOK, response))
	read, err := (Application{client: fixedClient(s.API())}).Read(t.Context(), infer.ReadRequest[ApplicationArgs, ApplicationState]{ID: "a1"})
	require.NoError(t, err)
	require.Equal(t, SourceGitHub, read.Inputs.Source.Type)

	diff, err := (Application{}).Diff(t.Context(), infer.DiffRequest[ApplicationArgs, ApplicationState]{Inputs: read.Inputs, State: read.State})
	require.NoError(t, err)
	require.False(t, diff.HasChanges)
	require.Empty(t, diff.DetailedDiff)
}

func TestApplicationGitHubSourceSchemaExposesTriggerType(t *testing.T) {
	spec, err := p.GetSchema(t.Context(), Name, Version, Provider())
	require.NoError(t, err)
	require.Contains(t, schemaProperty(spec, "dokploy:index:Application", "source.github.triggerType").Type, "string")
	require.Contains(t, schemaProperty(spec, "dokploy:index:Application", "source.github.integrationId").Type, "string")
}

func TestApplicationSourceIDOnlyReadReconstructsAllFields(t *testing.T) {
	tests := []struct {
		name            string
		response        string
		want            ApplicationSource
		registryID      *string
		buildRegistryID *string
	}{
		{
			name:     "git",
			response: `{"applicationId":"a1","type":"git","registryId":"registry-git","buildRegistryId":"build-registry-git","customGitUrl":"https://git.test/repo","customGitBranch":"release","customGitBuildPath":"services/api","customGitSSHKeyId":"ssh-1","watchPaths":["services/**","README.md"],"enableSubmodules":true,"buildType":"dockerfile","dockerfile":"Containerfile","dockerContextPath":"services/api","dockerBuildStage":"production"}`,
			want: ApplicationSource{Type: SourceGit, Git: &GitApplicationSource{
				URL: "https://git.test/repo", Branch: "release", BuildPath: stringPtr("services/api"), SSHKeyID: stringPtr("ssh-1"),
				WatchPaths: []string{"services/**", "README.md"}, EnableSubmodules: true,
				Build: ApplicationBuild{Type: BuildDockerfile, Dockerfile: stringPtr("Containerfile"), DockerContextPath: stringPtr("services/api"), DockerBuildStage: stringPtr("production")},
			}},
			registryID:      stringPtr("registry-git"),
			buildRegistryID: stringPtr("build-registry-git"),
		},
		{
			name:            "docker",
			response:        `{"applicationId":"a1","type":"docker","registryId":"registry-docker","buildRegistryId":"build-registry-docker","dockerImage":"registry.test/team/api:1","registryUrl":"https://registry.test","username":"alice"}`,
			want:            ApplicationSource{Type: SourceDocker, Docker: &DockerSource{Image: "registry.test/team/api:1", RegistryURL: stringPtr("https://registry.test"), Username: stringPtr("alice")}},
			registryID:      stringPtr("registry-docker"),
			buildRegistryID: stringPtr("build-registry-docker"),
		},
		{
			name:     "gitlab",
			response: `{"applicationId":"a1","type":"gitlab","registryId":"registry-gitlab","buildRegistryId":"build-registry-gitlab","gitlabId":"integration-1","gitlabProjectId":42,"gitlabOwner":"owner","gitlabPathNamespace":"platform/api","gitlabRepository":"service","gitlabBranch":"main","gitlabBuildPath":".","watchPaths":["cmd/**"],"enableSubmodules":true,"buildType":"dockerfile","dockerfile":"Dockerfile","dockerContextPath":".","dockerBuildStage":"release"}`,
			want: ApplicationSource{Type: SourceGitLab, GitLab: &GitLabAppSource{
				IntegrationID: "integration-1", ProjectID: 42, Owner: "owner", Namespace: "platform/api", Repository: "service", Branch: "main", BuildPath: stringPtr("."),
				WatchPaths: []string{"cmd/**"}, EnableSubmodules: true,
				Build: ApplicationBuild{Type: BuildDockerfile, Dockerfile: stringPtr("Dockerfile"), DockerContextPath: stringPtr("."), DockerBuildStage: stringPtr("release")},
			}},
			registryID:      stringPtr("registry-gitlab"),
			buildRegistryID: stringPtr("build-registry-gitlab"),
		},
		{
			name:     "github",
			response: `{"applicationId":"a1","type":"github","registryId":"registry-github","buildRegistryId":"build-registry-github","githubId":"integration-1","owner":"owner","repository":"service","branch":"main","buildPath":"services/api","watchPaths":["services/**"],"triggerType":"tag","enableSubmodules":true,"buildType":"dockerfile","dockerfile":"Dockerfile","dockerContextPath":".","dockerBuildStage":"release"}`,
			want: ApplicationSource{Type: SourceGitHub, GitHub: &GitHubAppSource{
				IntegrationID: "integration-1", Owner: "owner", Repository: "service", Branch: "main", BuildPath: stringPtr("services/api"),
				WatchPaths: []string{"services/**"}, TriggerType: ApplicationTriggerTag, EnableSubmodules: true,
				Build: ApplicationBuild{Type: BuildDockerfile, Dockerfile: stringPtr("Dockerfile"), DockerContextPath: stringPtr("."), DockerBuildStage: stringPtr("release")},
			}},
			registryID:      stringPtr("registry-github"),
			buildRegistryID: stringPtr("build-registry-github"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			s := newScriptedServer(t, expectGET("/api/application.one", map[string][]string{"applicationId": {"a1"}}, http.StatusOK, tc.response))
			got, err := (Application{client: fixedClient(s.API())}).Read(t.Context(), infer.ReadRequest[ApplicationArgs, ApplicationState]{ID: "a1"})
			require.NoError(t, err)
			assertApplicationSourceFields(t, tc.want, got.Inputs.Source)
			requireLiveEqual(t, "application.registryId", tc.registryID, got.Inputs.RegistryID)
			requireLiveEqual(t, "application.buildRegistryId", tc.buildRegistryID, got.Inputs.BuildRegistryID)
		})
	}
}

func TestApplicationSourceWriteOnlySecretsArePreserved(t *testing.T) {
	password := "write-only-password"
	s := newScriptedServer(t, expectGET("/api/application.one", map[string][]string{"applicationId": {"a1"}}, http.StatusOK, `{"applicationId":"a1","type":"docker","dockerImage":"nginx:1.27","registryUrl":"https://registry.test","username":"alice"}`))
	prior := ApplicationState{ApplicationArgs: ApplicationArgs{Source: ApplicationSource{Type: SourceDocker, Docker: &DockerSource{Password: &password}}}}
	got, err := (Application{client: fixedClient(s.API())}).Read(t.Context(), infer.ReadRequest[ApplicationArgs, ApplicationState]{ID: "a1", State: prior})
	require.NoError(t, err)
	requireLiveEqual(t, "application.source.docker.password", &password, got.Inputs.Source.Docker.Password)
}

func assertApplicationSourceFields(t *testing.T, want, got ApplicationSource) {
	t.Helper()
	require.Equal(t, want.Type, got.Type)
	switch want.Type {
	case SourceGit:
		requireLiveEqual(t, "application.source.git", true, got.Git != nil)
		if got.Git == nil {
			return
		}
		requireLiveEqual(t, "application.source.git.url", want.Git.URL, got.Git.URL)
		requireLiveEqual(t, "application.source.git.branch", want.Git.Branch, got.Git.Branch)
		requireLiveEqual(t, "application.source.git.buildPath", want.Git.BuildPath, got.Git.BuildPath)
		requireLiveEqual(t, "application.source.git.sshKeyId", want.Git.SSHKeyID, got.Git.SSHKeyID)
		requireLiveEqual(t, "application.source.git.watchPaths", want.Git.WatchPaths, got.Git.WatchPaths)
		requireLiveEqual(t, "application.source.git.enableSubmodules", want.Git.EnableSubmodules, got.Git.EnableSubmodules)
		assertApplicationBuildFields(t, "application.source.git.build", want.Git.Build, got.Git.Build)
	case SourceDocker:
		requireLiveEqual(t, "application.source.docker", true, got.Docker != nil)
		if got.Docker == nil {
			return
		}
		requireLiveEqual(t, "application.source.docker.image", want.Docker.Image, got.Docker.Image)
		requireLiveEqual(t, "application.source.docker.registryUrl", want.Docker.RegistryURL, got.Docker.RegistryURL)
		requireLiveEqual(t, "application.source.docker.username", want.Docker.Username, got.Docker.Username)
	case SourceGitLab:
		requireLiveEqual(t, "application.source.gitlab", true, got.GitLab != nil)
		if got.GitLab == nil {
			return
		}
		requireLiveEqual(t, "application.source.gitlab.integrationId", want.GitLab.IntegrationID, got.GitLab.IntegrationID)
		requireLiveEqual(t, "application.source.gitlab.projectId", want.GitLab.ProjectID, got.GitLab.ProjectID)
		requireLiveEqual(t, "application.source.gitlab.owner", want.GitLab.Owner, got.GitLab.Owner)
		requireLiveEqual(t, "application.source.gitlab.namespace", want.GitLab.Namespace, got.GitLab.Namespace)
		requireLiveEqual(t, "application.source.gitlab.repository", want.GitLab.Repository, got.GitLab.Repository)
		requireLiveEqual(t, "application.source.gitlab.branch", want.GitLab.Branch, got.GitLab.Branch)
		requireLiveEqual(t, "application.source.gitlab.buildPath", want.GitLab.BuildPath, got.GitLab.BuildPath)
		requireLiveEqual(t, "application.source.gitlab.watchPaths", want.GitLab.WatchPaths, got.GitLab.WatchPaths)
		requireLiveEqual(t, "application.source.gitlab.enableSubmodules", want.GitLab.EnableSubmodules, got.GitLab.EnableSubmodules)
		assertApplicationBuildFields(t, "application.source.gitlab.build", want.GitLab.Build, got.GitLab.Build)
	case SourceGitHub:
		requireLiveEqual(t, "application.source.github", true, got.GitHub != nil)
		if got.GitHub == nil {
			return
		}
		requireLiveEqual(t, "application.source.github.integrationId", want.GitHub.IntegrationID, got.GitHub.IntegrationID)
		requireLiveEqual(t, "application.source.github.owner", want.GitHub.Owner, got.GitHub.Owner)
		requireLiveEqual(t, "application.source.github.repository", want.GitHub.Repository, got.GitHub.Repository)
		requireLiveEqual(t, "application.source.github.branch", want.GitHub.Branch, got.GitHub.Branch)
		requireLiveEqual(t, "application.source.github.buildPath", want.GitHub.BuildPath, got.GitHub.BuildPath)
		requireLiveEqual(t, "application.source.github.watchPaths", want.GitHub.WatchPaths, got.GitHub.WatchPaths)
		requireLiveEqual(t, "application.source.github.triggerType", want.GitHub.TriggerType, got.GitHub.TriggerType)
		requireLiveEqual(t, "application.source.github.enableSubmodules", want.GitHub.EnableSubmodules, got.GitHub.EnableSubmodules)
		assertApplicationBuildFields(t, "application.source.github.build", want.GitHub.Build, got.GitHub.Build)
	}
}

func assertApplicationBuildFields(t *testing.T, field string, want, got ApplicationBuild) {
	t.Helper()
	requireLiveEqual(t, field+".type", want.Type, got.Type)
	requireLiveEqual(t, field+".dockerfile", want.Dockerfile, got.Dockerfile)
	requireLiveEqual(t, field+".dockerContextPath", want.DockerContextPath, got.DockerContextPath)
	requireLiveEqual(t, field+".dockerBuildStage", want.DockerBuildStage, got.DockerBuildStage)
}

func assertLiveApplicationSource(t *testing.T, label string, want, got ApplicationSource) {
	t.Helper()
	if want.Type != got.Type {
		t.Fatalf("live Application %s source type did not match", label)
	}
	assertApplicationSourceFields(t, want, got)
}

func TestApplicationGitSourceSchemaIncludesSSHKeyID(t *testing.T) {
	spec, err := p.GetSchema(t.Context(), Name, Version, Provider())
	require.NoError(t, err)
	require.Contains(t, schemaProperty(spec, "dokploy:index:Application", "source.git.sshKeyId").Type, "string")
}

func TestApplicationSourceValidate(t *testing.T) {
	tests := []struct {
		name   string
		source ApplicationSource
		want   string
	}{
		{"docker missing", ApplicationSource{Type: SourceDocker}, "source.docker is required when source.type is docker"},
		{"docker extra", ApplicationSource{Type: SourceDocker, Docker: &DockerSource{Image: "x"}, Git: &GitApplicationSource{URL: "x", Branch: "main"}}, "source.git must be omitted when source.type is docker"},
		{"git mismatch", ApplicationSource{Type: SourceGit, Docker: &DockerSource{Image: "x"}}, "source.git is required when source.type is git"},
		{"git empty url", ApplicationSource{Type: SourceGit, Git: &GitApplicationSource{Branch: "main"}}, "source.git.url must not be empty"},
		{"git empty branch", ApplicationSource{Type: SourceGit, Git: &GitApplicationSource{URL: "https://example.test/repo"}}, "source.git.branch must not be empty"},
		{"git nixpacks docker fields", ApplicationSource{Type: SourceGit, Git: &GitApplicationSource{URL: "x", Branch: "main", Build: ApplicationBuild{Type: BuildNixpacks, Dockerfile: appStringPtr("Dockerfile")}}}, "source.git.build.dockerfile must be omitted for nixpacks builds"},
		{"git dockerfile missing", ApplicationSource{Type: SourceGit, Git: &GitApplicationSource{URL: "x", Branch: "main", Build: ApplicationBuild{Type: BuildDockerfile}}}, "source.git.build.dockerfile is required for dockerfile builds"},
		{"valid docker", ApplicationSource{Type: SourceDocker, Docker: &DockerSource{Image: "nginx"}}, ""},
		{"valid git", ApplicationSource{Type: SourceGit, Git: &GitApplicationSource{URL: "x", Branch: "main", Build: ApplicationBuild{Type: BuildNixpacks}}}, ""},
		{"valid gitlab", ApplicationSource{Type: SourceGitLab, GitLab: &GitLabAppSource{IntegrationID: "i", ProjectID: 1, Owner: "o", Namespace: "n", Repository: "r", Branch: "main", Build: ApplicationBuild{Type: BuildDockerfile, Dockerfile: appStringPtr("Dockerfile")}}}, ""},
		{"github missing", ApplicationSource{Type: SourceGitHub}, "source.github is required when source.type is github"},
		{"github extra docker", ApplicationSource{Type: SourceGitHub, GitHub: validGitHub(), Docker: &DockerSource{Image: "x"}}, "source.docker must be omitted when source.type is github"},
		{"github extra git", ApplicationSource{Type: SourceGitHub, GitHub: validGitHub(), Git: &GitApplicationSource{URL: "x", Branch: "main"}}, "source.git must be omitted when source.type is github"},
		{"github extra gitlab", ApplicationSource{Type: SourceGitHub, GitHub: validGitHub(), GitLab: &GitLabAppSource{IntegrationID: "i"}}, "source.gitlab must be omitted when source.type is github"},
		{"github empty integrationId", ApplicationSource{Type: SourceGitHub, GitHub: withGitHub(func(s *GitHubAppSource) { s.IntegrationID = "" })}, "source.github.integrationId must not be empty"},
		{"github empty owner", ApplicationSource{Type: SourceGitHub, GitHub: withGitHub(func(s *GitHubAppSource) { s.Owner = "" })}, "source.github.owner must not be empty"},
		{"github empty repository", ApplicationSource{Type: SourceGitHub, GitHub: withGitHub(func(s *GitHubAppSource) { s.Repository = "" })}, "source.github.repository must not be empty"},
		{"github empty branch", ApplicationSource{Type: SourceGitHub, GitHub: withGitHub(func(s *GitHubAppSource) { s.Branch = "" })}, "source.github.branch must not be empty"},
		{"github invalid triggerType", ApplicationSource{Type: SourceGitHub, GitHub: withGitHub(func(s *GitHubAppSource) { s.TriggerType = "release" })}, "source.github.triggerType must be one of push or tag"},
		{"github empty triggerType", ApplicationSource{Type: SourceGitHub, GitHub: withGitHub(func(s *GitHubAppSource) { s.TriggerType = "" })}, "source.github.triggerType must be one of push or tag"},
		{"github dockerfile missing", ApplicationSource{Type: SourceGitHub, GitHub: withGitHub(func(s *GitHubAppSource) { s.Build = ApplicationBuild{Type: BuildDockerfile} })}, "source.github.build.dockerfile is required for dockerfile builds"},
		{"github valid push", ApplicationSource{Type: SourceGitHub, GitHub: validGitHub()}, ""},
		{"github valid tag", ApplicationSource{Type: SourceGitHub, GitHub: withGitHub(func(s *GitHubAppSource) { s.TriggerType = ApplicationTriggerTag })}, ""},
		{"unknown type", ApplicationSource{Type: "bitbucket", Git: &GitApplicationSource{URL: "x", Branch: "main"}}, "source.type must be one of docker, git, gitlab, or github"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.validate()
			if tt.want == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.want)
			}
		})
	}
}

func appStringPtr(s string) *string { return &s }

func validGitHub() *GitHubAppSource {
	return &GitHubAppSource{IntegrationID: "i", Owner: "o", Repository: "r", Branch: "main", TriggerType: ApplicationTriggerPush, Build: ApplicationBuild{Type: BuildNixpacks}}
}

func withGitHub(mutate func(*GitHubAppSource)) *GitHubAppSource {
	s := validGitHub()
	mutate(s)
	return s
}
