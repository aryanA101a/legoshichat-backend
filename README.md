# legoshichat-backend
![legoshichat](https://github.com/aryanA101a/legoshichat-backend/assets/23309033/51add021-5c59-42cf-be1e-d2bdedd04b20)

Welcome to the legoshichat API, the official interface for the legoshichat. Unlock seamless integration and access the unique features of legoshichat through this official API.

## Documentation
Check out documentation [here](https://aryana101a.github.io/legoshichat-backend/).

## Roadmap
- [x] API Spec
- [ ] API Implementation **WIP**
- [ ] Unit Testing
- [ ] CI
- [ ] CD

## Steps To Run

1. `docker run --name gofr-pgsql  -e POSTGRES_DB=legoshichat -e POSTGRES_USER=legoshiuser -e POSTGRES_PASSWORD=legoshipass -p 2006:5432 -d postgres:latest`

2. `go run .`

**To monitor database**
`docker exec -it gofr-pgsql psql --username=legoshiuser --dbname=legoshichat`



