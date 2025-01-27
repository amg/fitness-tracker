package utils

import "time"

const KEY_SESSION_TOKEN = "session_token"
const KEY_REFRESH_TOKEN = "refresh_token"
const SESSION_EXPIRATION_TIME = 30 * time.Minute
const REFRESH_TOKEN_EXPIRATION_TIME = 30 * 24 * time.Hour // 1 month
