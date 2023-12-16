# legoshichat-backend
![legoshichat](https://github.com/aryanA101a/legoshichat-backend/assets/23309033/51add021-5c59-42cf-be1e-d2bdedd04b20)

Welcome to the legoshichat API, the official interface for the legoshichat. Unlock seamless integration and access the unique features of legoshichat through this official API.

This API employs JSON Web Tokens (JWT) for authentication. Please note that all endpoints, with the exception of the login and create account functionalities, necessitate the inclusion of an authorized JWT bearer token for access.

## Documentation
Check out documentation [here](https://aryana101a.github.io/legoshichat-backend/).

Also checkout the **postman collection** in the [docs](https://github.com/aryanA101a/legoshichat-backend/tree/main/docs) folder.

## Roadmap
- [x] API Spec
- [x] API Implementation
- [ ] Unit Testing **30%**
- [ ] CI
- [ ] CD

## Steps To Run

1. `docker run --name gofr-pgsql  -e POSTGRES_DB=legoshichat -e POSTGRES_USER=legoshiuser -e POSTGRES_PASSWORD=legoshipass -p 2006:5432 -d postgres:latest`

2. `go run .`

**To monitor database**
`docker exec -it gofr-pgsql psql --username=legoshiuser --dbname=legoshichat`


**Note:** In the context of this project, the .env file is intentionally exposed publicly; however, it is crucial to emphasize that this practice is strongly discouraged in a production environment

