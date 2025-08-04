-- Truncate tables with cascade for FK relations
TRUNCATE audit_logs, print_jobs, labels, users RESTART IDENTITY CASCADE;
