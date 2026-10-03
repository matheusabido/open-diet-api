ALTER TABLE users
    ADD COLUMN sex            VARCHAR(8)       CHECK (sex IN ('male', 'female')),
    ADD COLUMN age            SMALLINT         CHECK (age BETWEEN 1 AND 120),
    ADD COLUMN weight         DOUBLE PRECISION CHECK (weight > 0),
    ADD COLUMN height         SMALLINT         CHECK (height > 0),
    ADD COLUMN activity_level VARCHAR(16)      CHECK (activity_level IN ('sedentary', 'light', 'moderate', 'heavy', 'athlete')),
    ADD COLUMN objective      VARCHAR(16)      CHECK (objective IN ('cut', 'maintain', 'bulk'));
