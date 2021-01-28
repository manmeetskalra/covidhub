CREATE DATABASE covidhub;
CREATE TABLE covidhub.users (
    id int NOT NULL AUTO_INCREMENT,
    lastname varchar(255) NOT NULL,
    firstname varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    password varchar(255) NOT NULL,
    phoneNumber varchar(255) NOT NULL,
    loggedIn boolean NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE covidhub.emails (
    email varchar(255) NOT NULL,
    frequency int NOT NULL,
    dueTime DATETIME(6),
    countries VARCHAR(1000),
    informationType VARCHAR(255),
    PRIMARY KEY (email)
); 

CREATE TABLE covidhub.texts (
    phoneNumber varchar(255) NOT NULL,
    frequency int NOT NULL,
    dueTime DATETIME,
    country VARCHAR(255),
    PRIMARY KEY (phoneNumber)
); 

INSERT INTO covidhub.users (lastname, firstname, email, password, phoneNumber, loggedIn) VALUES ('manan', 'maniyar', 'mm@gmail.com', 'xyz', '1111111111', false);
INSERT INTO covidhub.users (lastname, firstname, email, password, phoneNumber, loggedIn) VALUES ('michael', 'zaghi', 'mz@gmail.com', 'xyz', '2222222222', false);
INSERT INTO covidhub.users (lastname, firstname, email, password, phoneNumber, loggedIn) VALUES ('raghav', 'mittal', 'rm@gmail.com', 'xyz', '3333333333', false);
INSERT INTO covidhub.emails (email, frequency, countries, informationType) VALUES('mm@gmail.com', 300, "India|Canada", "Confirmed");
INSERT INTO covidhub.emails (email, frequency, countries, informationType) VALUES('rm@gmail.com', 300, "India|Canada", "Confirmed");
INSERT INTO covidhub.emails (email, frequency, countries, informationType) VALUES('mz@gmail.com', 300, "India|Canada", "Confirmed");
INSERT INTO covidhub.texts (phoneNumber, frequency, country) VALUES('1111111111', 300, "India");
INSERT INTO covidhub.texts (phoneNumber, frequency, country) VALUES('2222222222', 600, "India");
INSERT INTO covidhub.texts (phoneNumber, frequency, country) VALUES('3333333333', 900, "India");
