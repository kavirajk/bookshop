CREATE TABLE IF NOT EXISTS users (
	id    	              SERIAL PRIMARY KEY,
	first_name            CHARACTER VARYING(255)         NOT NULL,
	last_name             CHARACTER VARYING(255)         NOT NULL,
	email                 CHARACTER VARYING(255)         NOT NULL,
	username              CHARACTER VARYING(255)         NOT NULL,
	is_active             BOOLEAN                        NOT NULL,
	password_hash         CHARACTER VARYING(255)        NOT NULL,
	created_at            TIMESTAMP WITHOUT TIME ZONE     NOT NULL,
	updated_at            TIMESTAMP WITHOUT TIME ZONE     NOT NULL
);
