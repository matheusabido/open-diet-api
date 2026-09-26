CREATE TABLE users (
    id                BIGSERIAL PRIMARY KEY,
    username          VARCHAR(32)  NOT NULL UNIQUE,
    password          VARCHAR(255) NOT NULL,
    role              VARCHAR(16)  NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    kcal_goal         INTEGER      NOT NULL DEFAULT 0,
    protein_goal      INTEGER      NOT NULL DEFAULT 0,
    carbohydrate_goal INTEGER      NOT NULL DEFAULT 0,
    fat_goal          INTEGER      NOT NULL DEFAULT 0,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
