CREATE TYPE repeatable_task AS ENUM (
    'NEVER',
    'DAILY',
    'WEEKLY',
    'MONTHLY',
    'YEARLY'
);