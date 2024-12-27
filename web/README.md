Simple dockerised react app

## Building

`docker-compose.yaml` declares the services and `.env.development` (and others) specify some env variables.

2 services:
 - api
    GO lang api service
 - frontend
    ReactJS website

Run command below to recreate container:
`docker compose --env-file ./.env.development build --progress=plain --no-cache && docker compose --env-file ./.env.development up`

Run existing container:
`docker compose --env-file ./.env.development up`