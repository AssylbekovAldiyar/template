CREATE TABLE IF NOT EXISTS users_permissions (
                                                 user_id int NOT NULL REFERENCES employee ON DELETE CASCADE,
                                                 permission_id int NOT NULL REFERENCES permissions ON DELETE CASCADE,
                                                 PRIMARY KEY (user_id, permission_id)
);