CREATE TABLE Users(
  
  user_id varchar(200) NOT NULL,
  username varchar(100) UNIQUE NOT NULL,
  email varchar(200) NOT NULL,
  password varchar(100) NOT NULL,
  primary key(user_id)
);