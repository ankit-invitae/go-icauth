package main

import "fmt"

const AwsRegion = "us-east-1"

// login urls from okta web dashboard
const NonPrdOktaURL = "https://invitae.okta.com/home/amazon_aws/0oaqemgk28W8REabg0x7/272"
const PrdOktaURL = "https://invitae.okta.com/home/amazon_aws/0oaqj47083Kepm7Iv0x7/272"
const ItOktaURL = "https://invitae.okta.com/home/amazon_aws/0oarivbraqsPI3w9V0x7/272"

const NonPrdAcctId = "160990826323"
const PrdAcctId = "527246062414"
const ItAcctId = "761885676921"

var EnvOktaMap = map[string][]string{
	"dev-test": {NonPrdOktaURL, NonPrdAcctId},
	"dev":      {NonPrdOktaURL, NonPrdAcctId},
	"tst":      {NonPrdOktaURL, NonPrdAcctId},
	"stg":      {NonPrdOktaURL, NonPrdAcctId},
	"prd":      {PrdOktaURL, PrdAcctId},
	"it":       {ItOktaURL, ItAcctId},
}

const AwsAdminAuthRole = "aws-admin"
const AdminAuthRole = "admin"
const DeployerAuthRole = "deployer"
const ViewerAuthRole = "viewer"

var RoleList = []string{AwsAdminAuthRole, AdminAuthRole, DeployerAuthRole, ViewerAuthRole}

func GetAwsAdminRoleFmt(accountId string) string {
	return fmt.Sprintf("arn:aws:iam::%v:role/Administrator-SSO", accountId)
}

func GetRoleArnFmt(accountId, clusterName, authRole string) string {
	return fmt.Sprintf("arn:aws:iam::%v:role/%v.locusdev.net-%v-SSO", accountId, clusterName, authRole)
}
