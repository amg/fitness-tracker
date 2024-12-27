#!/bin/sh
# Inject some env variables into frontend. Careful this becomes public!
ENV_JS_FILE=./src/env.js
echo "Updating $ENV_JS_FILE ..."
echo '/* === generated env variables === */ ' > $ENV_JS_FILE
echo "window.env = { API_HOST: \"$API_HOST\" };" >> $ENV_JS_FILE

echo "Starting React web app ..."

npm start