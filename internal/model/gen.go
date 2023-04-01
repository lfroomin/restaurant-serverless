/*
  To generate the model/restaurant.gen.go file from the OAS3 spec:
	  1. Install oapi-codegen (one time installation)
				'go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest'
      2. Generate model from internal/model folder
				'go generate'
*/

package model

//go:generate oapi-codegen --config config.yaml restaurant-api.yaml
