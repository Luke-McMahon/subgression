DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id INT AUTO_INCREMENT NOT NULL,
  fullname VARCHAR(128) NOT NULL,
  beltrank VARCHAR(128) NOT NULL,
  degree INT,
  PRIMARY KEY (`id`)
);

INSERT INTO users
(fullname, beltrank, degree)
VALUES
('Luke McMahon', 'White', 4),
('Wayne Spence', 'Brown', 2),
('Josh Rogers', 'White', 4),
('Sam Jones', 'Blue', 2),
('John Morton', 'Purple', 2)