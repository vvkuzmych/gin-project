-- Create a table for storing todo_tasks
CREATE TABLE IF NOT EXISTS todo_tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    state BOOLEAN DEFAULT false
);
