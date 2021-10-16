package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
)

func getProps() Icauth {
	var qs = []*survey.Question{
		{
			Name: "authRole",
			Prompt: &survey.Select{
				Message: "Please select Role:",
				Options: RoleList,
			},
		},
		{
			Name:     "clusterName",
			Prompt:   &survey.Input{Message: "Please enter ClusterName:"},
			Validate: validateClusterName,
		},
	}

	var icauth Icauth
	err := survey.Ask(qs, &icauth)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	return icauth
}

func main() {
	forceRefresh := flag.Bool("f", false, "force refresh of the cached configuration and login session")
	flag.Parse()

	// Check if all the dependencies are present or not
	checkDependency()

	icauth := getProps()
	icauth.ForceRefresh = *forceRefresh

	err := icauth.AwsLogin()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}

	err = icauth.UpdateUserKubeConfig()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
