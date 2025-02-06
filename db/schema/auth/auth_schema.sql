CREATE TABLE refresh_token_jti (
	id uuid UNIQUE PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES user_info (id),
    fingerprint VARCHAR(64) NOT NULL,
    expires_at TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP
);