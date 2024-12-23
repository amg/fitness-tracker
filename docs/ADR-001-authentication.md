# Context

User authentication is required to store personalised data. 
For ease of signup, Google OAuth is used.

# Status

Accepted

# Decision

Web will get one off auth token from Google using the React Google lib.
Web will pass one off token to the GoLang backend to request session/refresh tokens pair.
Refresh token will be stored in the HTTP-only cookie while session token will be returned to the web client to persist in a session storage for easier access.

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
    Backend->>-Web: refresh http-only cookie + session token
    activate Web
    deactivate Web
:::

When session expires or not available, refresh request is made by Web to get the session/refresh token pair again.
Refresh token in http-only cookie will be used to attempt a session refresh.

::: mermaid
sequenceDiagram
    participant Web
    participant Backend
    participant Google

    Web->>Web: check session store
    Note over Web: no access token
    Web->>Backend: request tokens
    Backend->>Backend: is refresh token in a cookie
    Note over Backend: got refresh token
    Backend->>Google: get new pair of tokens
    Google->>Google: validate refresh token
    Note over Google: if valid
    Google->>Backend: fresh session/refresh tokens
    Backend->>Web: refresh http-only cookie + session token
    Note over Google: if not valid
    Google->>Backend: not authorized
    Backend->>Web: logged out
:::