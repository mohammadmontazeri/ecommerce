CREATE TABLE categories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) UNIQUE NOT NULL,
		parent_id INT,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		CONSTRAINT fk_parent_id
	  	FOREIGN KEY(parent_id)
	  	REFERENCES categories(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
        ) ;

        