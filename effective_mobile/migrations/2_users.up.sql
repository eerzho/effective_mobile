CREATE TABLE users
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL,
    surname     VARCHAR(255) NOT NULL,
    patronymic  VARCHAR(255),
    age         INT,
    sex         VARCHAR(255),
    nationality VARCHAR(255)
);
