package types

import (
	"config_server/config"
	"encoding/base64"
	"errors"
	"fmt"

	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
)

type serviceAccountGenerator struct {
	config config.ServerConfig
}

func NewServiceAccountGenerator(config config.ServerConfig) ValueGenerator {
	return serviceAccountGenerator{config: config}
}

func (s serviceAccountGenerator) Generate(parameters interface{}) (interface{}, error) {
	params := parameters.(map[string]interface{})
	accountID, ok := params["account_id"].(string)
	if !ok {
		return nil, errors.New("account_id is required")
	}

	accountDisplayName, ok := params["account_display_name"].(string)
	if !ok {
		accountDisplayName = accountID
	}

	rolesTemp, ok := params["roles"].([]interface{})
	if !ok {
		return nil, errors.New("roles are required")
	}
	roles := []string{}
	for _, r := range rolesTemp {
		roleName, ok := r.(string)
		if !ok {
			return nil, errors.New("Role name not of type string")
		}
		roles = append(roles, roleName)
	}

	projectPath := fmt.Sprintf("projects/%s", s.config.GCPProjectID)
	serviceAcct, err := s.config.GCPServiceAccounts.Create(projectPath, &iam.CreateServiceAccountRequest{
		AccountId: accountID,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: accountDisplayName,
		},
	}).Do()
	if err != nil {
		return nil, err
	}

	// Create keys
	serviceKey, err := s.config.GCPServiceAccounts.Keys.Create(serviceAcct.Name, &iam.CreateServiceAccountKeyRequest{}).Do()
	if err != nil {
		return nil, err
	}
	keyContent, err := base64.StdEncoding.DecodeString(serviceKey.PrivateKeyData)
	if err != nil {
		return nil, err
	}

	policy, err := s.config.GCPProjectManager.GetIamPolicy(s.config.GCPProjectID, &cloudresourcemanager.GetIamPolicyRequest{}).Do()
	if err != nil {
		return nil, err
	}

	for _, r := range roles {
		policy.Bindings = append(policy.Bindings, &cloudresourcemanager.Binding{
			Role:    r,
			Members: []string{fmt.Sprintf("serviceAccount:%s", serviceAcct.Email)},
		})
	}

	_, err = s.config.GCPProjectManager.SetIamPolicy(s.config.GCPProjectID, &cloudresourcemanager.SetIamPolicyRequest{
		Policy: policy,
	}).Do()
	if err != nil {
		return nil, err
	}

	return string(keyContent), nil
}
