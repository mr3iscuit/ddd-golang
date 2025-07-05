-- Drop trigger first
DROP TRIGGER IF EXISTS update_todos_updated_at ON todos;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_todos_created_at;
DROP INDEX IF EXISTS idx_todos_priority;
DROP INDEX IF EXISTS idx_todos_status;

-- Drop table
DROP TABLE IF EXISTS todos; 