ALTER TABLE users
    DROP COLUMN IF EXISTS streak_last_updated,
    DROP COLUMN IF EXISTS streak_count;
