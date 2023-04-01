restaurant-serverless is a microservice using a REST API
deployed in a serverless environment.

The API is documented using the OAS3 (Swagger) specification,
and the model is generated from the specification. The
specification is in the model/restaurant-api.yaml file.

The basic CRUD endpoints exist for the restaurant entity.
- Create - create a restaurant
- Read - get a restaurant
- Update - update a restaurant
- Delete - delete a restaurant

When a restaurant is created or updated, if it contains
an address, the address is used to look up the geocode
coordinates of the address (lat, lon).

The AWS services used:
- API Gateway
- Lambda functions
- Dynamo DB
- Location (used for geocoding)

A SAM (Serverless Application Model) template is used to organize
the service and deploy it to AWS.

To update the generated model when the OAS3 specification is
changed, do the following:
- From the internal/model folder execute `go generate`

**Build**
- From the project root folder execute `sam build`

**Deploy**
- From the project root folder execute `sam deploy`

**Unit Tests**
- From the project root folder execute `go test ./...`
