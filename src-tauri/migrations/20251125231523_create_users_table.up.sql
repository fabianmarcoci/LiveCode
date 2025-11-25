-- migrate:up
CREATE TABLE public.users (
  id uuid NOT NULL DEFAULT gen_random_uuid(),
  username varchar(17) NOT NULL,
  email varchar(255) NOT NULL,
  password_hash text,
  is_oauth bool DEFAULT false,
  created_at timestamptz(6) DEFAULT now(),
  updated_at timestamptz(6)
);

-- Unique constraints
ALTER TABLE public.users
    ADD CONSTRAINT users_email_key UNIQUE (email);

ALTER TABLE public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);

-- Primary key
ALTER TABLE public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
