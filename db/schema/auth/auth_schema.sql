CREATE TABLE refresh_token (
	id VARCHAR(1024) UNIQUE PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES user_info (id),
    fingerprint VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE (user_id, fingerprint)
);