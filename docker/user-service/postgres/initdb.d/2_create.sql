\c userservice;
CREATE SCHEMA userschema;

CREATE TABLE userschema.users (
              id SERIAL,
              name varchar(200) DEFAULT NULL,
              age int DEFAULT NULL,
              mail varchar(200) DEFAULT NULL,
              address varchar(1024) DEFAULT NULL,
              PRIMARY KEY (id)
);

ALTER SCHEMA userschema OWNER TO user_users;
ALTER TABLE userschema.users OWNER TO user_users;

