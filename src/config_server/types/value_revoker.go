package types

type ValueRevoker interface {
	Revoke(value string) error
}
