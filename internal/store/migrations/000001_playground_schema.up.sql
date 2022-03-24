CREATE TABLE channel (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255)
);

CREATE TABLE message (
  id SERIAL PRIMARY KEY,
  channel_id INT NOT NULL,
  title TEXT,
  CONSTRAINT fk_channel FOREIGN KEY(channel_id) REFERENCES channel(id)
);