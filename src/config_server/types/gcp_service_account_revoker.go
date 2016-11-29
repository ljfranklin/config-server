package types

import (
	"config_server/config"
	"encoding/json"
	"fmt"
)

type serviceAccountRevoker struct {
	config config.ServerConfig
}

func NewServiceAccountRevoker(config config.ServerConfig) ValueRevoker {
	return serviceAccountRevoker{config: config}
}

func (s serviceAccountRevoker) Revoke(value string) error {
	var keyFields map[string]string
	err := json.Unmarshal([]byte(value), &keyFields)
	if err != nil {
		return err
	}

	projectID := keyFields["project_id"]
	clientEmail := keyFields["client_email"]
	serviceAccountName := fmt.Sprintf("projects/%s/serviceAccounts/%s", projectID, clientEmail)

	_, err = s.config.GCPServiceAccounts.Delete(serviceAccountName).Do()
	if err != nil {
		return err
	}

	return nil
}
