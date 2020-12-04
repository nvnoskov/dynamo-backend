CREATE TABLE flight
(
    id            VARCHAR PRIMARY KEY,
    name          VARCHAR NOT NULL,
    number        VARCHAR NOT NULL,
    departure     VARCHAR NOT NULL,
    departure_time TIMESTAMP NOT NULL,
    destination   VARCHAR NOT NULL,
    arrival_time   TIMESTAMP NOT NULL,
    fare          VARCHAR NOT NULL,
    duration      VARCHAR NOT NULL,
    created_at    TIMESTAMP NOT NULL,
    updated_at    TIMESTAMP NOT NULL
);