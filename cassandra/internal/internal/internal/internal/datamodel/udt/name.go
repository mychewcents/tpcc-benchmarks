package udt

type Name struct {
	FirstName  string `mapstructure:"first_name"`
	MiddleName string `mapstructure:"middle_name"`
	LastName   string `mapstructure:"last_name"`
}
