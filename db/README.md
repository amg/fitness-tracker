Postgres:
`docker pull postgres:latest`

`docker volume create postgres-volume`

`docker run --name pgsql-dev --env POSTGRES_PASSWORD=test1234! --volume postgres-volume:/var/lib/postgresql/data --publish 5432:5432 postgres`

Admin:
`docker pull dpage/pgadmin4`

`docker rm /pgsql-admin4-dev && docker run --name pgsql-admin4-dev --env PGADMIN_DEFAULT_EMAIL=gvozdik.sanya@gmail.com --env PGADMIN_DEFAULT_PASSWORD=test1234! -p 9090:80 dpage/pgadmin4`

Inspect to find network ip address (required?)
`docker inspect pgsql-dev`

Server add, password as above and username `postgres`

Unix socks connection for postgres on gcloud run + sql auth proxy
`TF_VAR_postgres_fullstring='host=/cloudsql/learning-gcloud-444623:us-west1:postgres-instance-ft-staging user=%s password=%s dbname=%s sslmode=disable'`


Refs:
 - https://medium.com/google-cloud/managing-cloud-sql-resources-with-terraform-76cc044319e9
 - https://squaredup.com/blog/running-postgres-in-docker/
 - https://cloud.google.com/sql/docs/mysql/connect-run#terraform_1