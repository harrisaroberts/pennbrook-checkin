CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL
);

INSERT INTO users (username, password_hash, role) VALUES
  (
    'admin',
    '$2a$10$LL89B01E4hFcPDYCSJoEwuG88nnslBXJcq4oUlCl2bR11TEEg5A3a',
    'admin'
  ),
  (
    'guard',
    '$2a$10$wr/i22eJCfNs.ufUN48as.D3MzU.ijFHXrqU8Xk.kblrRv8Gif4IK',
    'guard'
  );

