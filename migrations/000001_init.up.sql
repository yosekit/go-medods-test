CREATE TABLE IF NOT EXISTS refresh_tokens (
    id serial PRIMARY KEY,
    user_id uuid NOT NULL,
    client_ip character varying(45) NOT NULL,
    hashed_token text NOT NULL,
    created_at timestamp without time zone DEFAULT now()
);
