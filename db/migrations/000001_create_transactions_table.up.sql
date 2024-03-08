CREATE TABLE IF NOT EXISTS transactions(
    key varchar PRIMARY KEY,
    amount      DECIMAL NOT NULL,
    account_key varchar NOT NULL,
    type        varchar NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    ); 