CREATE TABLE IF NOT EXISTS users (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  username VARCHAR(255) NOT NULL,
  email VARCHAR(320) NOT NULL,
  photo_url VARCHAR(348),
  pass_hash VARCHAR(72) NOT NULL
);

CREATE UNIQUE INDEX idx_username
ON users (username);

CREATE UNIQUE INDEX idx_email
ON users (email);