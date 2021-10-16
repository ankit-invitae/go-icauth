package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"gopkg.in/ini.v1"
)

// Check for dependencies which are required to run icauth
func checkDependency() {
	var missingDependency []string
	for _, command := range []string{"saml2aws", "aws-iam-authenticator", "kubectl"} {
		_, err := exec.LookPath(command)
		if err != nil {
			missingDependency = append(missingDependency, command)
		}
	}

	if len(missingDependency) > 0 {
		fmt.Printf("ERROR: Missing command-line dependency: %v\n", strings.Join(missingDependency, ", "))
		os.Exit(1)
	}
}

func contains(value string, array []string) bool {
	for _, r := range array {
		if r == value {
			return true
		}
	}
	return false
}

func validateClusterName(clusterName interface{}) error {
	keys := make([]string, 0, len(EnvOktaMap))
	for k := range EnvOktaMap {
		keys = append(keys, k)
	}

	reg, _ := regexp.Compile(`^[a-z0-9\-]+\.(` + strings.Join(keys, "|") + `)$`)
	if reg.MatchString(clusterName.(string)) {
		return nil
	}
	return fmt.Errorf("%v is not a valid cluster name", clusterName)
}

func SamlProfiles() []string {
	homeDir, _ := os.UserHomeDir()
	samlFilePath := path.Join(homeDir, ".saml2aws")
	samlCfgFile, err := ini.Load(samlFilePath)

	if err != nil {
		fmt.Printf("ERROR: Cannot read saml2aws file: %v\n", samlFilePath)
		os.Exit(1)
	}
	return samlCfgFile.SectionStrings()
}
