/*
 * Databricks
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package models

type ServicePrincipal struct {
	DisplayName   string         `json:"displayName,omitempty"`
	Groups        []Groups       `json:"groups,omitempty"`
	Id            string         `json:"id,omitempty"`
	Entitlements  []Entitlements `json:"entitlements,omitempty"`
	ApplicationId string         `json:"applicationId,omitempty"`
	Active        bool           `json:"active,omitempty"`
}
