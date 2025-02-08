CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

\i /docker-entrypoint-initdb.d/migrations/001_create_tables.sql
\i /docker-entrypoint-initdb.d/migrations/002_create_indexes.sql


-- We will not want to seed test data outside of dev
\if :'ENVIRONMENT' = 'development'
    \echo 'Running development seeds...'
    \i /docker-entrypoint-initdb.d/seeds/001_seed_test_data.sql
\endif
\i /docker-entrypoint-initdb.d/seeds/002_seed_funds.sql

