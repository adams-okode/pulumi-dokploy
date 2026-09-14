package dokploy

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func readProjectFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(content)
}

var registryLanguages = []string{"typescript", "python", "go", "csharp", "java", "yaml"}

func TestRegistryOverviewStructure(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	require.True(t, strings.HasPrefix(index, "---\n"))
	parts := strings.SplitN(index, "---\n", 3)
	require.Len(t, parts, 3)
	frontMatter := map[string]string{}
	require.NoError(t, yaml.Unmarshal([]byte(parts[1]), &frontMatter))
	require.Equal(t, "package", frontMatter["layout"])
	require.Equal(t, "Dokploy", frontMatter["title"])
	require.Contains(t, frontMatter["meta_desc"], "Dokploy")
	require.Contains(t, frontMatter["meta_desc"], "Pulumi")
	require.NotRegexp(t, regexp.MustCompile(`(?m)^# `), parts[2])

	headings := []string{"## Installation", "## Example Usage", "## Configuration"}
	previous := -1
	for _, heading := range headings {
		position := strings.Index(index, heading)
		require.Greater(t, position, previous, heading)
		previous = position
	}
}

func TestRegistryOverviewLanguageChoosers(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	chooser := `{{< chooser language "typescript,python,go,csharp,java,yaml" >}}`
	require.Equal(t, 2, strings.Count(index, chooser))
	require.Equal(t, 2, strings.Count(index, "{{< /chooser >}}"))
	for _, language := range registryLanguages {
		open := "{{% choosable language " + language + " %}}"
		require.Equal(t, 2, strings.Count(index, open), language)
	}
	require.Equal(t, 12, strings.Count(index, "{{% /choosable %}}"))
}

func TestRegistryOverviewCoordinatesExamplesAndConfiguration(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	for _, marker := range []string{
		"npm install @dimeskigj/pulumi-dokploy",
		"pip install pulumi-dokploy",
		"go get github.com/dimeskigj/pulumi-dokploy/sdk/go/dokploy",
		"dotnet add package Dimeskigj.Pulumi.Dokploy",
		"<groupId>net.dimeski.pulumi</groupId>",
		"implementation 'net.dimeski.pulumi:dokploy:",
		"pulumi package add github.com/dimeskigj/pulumi-dokploy dokploy",
		"new dokploy.Project", "pulumi_dokploy.Project", "dokploy.NewProject",
		"new Project", "new Project(\"example\"", "type: dokploy:index:Project",
		"pulumi config set dokploy:endpoint https://dokploy.example.invalid",
		"pulumi config set --secret dokploy:apiKey your-api-key",
		"`endpoint` (Required, Not secret)", "`apiKey` (Required, Secret)",
		"DOKPLOY_ENDPOINT", "DOKPLOY_API_KEY", "community-maintained",
	} {
		require.Contains(t, index, marker)
	}
	require.NotContains(t, index, "official Dokploy")
	require.NotContains(t, index, "official Pulumi")
}

func TestRegistryDocumentationCoordinatesStayConsistent(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	installation := readProjectFile(t, "../docs/installation-configuration.md")
	for _, marker := range []string{
		"@dimeskigj/pulumi-dokploy", "pulumi-dokploy",
		"Dimeskigj.Pulumi.Dokploy", "net.dimeski.pulumi:dokploy",
		"github.com/dimeskigj/pulumi-dokploy/sdk/go/dokploy",
		"dokploy:endpoint", "dokploy:apiKey", "DOKPLOY_ENDPOINT", "DOKPLOY_API_KEY",
	} {
		require.Contains(t, index, marker)
		require.Contains(t, installation, marker)
	}
}

func TestRegistryFrontMatterAndSupportLinks(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	installation := readProjectFile(t, "../docs/installation-configuration.md")
	for name, content := range map[string]string{"index": index, "installation": installation} {
		t.Run(name, func(t *testing.T) {
			parts := strings.SplitN(content, "---\n", 3)
			require.Len(t, parts, 3)
			frontMatter := map[string]string{}
			require.NoError(t, yaml.Unmarshal([]byte(parts[1]), &frontMatter))
			require.Equal(t, "package", frontMatter["layout"])
			require.NotEmpty(t, frontMatter["meta_desc"])
			require.NotEmpty(t, frontMatter["title"])
		})
	}
	for _, link := range []string{
		"https://github.com/dimeskigj/pulumi-dokploy",
		"https://github.com/dimeskigj/pulumi-dokploy/issues",
		"https://github.com/dimeskigj/pulumi-dokploy/blob/main/CONTRIBUTING.md",
	} {
		require.Contains(t, index, link)
	}
}

