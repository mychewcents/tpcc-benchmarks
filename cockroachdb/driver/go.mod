module driver

go 1.15

replace cockroachdb/model => ../model

replace cockroachdb/executors => ../executors

require (
	cockroachdb/executors v0.0.0-00010101000000-000000000000
	cockroachdb/model v0.0.0-00010101000000-000000000000
)
