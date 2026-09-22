-- Schema inicial do PostgreSQL para o ambiente local.

-- Eventos confirmados capturados pelo ingestor.
CREATE TABLE IF NOT EXISTS raw_transactions (
    id                 BIGSERIAL PRIMARY KEY,
    tx_hash            TEXT UNIQUE NOT NULL,
    from_address       TEXT,
    to_address         TEXT,
    function_signature TEXT,
    block_number       BIGINT,
    block_hash         TEXT,
    created_at         TIMESTAMP DEFAULT NOW()
);

-- Classificações de risco produzidas pelo serviço de analytics.
CREATE TABLE IF NOT EXISTS risk_analysis (
    id          BIGSERIAL PRIMARY KEY,
    tx_hash     TEXT UNIQUE NOT NULL REFERENCES raw_transactions(tx_hash),
    risk_score  INTEGER NOT NULL,
    analyzed_at TIMESTAMP DEFAULT NOW()
);
