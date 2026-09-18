# Web3 Shield — Cheat Sheet de Comandos (MVP)

## 1. Banco de Dados (Docker)

```bash
# Derruba o Postgres nativo do Ubuntu para liberar a porta 5432
sudo systemctl stop postgresql

# Sobe o container (o schema é criado automaticamente via infra/init.sql)
docker compose up -d
```

## 2. Nó Blockchain Local (Anvil)

```bash
# Fork da Base Mainnet (mantenha rodando em aba separada)
anvil --fork-url https://mainnet.base.org
```

## 3. Microsserviços (abas separadas)

```bash
# Ingestor Go (variáveis em ingestor-go/.env)
cd ingestor-go
go run cmd/main.go

# Analytics Python
cd analytics-py
source venv/bin/activate
python -m src.main
```

## 4. Testes

```bash
# Go (extração de assinatura)
cd ingestor-go/internal/domain
go test -v

# Python (heurística de risco)
cd analytics-py
source venv/bin/activate
python -m pytest src/domain/test_risk_scorer.py -v
```

## 5. Simulação de Ameaças (Novo terminal)

```bash
# A — Transferência nativa: risco esperado 10
cast send --value 0.1ether 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# B — Transferência ERC-20: risco esperado 30
cast send 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 "transfer(address,uint256)" 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC 100 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80

# C — Approve infinito (drain): risco esperado 90
cast send 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 "approve(address,uint256)" 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC 999999999999999 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

## 6. Auditoria de Dados (Validação Final)

```bash
# Transações cruas
docker exec -it web3_shield_db psql -U admin -d web3shield -c "SELECT id, left(tx_hash, 10) as hash, to_address, function_signature, block_number FROM raw_transactions;"

# Pipeline completo (Go -> Postgres -> Python -> Postgres)
docker exec -it web3_shield_db psql -U admin -d web3shield -c "
SELECT r.tx_hash, r.function_signature, a.risk_score, a.analyzed_at
FROM raw_transactions r
JOIN risk_analysis a ON r.tx_hash = a.tx_hash;
"
```