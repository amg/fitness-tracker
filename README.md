- [Context](#context)
- [Pillars](#pillars)
  - [Simplicity](#simplicity)
  - [Reliability](#reliability)
- [Rough MVP](#rough-mvp)
- [Technology](#technology)
  - [High level view](#high-level-view)
  - [Swagger spec](#swagger-spec)
- [Building](#building)
  - [Run on docker](#run-on-docker)
  - [Run on GCP](#run-on-gcp)

## Context

Do you love exercises? Do you record what you did and when? Does it look like a bunch of lines on the paper note book and you wish you could somehow summarise all the amazing effort you have put in? Or have you used an app just to realise that exercise you are doing doesn't match the standard set or that its input is hundred clicks away?

Meet Fitness tracker, a Fullstack distributed application that aims to eliminate the friction of recording and reviewing your exercises!
It’s being built using a number of languages and technologies: Go, Typescript, HCL Terraform and Docker. 
This project is about learning about backend development in a real-world setting.


## Pillars<a name="reqs"></a>

1. simplicity
2. reliability

### Simplicity

Doing exercises is hard enough, no more traversing multiple screens trying to record what you have achieved...

Once you have created a simple account recording is as simple as:
1. open the website (you will have shortcut on your mobile most likely!)
2. today is preselected, you see the list of generic exercises or your custom ones
3. you select one, you specify number of sets and reps for each and you hit done!

### Reliability

Data is safely stored in the cloud. Any device, any time, same simple experience.

## Rough MVP

Customer authentication:
  - (must) signup/login with google
  - (must) account info

Exercises input:
  - create/edit
    - (must) name, description
    - (must) visual steps
    - (nice-to-have) video
  - delete
    - (must) archive, existing recordings are safe

Schedule/Reminders:
  - (must) create
      ie. every x days, can be every second day for example
  - (must) delete
  - (nice) timed schedule
      ie. start, end on the date
  - (nice) notifications

## Technology<a name="technology"></a>

1. React JS (frontend)
2. GO lang for auth (backend)
3. NodeJS for data (backend)
4. Terraform Google Cloud (infrastructure)
5. Docker containers

<br />

### High level view
![High level view](./docs/high-level-arch.png)

### Swagger spec

[API spec](./docs/openapi-spec.yaml)

## Building<a name="building"></a>

Main entry point `docker-compose.yaml`

4 services:
 - api
    GO lang auth api service
 - node-api
    Typescript node express
 - frontend
    ReactJS single page app
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

Because of that don't recreate those when actively devving.

IMPORTANT: Extra steps to make this run.
Since switching from experimental cloud run domain mapping to external Load Balancer few steps became manual:
1. once LB is up, declare 2 A records in Cloud DNS for api and web endpoints pointing to that LB
2. Load balancer mapping is not specified in Terraform, edit load balancer and add routes manually 
Edit classic application load balancer -> Host and path rules (path /*, host [api/web].domain.you.own.com)
3. Wait for a while, check certificate provisioning process, will be a link in Frontend of LB (up to 24 hours)
