package types

import "config_server/config"

type ValueRevokerConcrete struct {
	config config.ServerConfig
}

type NoopRevoker struct{}

func (NoopRevoker) Revoke(string) error { return nil }

func NewValueRevokerConcrete(config config.ServerConfig) ValueRevokerConcrete {
	return ValueRevokerConcrete{config: config}
}

func (vgc ValueRevokerConcrete) GetRevoker(valueType string) (ValueRevoker, error) {
	switch valueType {
	case "gcp_service_account":
		return NewServiceAccountRevoker(vgc.config), nil
	default:
		return NoopRevoker{}, nil
	}
}
