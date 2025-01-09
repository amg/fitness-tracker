- [Requirements](#requirements)
- [Technology](#technology)
  - [Docs:](#docs)
- [Building](#building)
  - [Run on docker](#run-on-docker)
  - [Run on GCP](#run-on-gcp)
- [Other useful links](#other-useful-links)


## Requirements<a name="reqs"></a>

Customer authentication:
  - (must) signup with google
  - (must) login with google
  - (must) get profile
  - (nice) update nickname

Exercises input:
  - (must) create
      - (must) name
      - (must) simple description
      - (must) image
      - (nice) video
  - (must) delete
  - (nice) edit
  - (nice to have) end to end encryption using google account
      https://stackoverflow.com/questions/41939884/server-side-google-sign-in-way-to-encrypt-decrypt-data-with-google-managed-secr
      https://cloud.google.com/docs/security/key-management-deep-dive
  - (nice) exercises edit

Schedule builder:
  - (must) create new daily schedule
      ie. every x days, can be every second day for example
  - (must) set reps and sets goal (3 sets 10 reps each)
  - (must) finish schedule/end it/archive so it's remembered
  - (nice) timed schedule
      ie. start, end on the date
  - (nice) notifications for a workout
  - (nice) add to google calendar (web)


## Technology<a name="technology"></a>

1. React JS (https://www.googlecloudcommunity.com/gc/Community-Blogs/No-servers-no-problem-A-guide-to-deploying-your-React/ba-p/690760)
 - install node using brew
 - install npm
 - create react app

2. Go lang for backend
3. Google Cloud Run

Authentication:
https://developers.google.com/identity/gsi/web/guides/overview
(chrome only)https://developers.google.com/privacy-sandbox/cookies/fedcm

 ### Docs:

ReactJS
  - https://react.dev/learn/state-as-a-snapshot
OAuth
  - https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
OAuth Go lang
  - https://github.com/golang-jwt/jwt?tab=readme-ov-file


## Building<a name="building"></a>

`docker-compose.yaml` declares the services and `.env.development` (and others) specify some env variables.

2 services:
 - api
    GO lang api service
 - frontend
    ReactJS website

### Run on docker<a name="run-docker"></a>

Run command below to recreate container and run:

`eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' ./.secrets/.env.development) && docker compose watch`

// doesn't rebuild but prints logs as they come

`eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' ./.secrets/.env.development) && docker compose up --watch`

// runs printing logs

`... up logs`


### Run on GCP<a name="run-gcp"></a>

IMPORTANT: gcloud dosen't like images built on M1 so have to use `buildx bake` instead

`eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' ./.secrets/.env.staging) && docker buildx bake`

use ` --print` for dry run

Pushing the image to google registry:
https://cloud.google.com/artifact-registry/docs/docker/pushing-and-pulling

1. (if required)`./deploy_core.sh -a`
2. (from `deploy/`) `. ./.secrets/.env.staging && docker compose push`
3. (if required) `./deploy_core.sh -m apply`
4. `./deploy_main.sh -m apply`


Inspect manifest: `docker manifest inspect gcr.io/learning-gcloud-444623/web:latest`

## Other useful links<a name="links"></a>

Load balancer the hard way:
https://cloud.google.com/blog/topics/developers-practitioners/serverless-load-balancing-terraform-hard-way

Region picker:
https://googlecloudplatform.github.io/region-picker/
https://cloud.google.com/dns/docs/zones

Testing DNS propagation:
https://www.whatsmydns.net/#NS

Debugging GOlang with docker container:
https://blog.jetbrains.com/go/2020/05/06/debugging-a-go-application-inside-a-docker-container/

Adding PostgreSQL:
https://blog.logrocket.com/building-simple-app-go-postgresql/