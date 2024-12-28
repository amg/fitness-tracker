# (MVP) Fitness Workout Tracker

## Web:

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


### Technology:

React JS (https://www.googlecloudcommunity.com/gc/Community-Blogs/No-servers-no-problem-A-guide-to-deploying-your-React/ba-p/690760)
 - install node using brew
 - install npm
 - create react app
Google Cloud Run
Go lang for backend

Authentication:
    https://developers.google.com/identity/gsi/web/guides/overview
    (chrome only)https://developers.google.com/privacy-sandbox/cookies/fedcm


Future considerations:
 - look at using next.js

 Docs:
    ReactJS
     - https://react.dev/learn/state-as-a-snapshot
    OAuth
     - https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
    OAuth Go lang
     - https://github.com/golang-jwt/jwt?tab=readme-ov-file


### Building

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

--------------

IMPORTANT: gcloud dosen't like images built on M1 so have to use `buildx bake` instead
`eval $(sed -e '/^#/d' -e 's/^/export /' -e 's/$/;/' .env.development) && docker buildx bake`
use ` --print` for dry run

Pushing the image to google registry:
https://cloud.google.com/artifact-registry/docs/docker/pushing-and-pulling

`gcloud auth configure-docker`
`docker compose --env-file ./.env.development push`

Inspect manifest: `docker manifest inspect gcr.io/learning-gcloud-444623/web:latest`


//what's this about? https://www.reddit.com/r/docker/comments/13wgqgz/how_to_specify_provenance_with_docker_compose/

https://medium.com/@francisihe/how-to-get-google-cloud-run-service-url-programmatically-72964e2ce344


### Env example

```
ENV_PATH=./.env.development

ENV=dev-

# used for CORS
WEB_BASE_URL=http://127.0.0.1:3000
WEB_PORT=3000

# used for JWT cookie
COOKIE_DOMAIN=127.0.0.1

# used to make API calls from ReactJS
API_BASE_URL=http://127.0.0.1:8080
API_PORT=8080

# Google oauth
GOOGLE_PROJECT_ID=<ID>
GOOGLE_CLIENT_ID=<CLIENT_ID>
GOOGLE_CLIENT_SECRET=<CLIENT_SECRET>
GOOGLE_CLIENT_CALLBACK_URL=http://127.0.0.1:3000

# JWT token keys
FILE_KEY_PRIVATE=jwtRSA256-private.pem
FILE_KEY_PUBLIC=jwtRSA256-public.pem
```