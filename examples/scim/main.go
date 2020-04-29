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

	//printServicePrincipals(listServicePrincipals(endpoint))
	//printUserGroups(listGroups(endpoint))
	//printCreatedGroups(createGroup(endpoint))
	//printUsers(listUsers(endpoint))


}

func printUser(user *models.ScimUser) {
	fmt.Println("getting user")
	fmt.Println("user created with Id:",user.Id)
	fmt.Println("user created with name",user.DisplayName)
	for index, element := range user.Groups {
		fmt.Println(index, "User added to Group with displayName: =>", element.Display)
		fmt.Println(index, "User added to Group with Id:=>", element.Value)
	}
}


func printUsers(users *models.ListUserRequestScim) {
	fmt.Println("getting users")
	for index, element := range users.Resources {
		fmt.Println(index, "displayName: =>", element.DisplayName)
		fmt.Println(index, "User Family Name: =>", element.Name.FamilyName)
	}
}

func printGetGroup(group *models.ScimGroup) {
	fmt.Println("GroupId:",group.Id)
	fmt.Println("GroupDisplayName:",group.DisplayName)
	for index, element := range group.Members {
		fmt.Println(index, "displayName: =>", element.Display)
		fmt.Println(index, "id: =>", element.Value)
	}
}

func getServicePrincipal(endpoint scim.Endpoint, id string) *models.ServicePrincipal {
	fmt.Println("Getting Service Principals for id:", id)
	resp, err := endpoint.GetServicePrincipal(id)
	fmt.Println("Response:\n", resp)
	if err != nil {
		panic(err)
	}
	return resp
}

func deleteServicePrincipal(endpoint scim.Endpoint, id string) {
	fmt.Println("Deleting Service Principals for id:", id)
	err := endpoint.DeleteServicePrincipal(id)
	if err != nil {
		panic(err)
	}
}

func createServicePrincipal(endpoint scim.Endpoint, appId string, displayName string) *models.ServicePrincipal {
	fmt.Println("Attaching Service Principals for id:", appId)
	req := models.ServicePrincipalCreateRequest{}
	req.ApplicationId = appId
	req.DisplayName = displayName
	resp, err := endpoint.CreateServicePrincipal(&req)
	if err != nil {
		panic(err)
	}
	return resp
}

func updateServicePrincipal(endpoint scim.Endpoint, id string, updatedSP models.ServicePrincipal) *models.ServicePrincipal {
	fmt.Println("Updating Service Principals for SP id:", id)
	resp, err := endpoint.UpdateServicePrincipal(&updatedSP)
	if err != nil {
		panic(err)
	}
	return resp
}

func printServicePrincipals(principals []models.ServicePrincipal) {
	fmt.Println(principals)
}

func printUserGroups(groups []models.ScimGroup) {
	for index, element := range groups {
		fmt.Println(index, "displayName: =>", element.DisplayName)
		fmt.Println(index, "id: =>", element.Id)
		members := element.Members
		for _,element2 := range members {
			fmt.Println("value id: =>", element2.Value)
			fmt.Println("displayName with Member: =>", element2.Display)
		}

	}
	fmt.Println(groups)
}

func printCreatedGroups(displayName string, members []models.ScimMember) {
	fmt.Println("group displayName :",displayName)
	for index, element := range members {
		fmt.Println(index, "Member id: =>", element.Value)
		fmt.Println(index, "Member displayNAme: =>", element.Display)
		fmt.Println(index, "Member ref: =>", element.Ref)
	}
}

