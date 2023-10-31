module github.com/kozmod/omni/client

go 1.21

replace github.com/kozmod/omni/external => ../external

require (
	github.com/kozmod/omni/external v0.0.0-00010101000000-000000000000
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
