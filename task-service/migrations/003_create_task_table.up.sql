CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    title VARCHAR(255),
    description TEXT,
    status task_status NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    repeatable repeatable_task DEFAULT 'NEVER'
);

CREATE INDEX IF NOT EXISTS idx_id ON tasks(id);