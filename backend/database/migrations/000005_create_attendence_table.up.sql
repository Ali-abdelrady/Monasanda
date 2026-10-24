CREATE TABLE attendance(
    id SERIAL PRIMARY KEY NOT NULL,
    batch_id INT NOT NULL,
    name TEXT,
    location TEXT,
    phone_number VARCHAR(12),
    military_number VARCHAR(13),
    status attendance_status DEFAULT 'absent',
    notes TEXT,
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    checked_in_at TIMESTAMP,
    checked_in_by INT,
    created_by INT NOT NULL, -- Added data type (e.g., INT)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    -- Foreign Key for batch_id
    CONSTRAINT fk_attendance_batch 
        FOREIGN KEY (batch_id) 
        REFERENCES vacation_batches(id) 
        ON DELETE SET NULL,

    -- Foreign Key for checked_in_by (referencing users table)
    CONSTRAINT fk_attendance_checked_in_by 
        FOREIGN KEY (checked_in_by) 
        REFERENCES users(id) 
        ON DELETE RESTRICT,

    -- Foreign Key for created_by (referencing users table)
    CONSTRAINT fk_attendance_created_by 
        FOREIGN KEY (created_by) 
        REFERENCES users(id) 
        ON DELETE RESTRICT
);