#!/bin/sh
# Inject some env variables into frontend. Careful this becomes public!
ENV_JS_FILE=/app/src/env.js
echo "Updating $ENV_JS_FILE ..."
echo '/* === generated env variables === */' > $ENV_JS_FILE
echo '// eslint-disable-next-line no-undef' > $ENV_JS_FILE
echo "window.env = { API_BASE_URL: \"$API_BASE_URL\" };" >> $ENV_JS_FILE