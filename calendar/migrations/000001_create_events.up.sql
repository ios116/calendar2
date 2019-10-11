BEGIN;
CREATE TABLE IF NOT EXISTS events
(
    id          serial PRIMARY KEY,
    title       VARCHAR(300) NOT NULL,
    date        timestamp(8) WITHOUT TIME ZONE NOT NULL,
    duration    interval(16) NOT NULL,
    author      VARCHAR(300) NOT NULL,
    description TEXT,
    notify      interval(16)
);
COMMIT;



