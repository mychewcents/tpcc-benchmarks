package udt

type Address struct {
	Street1 string `mapstructure:"street_1"`
	Street2 string `mapstructure:"street_2"`
	City    string `mapstructure:"city"`
	State   string `mapstructure:"state"`
	Zip     string `mapstructure:"zip"`
}
