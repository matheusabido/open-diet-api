CREATE TABLE food_items (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT           NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    food_id    BIGINT           NOT NULL REFERENCES foods (id) ON DELETE CASCADE,
    date       DATE             NOT NULL,
    meal       VARCHAR(16)      NOT NULL CHECK (meal IN ('breakfast', 'lunch', 'dinner', 'snacks')),
    amount     DOUBLE PRECISION NOT NULL CHECK (amount > 0),
    unit       VARCHAR(4)       NOT NULL CHECK (unit IN ('g', 'kg', 'mg', 'ml', 'L', 'dL', 'oz', 'lb')),
    created_at TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX food_items_user_id_date_idx ON food_items (user_id, date);
CREATE INDEX food_items_food_id_idx ON food_items (food_id);
