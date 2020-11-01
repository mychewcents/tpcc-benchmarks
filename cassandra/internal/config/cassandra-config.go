package config

type CassandraConfig struct {
	Hosts            []string `xml:"Hosts>Host"`
	WriteConsistency string   `xml:"WriteConsistency"`
	ReadConsistency  string   `xml:"ReadConsistency"`
}
