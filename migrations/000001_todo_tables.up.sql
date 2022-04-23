CREATE TABLE todos (
    todo_id UUID PRIMARY KEY,
    title VARCHAR(128) NOT NULL,
    note TEXT,
    priority VARCHAR(16),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);

INSERT INTO 
    todos(todo_id, title, note, priority) 
VALUES ('29d40512-235f-47dd-8512-fd9cfcc817fc', 'yugurish', 'ertala ...','medium');