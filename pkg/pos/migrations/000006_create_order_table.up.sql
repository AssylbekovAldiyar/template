CREATE TABLE IF NOT EXISTS Orders (
    Id SERIAL PRIMARY KEY,
    Employee_ID VARCHAR(255),
    Total_Price FLOAT,
    Total_Paid FLOAT,
    Total_Return FLOAT,
    Receipt_ID VARCHAR(255),
    Created_At TIMESTAMP,
    Updated_At TIMESTAMP
);