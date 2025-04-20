CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status task_status NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    repeatable repeatable_task DEFAULT 'NEVER',
    parent_task_id UUID REFERENCES tasks(id) ON DELETE SET NULL    
);

CREATE INDEX IF NOT EXISTS idx_user_id ON tasks(user_id);