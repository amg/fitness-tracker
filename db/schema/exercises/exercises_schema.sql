CREATE TABLE exercise (
	id BIGSERIAL UNIQUE PRIMARY KEY,
	user_id uuid REFERENCES user_info (id),
    name VARCHAR(50) NOT NULL,
	description text NOT NULL
);