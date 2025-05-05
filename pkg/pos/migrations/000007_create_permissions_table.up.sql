CREATE TABLE IF NOT EXISTS permissions (
                                           id serial PRIMARY KEY,
                                           code text NOT NULL
);

INSERT INTO permissions (code)
VALUES ('menus:read'),
			 ('menus:write');