func TestRegistryInstallationDetails(t *testing.T) {
	installation := readProjectFile(t, "../docs/installation-configuration.md")
	require.Regexp(t, regexp.MustCompile(`(?s)<groupId>net\.dimeski\.pulumi</groupId>\s*<artifactId>dokploy</artifactId>\s*<version>\$\{DOKPLOY_VERSION\}</version>`), installation)
	require.Contains(t, installation, "pluginDownloadURL")
	require.Contains(t, installation, "github://api.github.com/dimeskigj/pulumi-dokploy")
}

func TestRegistryExamplesDoNotExposeSecrets(t *testing.T) {
	index := readProjectFile(t, "../docs/_index.md")
	configuration := strings.Index(index, "Configure the Dokploy endpoint and API key")
	require.NotEqual(t, -1, configuration)
	examples := index[:configuration]
	require.NotContains(t, examples, "apiKey")
	require.NotContains(t, examples, "DOKPLOY_API_KEY")
	require.NotContains(t, examples, "your-api-key")
}

func TestRegistryLicense(t *testing.T) {
	license := readProjectFile(t, "../LICENSE")
	require.Contains(t, license, "TERMS AND CONDITIONS FOR USE, REPRODUCTION, AND DISTRIBUTION")
	for _, section := range []string{"1. Definitions.", "2. Grant of Copyright License.", "3. Grant of Patent License.", "4. Redistribution.", "5. Submission of Contributions.", "6. Trademarks.", "7. Disclaimer of Warranty.", "8. Limitation of Liability.", "9. Accepting Warranty or Additional Liability."} {
		require.Contains(t, license, section)
	}
	require.Contains(t, license, "APPENDIX: How to apply the Apache License to your work")
}

func TestRegistryLogoAssetAndAttribution(t *testing.T) {
	logo := readProjectFile(t, "../website/public/logo.svg")
	attribution := readProjectFile(t, "../website/public/logo-LICENSE.md")
	require.Contains(t, logo, "<svg")
	require.Contains(t, logo, "viewBox=")
	lowerLogo := strings.ToLower(logo)
	for _, forbidden := range []string{"<script", "javascript:", "http://", "https://", "data:image/", "<image", "<foreignobject"} {
		require.NotContains(t, lowerLogo, forbidden)
	}
	for _, marker := range []string{"## Source", "## License", "## Modifications", "SPDX-License-Identifier:"} {
		require.Contains(t, attribution, marker)
	}
	require.Regexp(t, regexp.MustCompile(`SPDX-License-Identifier: (MIT|Apache-2\.0|BSD-2-Clause|BSD-3-Clause|CC0-1\.0)`), attribution)
}

func TestCodegenUsesCheckedInLogoSource(t *testing.T) {
	makefile := readProjectFile(t, "../Makefile")
	mise := readProjectFile(t, "../.mise.toml")
	setupTools := readProjectFile(t, "../.github/actions/setup-tools/action.yml")
	var action map[string]any
	require.NoError(t, yaml.Unmarshal([]byte(setupTools), &action))
	require.Contains(t, makefile, "rsvg-convert -w 175 -h 175 -o sdk/dotnet/logo.png website/public/logo.svg")
	require.Contains(t, makefile, `test "$$(rsvg-convert --version | awk 'NR==1 {print $$3}')" = "2.58.0"`)
	require.Contains(t, mise, `RSVG_CONVERT_PACKAGE = "librsvg2-bin=2.58.0+dfsg-1build1"`)
	require.Contains(t, mise, `sudo apt-get install --yes ${RSVG_CONVERT_PACKAGE}`)
	require.Contains(t, mise, `test \"$(rsvg-convert --version | awk 'NR==1 {print $3}')\" = \"2.58.0\"`)
	require.Contains(t, setupTools, "run: mise run setup-svg-renderer")
	runs, ok := action["runs"].(map[string]any)
	require.True(t, ok)
	steps, ok := runs["steps"].([]any)
	require.True(t, ok)
	for _, rawStep := range steps {
		step, ok := rawStep.(map[string]any)
		if ok && step["name"] == "Setup SVG renderer" {
			require.Equal(t, "mise run setup-svg-renderer", step["run"])
			return
		}
	}
	require.Fail(t, "Setup SVG renderer step not found")
}
