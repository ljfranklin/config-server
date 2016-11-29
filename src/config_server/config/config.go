package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/cloudfoundry/bosh-utils/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudresourcemanager "google.golang.org/api/cloudresourcemanager/v1"
	iam "google.golang.org/api/iam/v1"
)

type ServerConfig struct {
	Port                   int
	CertificateFilePath    string `json:"certificate_file_path"`
	PrivateKeyFilePath     string `json:"private_key_file_path"`
	JwtVerificationKeyPath string `json:"jwt_verification_key_path"`
	CACertificateFilePath  string `json:"ca_certificate_file_path"`
	CAPrivateKeyFilePath   string `json:"ca_private_key_file_path"`
	GCPProjectID           string `json:"gcp_project_id"`
	Store                  string
	Database               DBConfig
	GCPServiceAccounts     *iam.ProjectsServiceAccountsService
	GCPProjectManager      *cloudresourcemanager.ProjectsService
}

type DBConnectionConfig struct {
	MaxOpenConnections int `json:"max_open_connections"`
	MaxIdleConnections int `json:"max_idle_connections"`
}

type DBConfig struct {
	Adapter           string
	User              string
	Password          string
	Host              string
	Port              int
	Name              string             `json:"db_name"`
	ConnectionOptions DBConnectionConfig `json:"connection_options"`
}

func ParseConfig(filename string) (ServerConfig, error) {
	config := ServerConfig{}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, errors.WrapError(err, "Failed to read config file")
	}

	err = json.Unmarshal([]byte(data), &config)
	if err != nil {
		return config, errors.WrapError(err, "Failed to parse config file")
	}

	if config.CertificateFilePath == "" || config.PrivateKeyFilePath == "" {
		return config, errors.Error("Certificate file path and key file path should be defined")
	}

	if config.CACertificateFilePath == "" || config.CAPrivateKeyFilePath == "" {
		return config, errors.Error("CA Certificate file path and key file path should be defined")
	}

	if (&config.Database != nil) && (&config.Database.Adapter != nil) {
		config.Database.Adapter = strings.ToLower(config.Database.Adapter)
	}

	client, err := google.DefaultClient(context.Background(), iam.CloudPlatformScope)
	if err != nil {
		return config, err
	}
	service, err := iam.New(client)
	if err != nil {
		return config, err
	}
	crm, err := cloudresourcemanager.New(client)
	if err != nil {
		return config, err
	}

	config.GCPServiceAccounts = service.Projects.ServiceAccounts
	config.GCPProjectManager = crm.Projects

	return config, nil
}
