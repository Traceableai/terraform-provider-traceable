package provider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTeamsSamlConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTeamsSamlConfigCreate,
		Read:   resourceTeamsSamlConfigRead,
		Update: resourceTeamsSamlConfigUpdate,
		Delete: resourceTeamsSamlConfigDelete,

		Schema: map[string]*schema.Schema{
			"group_attribute_name": {
				Type:        schema.TypeString,
				Description: "The name of the group attribute",
				Required:    true,
			},
			"group_to_roles_collection": {
				Type:        schema.TypeList,
				Description: "Collection of groups to roles mapping",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:        schema.TypeString,
							Description: "The name of the group",
							Required:    true,
						},
						"assigned_roles": {
							Type:        schema.TypeList,
							Description: "List of assigned roles",
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Description: "The role ID",
										Required:    true,
									},
									"scope": {
										Type:        schema.TypeList,
										Description: "The scope of the role",
										Optional:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"entity_scopes": {
													Type:        schema.TypeList,
													Description: "Entity scopes",
													Required:    true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Description: "The type of the scope",
																Required:    true,
															},
															"ids": {
																Type:        schema.TypeList,
																Description: "The IDs of the scope",
																Required:    true,
																Elem: &schema.Schema{
																	Type: schema.TypeString,
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
func resourceTeamsSamlConfigCreate(d *schema.ResourceData, meta interface{}) error {
	groupAttributeName := d.Get("group_attribute_name").(string)

	groupToRolesCollection := d.Get("group_to_roles_collection").([]interface{})
	groupToRolesCollectionStr := buildGroupToRolesCollection(groupToRolesCollection)

	query := fmt.Sprintf(`mutation {
		createSamlGroupMapping(
			input: {
				groupAttributeName: "%s",
				groupToRolesCollection: %s
			}
		) {
			id
			groupAttributeName
		}
	}`, groupAttributeName, groupToRolesCollectionStr)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling GraphQL response: %s", err)
	}

	if response["data"] != nil && response["data"].(map[string]interface{})["createSamlGroupMapping"] != nil {
		id := response["data"].(map[string]interface{})["createSamlGroupMapping"].(map[string]interface{})["id"].(string)
		d.SetId(id)
	} else {
		return fmt.Errorf("could not create SAML group mapping, no ID returned")
	}

	return nil
}

func buildGroupToRolesCollection(groups []interface{}) string {
	groupStr := "["
	for _, group := range groups {
		groupMap := group.(map[string]interface{})
		groupName := groupMap["group_name"].(string)
		assignedRoles := groupMap["assigned_roles"].([]interface{})
		assignedRolesStr := buildAssignedRoles(assignedRoles)
		groupStr += fmt.Sprintf(`{
			groupName: "%s",
			assignedRoles: %s
		},`, groupName, assignedRolesStr)
	}
	groupStr = groupStr[:len(groupStr)-1] // Remove trailing comma
	groupStr += "]"
	return groupStr
}

func buildAssignedRoles(roles []interface{}) string {
	roleStr := "["
	for _, role := range roles {
		roleMap := role.(map[string]interface{})
		roleId := roleMap["id"].(string)
		scope := roleMap["scope"].([]interface{})
		scopeStr := buildScope(scope)
		roleStr += fmt.Sprintf(`{
			id: "%s",
			scope: %s
		},`, roleId, scopeStr)
	}
	roleStr = roleStr[:len(roleStr)-1] // Remove trailing comma
	roleStr += "]"
	return roleStr
}

func buildScope(scopes []interface{}) string {
	if len(scopes) == 0 {
		return "null"
	}
	scopeStr := "{ entityScopes: ["
	for _, scope := range scopes {
		scopeMap := scope.(map[string]interface{})
		scopeType := scopeMap["type"].(string)
		ids := scopeMap["ids"].([]interface{})
		idsStr := buildStringArray(interfaceSliceToStringSlice(ids))
		scopeStr += fmt.Sprintf(`{
			type: "%s",
			ids: %s
		},`, scopeType, idsStr)
	}
	scopeStr = scopeStr[:len(scopeStr)-1] // Remove trailing comma
	scopeStr += "]}"
	return scopeStr
}
func resourceTeamsSamlConfigRead(d *schema.ResourceData, meta interface{}) error {
	readQuery := `{
		samlGroupMappingsMetadata {
			count
			results {
				id
				groupToRolesCollection {
					groupName
					assignedRoles {
						roleId: id
						scope {
							entityScopes {
								type
								ids
							}
						}
					}
				}
				groupAttributeName
			}
			total
		}
	}`

	var response map[string]interface{}
	responseStr, err := executeQuery(readQuery, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling GraphQL response: %s", err)
	}

	id := d.Id()
	ruleDetails := getRuleDetailsFromRulesListUsingIdName(response, "samlGroupMappingsMetadata", id, "id", "groupAttributeName")
	if ruleDetails == nil {
		d.SetId("")
		return nil
	}

	d.Set("group_attribute_name", ruleDetails["groupAttributeName"])

	groupToRolesCollection := ruleDetails["groupToRolesCollection"].([]interface{})
	var groupToRolesList []map[string]interface{}
	for _, group := range groupToRolesCollection {
		groupMap := group.(map[string]interface{})
		groupName := groupMap["groupName"].(string)
		assignedRoles := groupMap["assignedRoles"].([]interface{})
		var assignedRolesList []map[string]interface{}
		for _, role := range assignedRoles {
			roleMap := role.(map[string]interface{})
			roleId := roleMap["roleId"].(string)
			scope := roleMap["scope"].(map[string]interface{})
			var scopeList []map[string]interface{}
			if entityScopes, ok := scope["entityScopes"].([]interface{}); ok {
				for _, entityScope := range entityScopes {
					entityScopeMap := entityScope.(map[string]interface{})
					scopeType := entityScopeMap["type"].(string)
					ids := entityScopeMap["ids"].([]interface{})
					scopeList = append(scopeList, map[string]interface{}{
						"type": scopeType,
						"ids":  ids,
					})
				}
			}
			assignedRolesList = append(assignedRolesList, map[string]interface{}{
				"id":    roleId,
				"scope": scopeList,
			})
		}
		groupToRolesList = append(groupToRolesList, map[string]interface{}{
			"group_name":     groupName,
			"assigned_roles": assignedRolesList,
		})
	}
	d.Set("group_to_roles_collection", groupToRolesList)

	return nil
}
func resourceTeamsSamlConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	groupAttributeName := d.Get("group_attribute_name").(string)

	groupToRolesCollection := d.Get("group_to_roles_collection").([]interface{})
	groupToRolesCollectionStr := buildGroupToRolesCollection(groupToRolesCollection)

	query := fmt.Sprintf(`mutation {
		updateSamlGroupMapping(
			input: {
				id: "%s",
				groupAttributeName: "%s",
				groupToRolesCollection: %s
			}
		) {
			id
			groupAttributeName
		}
	}`, id, groupAttributeName, groupToRolesCollectionStr)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling GraphQL response: %s", err)
	}

	if response["data"] != nil && response["data"].(map[string]interface{})["updateSamlGroupMapping"] != nil {
		d.SetId(id)
	} else {
		return fmt.Errorf("could not update SAML group mapping, no ID returned")
	}

	return nil
}
func resourceTeamsSamlConfigDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()

	query := fmt.Sprintf(`mutation {
		deleteSamlGroupMapping(id: "%s")
	}`, id)

	var response map[string]interface{}
	responseStr, err := executeQuery(query, meta)
	if err != nil {
		return fmt.Errorf("Error while executing GraphQL query: %s", err)
	}
	err = json.Unmarshal([]byte(responseStr), &response)
	if err != nil {
		return fmt.Errorf("Error while unmarshalling GraphQL response: %s", err)
	}

	if success, ok := response["data"].(map[string]interface{})["deleteSamlGroupMapping"].(bool); !ok || !success {
		return fmt.Errorf("failed to delete SAML group mapping")
	}

	d.SetId("")
	return nil
}
func buildStringArray(input []string) string {
	if len(input) == 0 {
		return "[]"
	}
	output := "["
	for _, v := range input {
		output += fmt.Sprintf(`"%s",`, v)
	}
	output = output[:len(output)-1] // Remove trailing comma
	output += "]"
	return output
}

func interfaceSliceToStringSlice(input []interface{}) []string {
	var output []string
	for _, v := range input {
		output = append(output, v.(string))
	}
	return output
}
