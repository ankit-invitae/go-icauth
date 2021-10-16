package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Icauth struct {
	AuthRole     string
	ClusterName  string
	ForceRefresh bool
}

func (icauth Icauth) AwsProfileName() string {
	return "icauth-" + icauth.ClusterName
}

func (icauth Icauth) EnvName() string {
	clusterNameSplit := strings.Split(icauth.ClusterName, ".")
	return clusterNameSplit[len(clusterNameSplit)-1]
}

func (icauth Icauth) configureClusterAuth() error {
	// Do install configuration fof saml2aws.
	// Configuration okta login details and local aws cli profile for saml2aws to populate.
	saml2awsCfgArgs := []string{"configure", "--skip-prompt", "--idp-provider=Okta", "--idp-provider=Okta"}

	saml2awsCfgArgs = append(saml2awsCfgArgs, "--url="+EnvOktaMap[icauth.EnvName()][0])
	saml2awsCfgArgs = append(saml2awsCfgArgs, icauth.AwsProfileName())
	// this defines a distinct config block in ~/.saml2aws
	saml2awsCfgArgs = append(saml2awsCfgArgs, "--idp-account="+icauth.AwsProfileName())
	saml2awsCfgArgs = append(saml2awsCfgArgs, "--session-duration=43200")

	fmt.Printf("... Configuring authentication for %v ...\n", icauth.ClusterName)
	cmd := exec.Command("saml2aws", saml2awsCfgArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (icauth Icauth) getRoleArn() string {
	var roleArn string

	accountId := EnvOktaMap[icauth.EnvName()][1]
	if icauth.AuthRole == AwsAdminAuthRole {
		roleArn = GetAwsAdminRoleFmt(accountId)
	} else {
		roleArn = GetRoleArnFmt(accountId, icauth.ClusterName, icauth.AuthRole)
	}
	return roleArn
}

func (icauth Icauth) AwsLogin() error {
	// Create an AWS session through saml2aws as he specified IAM role
	existingLoginProfile := SamlProfiles()

	if !contains(icauth.AwsProfileName(), existingLoginProfile) || icauth.ForceRefresh {
		return icauth.configureClusterAuth()
	}

	roleArn := icauth.getRoleArn()
	saml2awsLoginArgs := []string{"login", "--idp-account=" + icauth.AwsProfileName(), "--role=" + roleArn}

	if icauth.ForceRefresh {
		saml2awsLoginArgs = append(saml2awsLoginArgs, "--force")
	}
	fmt.Printf("... Authenticating with Okta as role %v ...\n", roleArn)

	cmd := exec.Command("saml2aws", saml2awsLoginArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func (icauth Icauth) getClusterFullName() (string, error) {
	var err error
	// Get the full cluster name from the v1.0+ cluster metadata table using aws cli dynamodb functions.
	awsCliArgs := []string{
		"--profile",
		icauth.AwsProfileName(),
		"dynamodb",
		"--region",
		AwsRegion,
		"query",
		"--table-name",
		"invitae-cloud-clusters",
		"--key-condition-expression",
		"cluster_fqdn = :cluster_fqdn",
		"--expression-attribute-values",
		fmt.Sprintf(`{":cluster_fqdn": {"S": "%v.locusdev.net"}}`, icauth.ClusterName),
		"--no-scan-index-forward",
	}

	out, _ := exec.Command("/usr/local/bin/aws", awsCliArgs...).Output()

	var result Aws
	json.Unmarshal(out, &result)

	var deploymentDate string
	if result.Count == 1 {
		backendClusters := result.Items[0].Clusters.L
		if len(backendClusters) == 1 {
			deploymentDate = backendClusters[0].M.DeploymentDate.S
		} else if len(backendClusters) == 2 {
			for _, backendCluster := range backendClusters {
				if backendCluster.M.Active.Bool {
					deploymentDate = backendClusters[0].M.DeploymentDate.S
					break
				}
			}
		} else {
			err = fmt.Errorf("more than 2 backend clusters for Invitae Cloud cluster %v.locusdev.net found, something weird is going on", icauth.ClusterName)
		}
	}

	if deploymentDate == "" {
		err = fmt.Errorf("unable to determine target EKS cluster to connect to for requested cluster named %v.locusdev.net", icauth.ClusterName)
	} else {
		deploymentDate = fmt.Sprintf("%v.%v.locusdev.net", deploymentDate, icauth.ClusterName)
	}
	return deploymentDate, err
}

func (icauth Icauth) updateKubeConfig(fullClusterName string) error {
	// Update the local kubeconfig for an v1 cluster, or fail if no cluster found
	fullClusterName = strings.ReplaceAll(fullClusterName, ".", "_")
	awsCliArgs := []string{
		"--profile",
		icauth.AwsProfileName(),
		"eks",
		"--region",
		AwsRegion,
		"update-kubeconfig",
		"--name",
		fullClusterName,
		"--alias",
		icauth.ClusterName,
	}
	out, err := exec.Command("/usr/local/bin/aws", awsCliArgs...).Output()
	fmt.Println(string(out))
	return err
}

func (icauth Icauth) UpdateUserKubeConfig() error {
	// Update local kubeconfig using cluster short name and aws session profile data.
	// First retrieves the metadata (cluster full name).
	// Then attempt to retrieve eks data for the named cluster, if no eks cluster found fallback to kops-based cluster.

	fmt.Println("... Retrieving cluster metadata for cli access ...")
	fullclusterName, err := icauth.getClusterFullName()
	if err != nil {
		return err
	}
	fmt.Printf("... Configuring kubectl to access to %v ...\n", fullclusterName)
	return icauth.updateKubeConfig(fullclusterName)
}
