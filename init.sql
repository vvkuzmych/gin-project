-- Create a table for storing todo tasks
CREATE TABLE IF NOT EXISTS todo_tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    state BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Insert initial data into the todo_tasks table
INSERT INTO todo_tasks (title, description, state) VALUES
    ('Task 1', 'Description of task 1', false),
    ('Task 2', 'Description of task 2', true),
    ('Task 3', 'Description of task 3', false);
