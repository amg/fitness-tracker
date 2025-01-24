DROP TABLE IF EXISTS user_info;

EXEC SQL INCLUDE '../schema/user_info_schema.sql'

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

INSERT INTO user_info (external_id, email, first_name, last_name, dob)
VALUES 
	(uuid_generate_v4(), 'vi_kiramman@gmail.com', 'Vi', 'Kiramman', '967-12-19'),
	(uuid_generate_v4(), 'jinx@gmail.com', 'Jinx', 'Unknown', '972-10-10')
RETURNING *

-- CREATE UNIQUE INDEX external_id_idx ON user_info (external_id);