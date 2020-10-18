package udt

import "fmt"

type Name struct {
	FirstName  string `mapstructure:"first_name"`
	MiddleName string `mapstructure:"middle_name"`
	LastName   string `mapstructure:"last_name"`
}

func (n *Name) GetNameString() string {
	return fmt.Sprintf("\"{first_name:'%s', middle_name:'%s', last_name:'%s'}\"", n.FirstName, n.MiddleName, n.LastName)
}
