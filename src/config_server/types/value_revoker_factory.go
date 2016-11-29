package types

type ValueRevokerFactory interface {
	GetRevoker(valueType string) (ValueRevoker, error)
}
