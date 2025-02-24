DROP TABLE IF EXISTS user_info;

CREATE TABLE user_info (
	id uuid UNIQUE PRIMARY KEY,
    email VARCHAR(100) UNIQUE NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    picture_url VARCHAR(2048)
);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO user_info (id, email, first_name, last_name)
VALUES 
	(uuid_generate_v4(), 'vi_kiramman@gmail.com', 'Vi', 'Kiramman'),
	(uuid_generate_v4(), 'jinx@gmail.com', 'Jinx', 'Unknown')
RETURNING *;