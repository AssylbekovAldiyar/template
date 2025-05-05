CREATE TABLE IF NOT EXISTS Employee (
    Id SERIAL PRIMARY KEY,
    Name VARCHAR(255),
    Surname VARCHAR(255),
    Password VARCHAR(255),
    Is_Admin BOOLEAN,
    Activated BOOLEAN,
    Phone_Number VARCHAR(255),
    Enrolled TIMESTAMP
);
