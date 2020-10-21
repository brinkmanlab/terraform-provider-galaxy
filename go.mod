module terraform-provider-galaxy

go 1.15

replace (
	github.com/brinkmanlab/blend4go => ../blend4go
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/brinkmanlab/blend4go v0.2.0
	github.com/fatih/color v1.9.0 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/go-multierror v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/hcl/v2 v2.6.0 // indirect
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.4
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/mapstructure v1.3.3
	github.com/oklog/run v1.1.0 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	github.com/zclconf/go-cty v1.6.1 // indirect
	golang.org/x/net v0.0.0-20200925080053-05aa5d4ee321 // indirect
	golang.org/x/sys v0.0.0-20200923182605-d9f96fdee20d // indirect
	google.golang.org/genproto v0.0.0-20200925023002-c2d885f95484 // indirect
	google.golang.org/grpc v1.32.0 // indirect
)
