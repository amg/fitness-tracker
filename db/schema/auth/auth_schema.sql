CREATE TABLE refresh_token (
	id char(128) UNIQUE PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES user_info (id),
    fingerprint char(32) NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE (user_id, fingerprint)
);