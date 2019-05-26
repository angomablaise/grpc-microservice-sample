DROP DATABASE IF EXISTS userservice;
CREATE DATABASE userservice;
USE userservice;

CREATE TABLE `users` (
              `id` bigint(20) NOT NULL AUTO_INCREMENT,
              `name` varchar(200) DEFAULT NULL,
              `age` bigint(20) DEFAULT NULL,
              `mail` varchar(200) DEFAULT NULL,
              `address` varchar(1024) DEFAULT NULL,
              PRIMARY KEY (`ID`),
              UNIQUE KEY `ID_UNIQUE` (`ID`)
);

CREATE USER `user-users`@`%` IDENTIFIED BY 'password';
GRANT SELECT,INSERT,UPDATE,DELETE ON userservice.* TO `user-users`@`%`;
