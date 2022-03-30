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

CREATE TABLE replie (
  id SERIAL PRIMARY KEY,
  message_id INT NOT NULL,
  title TEXT,
  CONSTRAINT fk_message FOREIGN KEY(message_id) REFERENCES message(id)
);
