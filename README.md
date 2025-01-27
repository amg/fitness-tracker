- [Context](#context)
- [Requirements](#requirements)
- [Technology](#technology)
  - [High level view:](#high-level-view)
- [Building](#building)
  - [Run on docker](#run-on-docker)
  - [Run on GCP](#run-on-gcp)

## Context

This project is meant as a playground for learning Full-Stack development but it is my hope that it will continue to evolve into something useful one day.

It is dockerized and ready for deployment in GCP.
See docker-compose.yml and /deploy Terraform files.

Terraform has currently some sections commented out to speed up iteration but for initial run, they are required.

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
  - (nice) end to end encryption using google account
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

1. React JS for FE (frontend)
2. Go lang for BE (backend)
3. Google Cloud (infrastructure)
4. Docker

<br />

### High level view:
![High level view](./docs/high-level-arch.png)

## Building<a name="building"></a>

Main entry point `docker-compose.yaml`

4 services:
 - api
    GO lang api service
 - node-api
    Typescript express based service (only local docker at this point)
 - frontend
    ReactJS website
 - postres DB

`.secrets` folder is required to specify multiple variables used by the stack but not committed to the source code. Can copy starting point from `secrets-example` and fill in the blanks.

### Run on docker<a name="run-docker"></a>

Script to build. See script for details:<br/>
`./build.sh -a build -e dev -p arm`

Script to watch:<br/>
`./build.sh -a watch -e dev`

Script to log:<br/>
`./build.sh -a logs -e dev`

### Run on GCP<a name="run-gcp"></a>

IMPORTANT: gcloud doesn't like images built on M1 so have to use `-p amd` 

1. (if required)`gcloud auth configure-docker`
2. Script to build and push to registry. See script for details:<br/>
`./build.sh -a build -e staging -p amd -u true`
1. `./deploy_main.sh -m apply` (yes before core, see below)
2. (if required) `./deploy_core.sh -m apply`

NOTE:
1. DB takes 15+ min to create
2. DNS records and loadbalancer certificate provisioning can take 24 hours

Becanse of that don't recreate those when actively devving.

IMPORTANT: Extra steps to make this run.
Since switching from experimental cloud run domain mapping to external Load Balancer few steps became manual:
1. once LB is up, declare 2 A records in Cloud DNS for api and web endpoints pointing to that LB
2. Load balancer mapping is not specified in Terraform, edit load balancer and add routes manually 
Edit classic application load balancer -> Host and path rules (path /*, host [api/web].domain.you.own.com)
3. Wait for a while, check certificate provisioning process, will be a link in Frontend of LB (up to 24 hours)


Inspect manifest: `docker manifest inspect gcr.io/learning-gcloud-444623/web:latest`
