package domain

type Route struct {
	Id           int64
	HttpMethod   string
	HttpAddress  string
	Instance     string
	ProtoService string
	ProtoMethod  string
	AccessRole   AccessRole
	IsActive     bool
}

type RoutesRepository interface {
	All() ([]*Route, error)
	GetByAddress(addr string) (*Route, error)
	Insert(*Route) error
	Delete(int64) error
}
