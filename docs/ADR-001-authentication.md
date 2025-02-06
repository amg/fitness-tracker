# Context

User authentication is required to store personalised data. 
For ease of signup, Google OAuth is used to verify user, create a record and issue a pair of json web tokens.

# Status

Accepted

# Decision

Using Google for sign up and sign in:
 - get initial info and email address to uniquely identify user
 - sign in will also be matched by unique email and issues tokens

2 tokens will be used: short lived JWT session token and long lived refresh token.
JWT session token will encode userId to allow fast data lookup as long as token is valid.
Refresh token will be used only to refresh session token.

Both tokens will be sent to the client. Web will store them in secure http only cookie but with different paths.
Session will be sent back with every authenticated data request but refresh token will only be sent to auth endpoint.

The signup flow is: 
 1. web will get one off auth token from Google as part of oauth flow using the React Google lib
 2. web will pass one off token to the auth endpoint to request session/refresh tokens pair and create user account if required

Additionally customer can have more than 1 active session. Each session will have unique pair of tokens and a fingerprint to loosely identify the device.

::: mermaid
sequenceDiagram
    participant Web
    participant Backend
    participant Google

    Web->>+Google: gClientId
    activate Web
    Google->>-Web: one-off auth token
    Web->>-Backend: auth token
    activate Backend
    Backend->>+Google: auth token
    activate Google
    Google->>-Backend: session/refresh tokens
    Backend->>Backend: find matching or create new user
    Backend->>Backend: issue a pair of tokens
    Backend->>Backend: persist refresh token jti and expiration
    Backend->>-Web: session/refresh http-only cookies + user name
    activate Web
    Web->>Web: store user name in localstorage
    deactivate Web
:::

### Refresh token:

When session expires or not available, refresh request is made by Web to get the session/refresh token pair again.
Refresh token in http-only cookie will be used to attempt a session refresh.

Refresh token will contain unique identifier (`jti`) to allow blacklisting in the database without storing the actual token (to avoid the need to encrypt it at rest). DB will store `jti` and `expires_at` to allow clean up in the future.

::: mermaid
sequenceDiagram
    participant Web
    participant Backend

    Web->>Web: check session store for a record
    Note over Web: no record, probably no session. See flow above
    Note over Web: have record, try getting data
    Web->>Backend: get some data
    Backend->>Web: 401 session expired
    Web->>Backend: trying to get new token (refresh_token cookie)
    Backend->>Backend: validate jwt refresh token
    Note over Backend: invalid refresh token
    Backend->>Web: failure 403
    Note over Backend: valid refresh token
    Backend->>Backend: checking refresh token jti in db
    Note over Backend: got no refresh token, sign out
    Backend->>Web: 403 sign-out
    Note over Backend: got refresh token
    Backend->>Backend: issue new tokens
    Backend->>Web: session and refresh tokens in secure http-only cookie
    Web->>Backend: get some data
    Backend->>Backend: is session token valid?
    Note over Backend: valid
    Backend->>Web: data
:::