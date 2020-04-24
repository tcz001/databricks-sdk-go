package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tcz001/databricks-sdk-go/api/scim"
	"github.com/tcz001/databricks-sdk-go/client"
	"github.com/tcz001/databricks-sdk-go/models"
	"io/ioutil"
)


func main() {
	flag.Parse() // required to suppress warnings from glog
	
	secrets := loadSecrets()

	cl, err := client.NewClient(client.Options{
		Domain:                              &secrets.Domain,
		Token:                               &secrets.Token,
		XDatabricksAzureWorkspaceResourceId: &secrets.DBWorkspaceResourceId,
		XDatabricksAzureSPManagementToken:   &secrets.SPMgmtToken,
	})
	if err != nil {
		panic(err)
	}


	endpoint := scim.Endpoint{
		Client: cl,
	}

	printServicePrincipals(listServicePrincipals(endpoint))

}

func printServicePrincipals(principals []models.ServicePrincipal) {
	fmt.Println(principals)
}

func listServicePrincipals(endpoint scim.Endpoint) []models.ServicePrincipal {
	fmt.Println("Listing Service Principals")
	resp,err := endpoint.List()
	if err != nil {
		panic(err)
	}
	return resp.Resources
}

type secrets struct {
	Domain                string `json:"domain"`
	Token                 string `json:"token"`
	ClusterName           string `json:"cluster_name"`
	SPToken               string `json:"sp_token"`
	SPMgmtToken           string `json:"sp_mgmt_token"`
	DBWorkspaceResourceId string `json:"workspace_resource_id"`
}

func loadSecrets() *secrets {
	content, err := ioutil.ReadFile("secrets.json")
	if err != nil {
		panic(err)
	}

	var sc secrets
	err = json.Unmarshal(content, &sc)
	if err != nil {
		panic(err)
	}

	return &sc
}
