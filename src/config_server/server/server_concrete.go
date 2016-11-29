package server

import (
	"config_server/config"
	"config_server/store"
	"config_server/types"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cloudfoundry/bosh-utils/errors"
)

type configServer struct {
	config                config.ServerConfig
	valueGeneratorFactory types.ValueGeneratorFactory
	valueRevokerFactory   types.ValueRevokerFactory
}

func NewConfigServer(config config.ServerConfig) ConfigServer {
	return configServer{config: config}
}

func (cs configServer) Start() error {
	if err := cs.configureHandler(); err != nil {
		return err
	}

	return http.ListenAndServeTLS(":"+strconv.Itoa(cs.config.Port),
		cs.config.CertificateFilePath,
		cs.config.PrivateKeyFilePath, nil)
}

func (cs configServer) configureHandler() error {
	jwtTokenValidator, err := NewJwtTokenValidator(cs.config.JwtVerificationKeyPath)
	if err != nil {
		return errors.WrapError(err, "Failed to create JWT token validator")
	}

	store, err := store.CreateStore(cs.config)
	if err != nil {
		return errors.WrapError(err, "Failed to create data store")
	}

	cs.valueGeneratorFactory = types.NewValueGeneratorConcrete(cs.config)
	cs.valueRevokerFactory = types.NewValueRevokerConcrete(cs.config)
	requestHandler, err := NewRequestHandler(
		store,
		cs.valueGeneratorFactory,
		cs.valueRevokerFactory,
	)
	if err != nil {
		return errors.WrapError(err, "Failed to create Request Handler")
	}
	authenticationHandler := NewAuthenticationHandler(jwtTokenValidator, requestHandler)

	http.Handle("/v1/data", authenticationHandler)
	http.Handle("/v1/data/", authenticationHandler)

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ = range ticker.C {
			fmt.Print("Checking for expired items...\n")
			expiredItems, err := store.GetExpired()
			if err != nil {
				fmt.Printf("Failed getting expired items: %s\n", err.Error())
			}
			fmt.Printf("Found %d expired items\n", len(expiredItems))
			for _, item := range expiredItems {
				fmt.Printf("Revoking item with ID: %s...\n", item.ID)
				err := cs.revoke(item)
				if err != nil {
					fmt.Printf("Failed revoking item: %s\n", err.Error())
				}

				_, err = store.Delete(item.Name)
				if err != nil {
					fmt.Printf("Failed deleting item: %s\n", err.Error())
				}
			}
		}
	}()

	return nil
}

func (cs configServer) revoke(value store.Configuration) error {
	var storedConfig store.Configuration
	configStr, err := value.StringifiedJSON()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(configStr), &storedConfig)
	if err != nil {
		return err
	}

	revoker, err := cs.valueRevokerFactory.GetRevoker(storedConfig.Type)
	if err != nil {
		return err
	}

	err = revoker.Revoke(storedConfig.Value)
	if err != nil {
		return err
	}

	return nil
}
