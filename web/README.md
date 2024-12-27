Simple dockerised react app

## Building

`docker-compose.yaml` declares the services and `.env.development` (and others) specify some env variables.

2 services:
 - api
    GO lang api service
 - frontend
    ReactJS website

Run command below to recreate container and pass in args to build script so they are accessible in Dockerfile.
For some reason `docker compose` doesn't auto pass `.env` to Dockerfiles declared.
`env_file` declaration passes env to the container itself BUT not to the build scripts.

Command:
`. ./.env.development && docker compose build --build-arg API_HOST="$API_HOST" --progress=plain --no-cache && docker compose --env-file ./.env.development up`

Run existing container:
`docker compose --env-file ./.env.development up`