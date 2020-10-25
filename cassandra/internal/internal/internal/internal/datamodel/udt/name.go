package udt

type Name struct {
	FirstName  string `cql:"first_name" mapstructure:"first_name"`
	MiddleName string `cql:"middle_name" mapstructure:"middle_name"`
	LastName   string `cql:"last_name" mapstructure:"last_name"`
}
