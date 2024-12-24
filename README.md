(MVP) Fitness Workout Tracker

Web:

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


Technology:

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