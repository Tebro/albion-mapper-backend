
CREATE TABLE portals (
  id int AUTO_INCREMENT PRIMARY KEY,
  source varchar(100) NOT NULL,
  target varchar(100) NOT NULL,
  size int NOT NULL,
  expires DATETIME NOT NULL,

  UNIQUE(source, target)
);