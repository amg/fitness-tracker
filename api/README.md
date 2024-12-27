Generate token using:
>openssl genrsa -out jwtRSA256-private.pem 2048

Get public token out of private:
>openssl rsa -in jwtRSA256-private.pem -pubout -outform PEM -out jwtRSA256-public.pem

NOTE: these tokens are test only as they are committed to the repository.
 For production environment keep private key secret.

Refs:
 https://techdocs.akamai.com/iot-token-access-control/docs/generate-rsa-keys
