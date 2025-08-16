-- Improve performance for heat number queries over recent data
-- Creates a composite index on (heat_no, created_at DESC)
-- Safe to run multiple times due to IF NOT EXISTS

BEGIN;

-- Create composite index to support queries like:
-- SELECT ... FROM print_jobs WHERE heat_no = $1 AND created_at >= NOW() - INTERVAL '10 days'
-- The ordering with created_at DESC helps if you later sort by recency.
CREATE INDEX IF NOT EXISTS idx_print_jobs_heat_no_created_at
    ON print_jobs (heat_no, created_at DESC);

COMMIT;
