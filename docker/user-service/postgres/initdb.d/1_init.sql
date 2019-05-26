CREATE USER user_users WITH PASSWORD 'password';

DROP DATABASE IF EXISTS userservice;
CREATE DATABASE userservice;
GRANT ALL PRIVILEGES ON DATABASE userservice TO user_users
