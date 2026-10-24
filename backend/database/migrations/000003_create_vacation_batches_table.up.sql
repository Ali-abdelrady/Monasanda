CREATE TABLE vacation_batches (
    id SERIAL PRIMARY KEY NOT NULL,
    name TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    created_by int NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT fk_vacation_batches_user 
        FOREIGN KEY (created_by) 
        REFERENCES users(id) 
        ON DELETE CASCADE
)