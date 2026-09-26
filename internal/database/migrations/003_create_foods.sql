CREATE TABLE foods (
    id           BIGSERIAL PRIMARY KEY,
    name         VARCHAR(100)     NOT NULL,
    brand        VARCHAR(100),
    kcal         DOUBLE PRECISION NOT NULL CHECK (kcal >= 0),
    protein      DOUBLE PRECISION NOT NULL CHECK (protein >= 0),
    carbohydrate DOUBLE PRECISION NOT NULL CHECK (carbohydrate >= 0),
    fat          DOUBLE PRECISION NOT NULL CHECK (fat >= 0),
    amount       DOUBLE PRECISION NOT NULL CHECK (amount > 0),
    unit         VARCHAR(4)       NOT NULL CHECK (unit IN ('g', 'kg', 'mg', 'ml', 'L', 'dL', 'oz', 'lb')),
    created_by   BIGINT           REFERENCES users (id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);
