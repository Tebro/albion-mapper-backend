CREATE TABLE zones (
  name varchar(100) NOT NULL PRIMARY KEY,
  color varchar(20) NOT NULL,
  tier int NOT NULL
);

CREATE TABLE portals (
  id int AUTO_INCREMENT PRIMARY KEY,
  source varchar(100) NOT NULL,
  CONSTRAINT fk_source
  FOREIGN KEY (source) 
        REFERENCES zones(name),
  target varchar(100) NOT NULL,
  CONSTRAINT fk_target
  FOREIGN KEY (target) 
        REFERENCES zones(name),
  size int NOT NULL,
  expires DATETIME NOT NULL,

  UNIQUE(source, target)
);