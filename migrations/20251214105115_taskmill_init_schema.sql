-- +goose Up
-- +goose StatementBegin

-- Create schema
CREATE SCHEMA taskmill;

-- Main task queue table
CREATE TABLE taskmill.task_queue (
    id BIGSERIAL PRIMARY KEY,

    -- Routing
    queue_name VARCHAR(255) NOT NULL,
    task_group_id VARCHAR(255),
    operation_id VARCHAR(255) NOT NULL,

    -- Content
    meta JSONB,
    payload JSONB NOT NULL,

    -- Timing Control
    scheduled_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    visible_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMPTZ,

    -- Processing
    priority INT NOT NULL DEFAULT 0,
    attempts INT NOT NULL DEFAULT 0,
    max_attempts INT NOT NULL DEFAULT 3,

    -- Idempotency
    idempotency_key VARCHAR(255) NOT NULL,

    -- Audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- DLQ
    dlq_at TIMESTAMPTZ,
    dlq_reason JSONB,

    -- Ephemeral flag
    ephemeral BOOLEAN NOT NULL DEFAULT FALSE
);

-- Task results table for completed tasks
CREATE TABLE taskmill.task_results (
    id BIGINT PRIMARY KEY,

    -- Routing
    queue_name VARCHAR(255) NOT NULL,
    task_group_id VARCHAR(255),
    operation_id VARCHAR(255) NOT NULL,

    -- Content
    meta JSONB,
    payload JSONB NOT NULL,

    -- Processing
    priority INT NOT NULL,
    attempts INT NOT NULL,
    max_attempts INT NOT NULL,

    -- Idempotency
    idempotency_key VARCHAR(255) NOT NULL,

    -- Timing
    scheduled_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Task schedules table for cron-based scheduling
CREATE TABLE taskmill.task_schedules (
    id BIGSERIAL PRIMARY KEY,

    -- Identity
    operation_id VARCHAR(255) NOT NULL UNIQUE,
    queue_name VARCHAR(255) NOT NULL,
    cron_pattern VARCHAR(100) NOT NULL,

    -- State
    next_run_at TIMESTAMPTZ NOT NULL,

    -- Execution tracking
    last_run_at TIMESTAMPTZ,
    last_run_status VARCHAR(20),
    last_error TEXT,
    run_count BIGINT NOT NULL DEFAULT 0,

    -- Audit
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Critical indexes for performance

-- Dequeue index (hot path) - optimized for the main dequeue query
-- Partial index only on active tasks (not in DLQ)
CREATE INDEX idx_task_queue_dequeue
    ON taskmill.task_queue (queue_name, priority DESC, visible_at, scheduled_at, id ASC)
    WHERE dlq_at IS NULL;

-- Task group index for FIFO ordering
CREATE INDEX idx_task_queue_group
    ON taskmill.task_queue (task_group_id, priority DESC, id ASC)
    WHERE task_group_id IS NOT NULL AND dlq_at IS NULL;

-- Idempotency index for duplicate detection (unique constraint)
CREATE UNIQUE INDEX idx_task_queue_idempotency
    ON taskmill.task_queue (queue_name, idempotency_key)
    WHERE dlq_at IS NULL;

-- Scheduled tasks index
CREATE INDEX idx_task_queue_scheduled
    ON taskmill.task_queue (scheduled_at)
    WHERE dlq_at IS NULL;

-- DLQ index for querying dead letter queue
CREATE INDEX idx_task_queue_dlq
    ON taskmill.task_queue (queue_name, dlq_at DESC)
    WHERE dlq_at IS NOT NULL;

-- Operation ID index for filtering by task type
CREATE INDEX idx_task_queue_operation
    ON taskmill.task_queue (queue_name, operation_id)
    WHERE dlq_at IS NULL;

-- Indexes for task_results table

-- Index for querying results by queue and completion time
CREATE INDEX idx_task_results_queue_completed
    ON taskmill.task_results (queue_name, completed_at DESC);

-- Index for cleanup operations
CREATE INDEX idx_task_results_completed
    ON taskmill.task_results (completed_at);

-- Index for querying by idempotency key
CREATE INDEX idx_task_results_idempotency
    ON taskmill.task_results (queue_name, idempotency_key);

-- Index for querying by operation ID
CREATE INDEX idx_task_results_operation
    ON taskmill.task_results (queue_name, operation_id, completed_at DESC);

-- Index for finding due schedules
CREATE INDEX idx_task_schedules_due
    ON taskmill.task_schedules (next_run_at);

-- Index for querying by queue
CREATE INDEX idx_task_schedules_queue
    ON taskmill.task_schedules (queue_name);

-- Update timestamp trigger function
CREATE FUNCTION taskmill.update_task_queue_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Update timestamp trigger
CREATE TRIGGER trigger_task_queue_updated_at
    BEFORE UPDATE ON taskmill.task_queue
    FOR EACH ROW
    EXECUTE FUNCTION taskmill.update_task_queue_updated_at();

-- Stats view for monitoring
CREATE VIEW taskmill.task_queue_stats AS
SELECT
    queue_name,
    COUNT(*) FILTER (WHERE dlq_at IS NULL) as total,
    COUNT(*) FILTER (WHERE visible_at <= CURRENT_TIMESTAMP AND scheduled_at <= CURRENT_TIMESTAMP AND dlq_at IS NULL) as available,
    COUNT(*) FILTER (WHERE visible_at > CURRENT_TIMESTAMP AND dlq_at IS NULL) as in_flight,
    COUNT(*) FILTER (WHERE scheduled_at > CURRENT_TIMESTAMP AND dlq_at IS NULL) as scheduled,
    COUNT(*) FILTER (WHERE dlq_at IS NOT NULL) as in_dlq,
    MIN(created_at) FILTER (WHERE dlq_at IS NULL) as oldest_task,
    AVG(attempts) FILTER (WHERE dlq_at IS NULL) as avg_attempts,
    PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY attempts) FILTER (WHERE dlq_at IS NULL) as p95_attempts
FROM taskmill.task_queue
GROUP BY queue_name;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop view first (depends on table)
DROP VIEW IF EXISTS taskmill.task_queue_stats;

-- Drop trigger (depends on function and table)
DROP TRIGGER IF EXISTS trigger_task_queue_updated_at ON taskmill.task_queue;

-- Drop function
DROP FUNCTION IF EXISTS taskmill.update_task_queue_updated_at();

-- Drop indexes for task_schedules
DROP INDEX IF EXISTS taskmill.idx_task_schedules_queue;
DROP INDEX IF EXISTS taskmill.idx_task_schedules_due;

-- Drop indexes for task_results
DROP INDEX IF EXISTS taskmill.idx_task_results_operation;
DROP INDEX IF EXISTS taskmill.idx_task_results_idempotency;
DROP INDEX IF EXISTS taskmill.idx_task_results_completed;
DROP INDEX IF EXISTS taskmill.idx_task_results_queue_completed;

-- Drop indexes for task_queue
DROP INDEX IF EXISTS taskmill.idx_task_queue_operation;
DROP INDEX IF EXISTS taskmill.idx_task_queue_dlq;
DROP INDEX IF EXISTS taskmill.idx_task_queue_scheduled;
DROP INDEX IF EXISTS taskmill.idx_task_queue_idempotency;
DROP INDEX IF EXISTS taskmill.idx_task_queue_group;
DROP INDEX IF EXISTS taskmill.idx_task_queue_dequeue;

-- Drop tables (reverse order of creation)
DROP TABLE IF EXISTS taskmill.task_schedules;
DROP TABLE IF EXISTS taskmill.task_results;
DROP TABLE IF EXISTS taskmill.task_queue;

-- Drop schema
DROP SCHEMA IF EXISTS taskmill;

-- +goose StatementEnd