func listServicePrincipals(endpoint scim.Endpoint) []models.ServicePrincipal {
	fmt.Println("Listing Service Principals")
	resp, err := endpoint.ListServicePrincipal()
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

func listGroups(endpoint scim.Endpoint) []models.ScimGroup {
	fmt.Println("Listing User Groups")
	resp, err := endpoint.ListGroups()
	if err != nil {
		panic(err)
	}
	return resp.Resources
}

func createGroup(endpoint scim.Endpoint,memberId string) (string,[]models.ScimMember) {
	fmt.Println("Creating User Group")
	member := models.ScimMember{
		Display: "",
		Value:   memberId,
		Ref:     "",
	}
	group := models.ScimGroup{
		Entitlements: nil,
		DisplayName:  "blah6",
		Members:      []models.ScimMember{member},
		Groups:       nil,
		Id:           "",
	}
	resp, err := endpoint.CreateGroup(&group)
	if err != nil {
		panic(err)
	}
	return resp.DisplayName,resp.Members
}

func getGroup(endpoint scim.Endpoint,id string) *models.ScimGroup {
	fmt.Println("Getting User Group with id:",id)
	resp, err := endpoint.GetGroup(id)
	if err != nil {
		panic(err)
	}
	return resp
}

func deleteGroup(endpoint scim.Endpoint,id string) {
	fmt.Println("Deleting User Group with id:",id)
	err := endpoint.DeleteGroup(id)
	if err != nil {
		panic(err)
	}
}

func updateGroup(endpoint scim.Endpoint,id string,memberId string,memberName string,groupName string) *models.ScimGroup {
	fmt.Println("Updating User Group with id:",id)
	fmt.Println("Creating User Group")

	updateMemberRef := fmt.Sprintf("https://eastus2.azuredatabricks.net/api/2.0/scim/v2/ServicePrincipals/%s", memberId)
	member := models.ScimMember{
		Display: memberName,
		Value:   memberId,
		Ref:     updateMemberRef,
	}

	entitlement1 := models.Entitlements{Value: "allow-cluster-create"}
	entitlement2 := models.Entitlements{Value: "allow-instance-pool-create"}
	group := models.ScimGroup{
		Entitlements: []models.Entitlements{entitlement1,entitlement2},
		DisplayName:  groupName,
		Members:      []models.ScimMember{member},
		Groups:       nil,
		Id:           id,
	}


	resp,err := endpoint.UpdateGroup(id,group)
	if err != nil {
		panic(err)
	}
	return resp
}

func listUsers(endpoint scim.Endpoint) *models.ListUserRequestScim {
	fmt.Println("Listing Users ")
	resp, err := endpoint.ListUsers()
	if err != nil {
		panic(err)
	}
	return resp
}

func createUser(endpoint scim.Endpoint,userName string,group string,) *models.ScimUser {
	fmt.Println("Creating Users ")

	entitlements := models.Entitlements{Value: "allow-cluster-create"}
	groups := models.Groups{
		Display: "",
		Value:   group,
		Ref:     "",
	}

	name := models.Name{
		FamilyName: "",
		GivenName:  userName,
	}
	user := models.ScimUser{
		Entitlements: []models.Entitlements{entitlements},
		DisplayName:  userName,
		Groups:       []models.Groups{groups},
		Emails:       nil,
		Id:           "",
		Name:         &name,
		Active:       false,
		UserName:     userName,
	}
	resp, err := endpoint.CreateUser(user)
	if err != nil {
		panic(err)
	}
	return resp
}

func getUser(endpoint scim.Endpoint,id string) *models.ScimUser {
	fmt.Println("Getting User  with id:",id)
	resp, err := endpoint.GetUser(id)
	if err != nil {
		panic(err)
	}
	return resp
}

func deleteUser(endpoint scim.Endpoint,id string)  {
	fmt.Println("Getting User  with id:",id)
	err := endpoint.DeleteUser(id)
	if err != nil {
		panic(err)
	}
}

func updateUser(endpoint scim.Endpoint,id string,group string,userName string) *models.ScimUser  {
	fmt.Println("Getting User  with id:",id)
	entitlement1 := models.Entitlements{Value: "allow-cluster-create"}
	entitlement2 := models.Entitlements{Value: "allow-instance-pool-create"}
	groups := models.Groups{
		Display: "",
		Value:   group,
		Ref:     "",
	}

	name := models.Name{
		FamilyName: "",
		GivenName:  userName,
	}
	user := models.ScimUser{
		Entitlements: []models.Entitlements{entitlement1,entitlement2},
		DisplayName:  userName,
		Groups:       []models.Groups{groups},
		Emails:       nil,
		Id:           "",
		Name:         &name,
		Active:       false,
		UserName:     userName,
	}
	resp,err := endpoint.UpdateUser(id,user)
	if err != nil {
		panic(err)
	}
	return resp
}