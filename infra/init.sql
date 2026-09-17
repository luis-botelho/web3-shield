-- Web3 Shield — Esquema do Banco de Dados
-- Executado automaticamente na primeira subida do container Postgres
-- (mount em /docker-entrypoint-initdb.d/init.sql via docker-compose).

-- Transações cruas capturadas pelo Ingestor (Go) via WebSocket
CREATE TABLE IF NOT EXISTS raw_transactions (
    id                 BIGSERIAL PRIMARY KEY,
    tx_hash            TEXT UNIQUE NOT NULL,
    to_address         TEXT,
    function_signature TEXT,
    block_number       BIGINT,
    block_hash         TEXT,
    created_at         TIMESTAMP DEFAULT NOW()
);

-- Análises de risco geradas pelo Analytics Engine (Python)
CREATE TABLE IF NOT EXISTS risk_analysis (
    id          BIGSERIAL PRIMARY KEY,
    tx_hash     TEXT UNIQUE NOT NULL REFERENCES raw_transactions(tx_hash),
    risk_score  INTEGER NOT NULL,
    analyzed_at TIMESTAMP DEFAULT NOW()
);