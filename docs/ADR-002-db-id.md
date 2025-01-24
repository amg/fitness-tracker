# Context

DB tables related decisions grouped in one file.

# Status

Accepted

# Decisions

1. userId will be uuid and will be a PK

Main reason is that user identity is included in the signed JWT token to avoid trip to db everytime to validate customer. 

 - having auto-incremented int `id` would expose that to customer as JWT can be decoded. Not ideal as people would know the count
 - using 2 ids: internal int and external uuid would mean data requests will need to fetch internal id from users table first to reference data

So while uuid is bigger than standard `id` it solves a few problems.
