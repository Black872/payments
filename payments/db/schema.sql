CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL UNIQUE,
  name VARCHAR NOT NULL,
  email VARCHAR NOT NULL UNIQUE,
  password_hash BYTEA NOT NULL,
  is_active BOOL DEFAULT true
);

CREATE TABLE IF NOT EXISTS admins (
	id SERIAL UNIQUE,
  user_id INT REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS accounts (
  id BIGSERIAL UNIQUE, 
  name BYTEA NOT NULL UNIQUE, /* IBAN */
  holder_id INT REFERENCES users (id),
  balance NUMERIC DEFAULT 0.00,
  is_active BOOL DEFAULT true
);

CREATE TYPE status AS ENUM ('prepared', 'completed');

CREATE TABLE IF NOT EXISTS payments (
  id BIGSERIAL NOT NULL,
  sender_id INT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE CASCADE,
  recipient_id INT REFERENCES accounts (id) ON UPDATE CASCADE ON DELETE CASCADE,
  amount NUMERIC NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE,
  completed_at TIMESTAMP WITH TIME ZONE,
  status status DEFAULT 'prepared'
);

CREATE TYPE category AS ENUM ('block_card', 'unblock_card', 'create_user', 'block_user', 'unblock_user', 'balance');

CREATE TABLE IF NOT EXISTS logs (
  id BIGSERIAL,
  time TIMESTAMP WITH TIME ZONE, 
  category category,
  user_id INT REFERENCES users (id),
  description TEXT
);