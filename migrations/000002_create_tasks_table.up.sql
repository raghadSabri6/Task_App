CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4() UNIQUE,
    title TEXT NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_by_id UUID NOT NULL,
    assigned_to_id UUID,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,


    FOREIGN KEY (created_by_id) REFERENCES users(uuid),
    FOREIGN KEY (assigned_to_id) REFERENCES users(uuid)
);
