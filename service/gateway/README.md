## Updating the API Schema

When updating the `api.graphql` schema file regenerate the generated code using [gqlgen](https://gqlgen.com/getting-started/):
```
cd service/gateway 
go run github.com/99designs/gqlgen generate .
```
Make sure there are no errors.
