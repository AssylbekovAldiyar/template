# POS Cash Register API

This is a RESTful API built with Golang for managing employees and products in a Point of Sale (POS) cash register system.
The POS Cash Register API, developed using Golang, facilitates the management of employees and products within a Point of Sale system. Below are detailed descriptions of the endpoints and database schemas.

## Endpoints

### Employees

- **GET /employees**: Retrieve all employees.
- **GET /employees/{id}**: Retrieve an employee by ID.
- **POST /employees**: Register a new employee.
- **PUT /employees/{id}**: Update an existing employee.
- **DELETE /employees/{id}**: Delete an employee.

### Products

- **GET /products**: Retrieve all products.
- **GET /products/{productId}**: Retrieve a product by ID.
- **POST /products**: Create a new product.
- **PUT /products/{productId}**: Update an existing product.
- **DELETE /products/{productId}**: Delete a product.

### Employee Table

```sql
employee (
    id VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    surname VARCHAR(255),
    password BYTEA,
    is_admin BOOLEAN,
    phone_number VARCHAR(20),
    enrolled TIMESTAMP
)

categories (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)

product (
    id VARCHAR(255) PRIMARY KEY,
    title VARCHAR(255),
    category_id VARCHAR(255),
    price INT,
    description TEXT,
    amount INT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)

orders (
   id VARCHAR(255) PRIMARY KEY,
    user_id VARCHAR(255) REFERENCES employee(id),
    payment_type_id VARCHAR(255) REFERENCES payments(id),
    total_price INTEGER,
    total_paid INTEGER,
    total_return INTEGER,
    receipt_id VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);



order_product (
    id VARCHAR(255),
    order_id VARCHAR(255) REFERENCES orders(id),
    product_id VARCHAR(255) REFERENCES product(id),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

payments (
    id VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

### Team
Zholdybay Altynbek - 22B030543
