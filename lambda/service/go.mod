module github.com/pennsieve/drs-service/service

//TODO update module name

go 1.21

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)


replace github.com/pennsieve/drs-service/service => ./service
