# Web3 Shield 🛡️

O **Web3 Shield** é um motor open-source de prevenção de fraudes e análise de risco para redes EVM, com foco nativo na **Base Mainnet**. Projetado sob os princípios de *Clean Architecture* e microsserviços, o sistema captura transações em tempo real, analisa assinaturas de contratos e alerta sobre riscos críticos (como ataques de *phishing* e drenagem de carteiras) antes que o estrago seja irreversível.

---

## 🏗️ Arquitetura do Sistema

O projeto é um monorepo que orquestra serviços assíncronos e altamente desacoplados:

1. **Ingestor (Go):** O "Cão de Guarda". Conecta-se à blockchain via WebSocket (`SubscribeNewHead`), aplica tratamento defensivo contra reorgs e persiste dados brutos. Construído com resiliência para reconexão automática e estruturado em *Domain*, *Repository* e *Worker*.
2. **Analytics Engine (Python):** O "Detetive". Roda em loop isolado consumindo as transações pendentes. Aplica heurísticas de segurança (ex: identificação de `approve` maliciosos) e gera um *Risk Score* (0-100).
3. **Database (PostgreSQL):** Ponto de sincronia central, operando com esquemas normalizados (`raw_transactions`, `risk_analysis`) e inserções idempotentes.
4. **REST API (Node.js/Fastify + Prisma):** A "Recepção". Porta de entrada para integrações B2B e frontends. Expõe endpoints de *healthcheck*, consulta das análises mais recentes (`/risks`) e risco consolidado por carteira (`/wallet/:address/risk`).

---

## 🚀 Quickstart (Ambiente de Desenvolvimento)

### Pré-requisitos

* [Docker & Docker Compose](https://docs.docker.com/compose/)
* [Go 1.24+](https://go.dev/)
* [Python 3.10+](https://www.python.org/)
* [Node.js 18+](https://nodejs.org/)
* [Foundry (Anvil & Cast)](https://book.getfoundry.sh/)

### 1. Clonar e Configurar

```bash
git clone https://github.com/luis-botelho/web3-shield.git
cd web3-shield

# Subir a infraestrutura (PostgreSQL com init.sql automático)
docker compose up -d
```

Configure as variáveis de ambiente (`.env`) em `ingestor-go`, `analytics-py` e `api-node` utilizando os arquivos `.env.example` como base.

### 2. Iniciar a Blockchain Local (Fork da Base)

Em um terminal dedicado, inicie o nó local espelhando a Mainnet:

```bash
anvil --fork-url https://mainnet.base.org
```

### 3. Rodar os Microsserviços

Abra três terminais independentes:

**Terminal A (Ingestor - Go)**

```bash
cd ingestor-go
go run cmd/main.go
```

**Terminal B (Analytics - Python)**

```bash
cd analytics-py
source venv/bin/activate
python -m src.main
```

**Terminal C (REST API - Node.js)**

```bash
cd api-node
npm install          # instala dependências
npx prisma generate  # gera o client a partir do schema
npm run dev          # sobe o servidor em http://localhost:3000
```

### 4. Simular Ameaças

Em um novo terminal, utilize o `cast` para forjar transações e observar a esteira de análise em tempo real:

**Transferência Comum (Risco 30):**

```bash
cast send 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 "transfer(address,uint256)" 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC 100 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

**Tentativa de Dreno / Phishing (Risco 90):**

```bash
cast send 0x70997970C51812dc3A010C7d01b50e0d17dc79C8 "approve(address,uint256)" 0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC 999999999999999 --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
```

---

## 🧪 Qualidade de Código e Testes

O projeto segue a cultura de **Test-Driven Development (TDD)** em suas regras de negócio centrais, garantindo que a lógica não dependa de infraestrutura externa.

* **Testes Go (Domínio):** `cd ingestor-go/internal/domain && go test -v`
* **Testes Python (Heurística):** `cd analytics-py && python -m pytest src/domain/test_risk_scorer.py -v`

---

## 🛣️ Roadmap e Escala (Produção)

O Web3 Shield foi arquitetado para evoluir de um projeto de pesquisa para uma infraestrutura crítica escalável, buscando sustentabilidade através de programas como o **Base Builder Grant** e **Gitcoin Grants**.

- [x] **Fase 1 (MVP):** Pipeline de ingestão, persistência resiliente e motor de pontuação isolado (Ambiente Local/Anvil).
- [x] **Fase 2 (Integração):** API REST em Node.js (Fastify + Prisma) para consumo de dados — *healthcheck*, análises recentes e risco por carteira.
- [ ] **Fase 3 (Sniper Mode):** Transição do modelo *Firehose* (escutar toda a rede) para um filtro em memória no Go, escutando exclusivamente alvos cadastrados (redução de 90%+ em custo de RPC).
- [ ] **Fase 4 (Alertas):** Integração com webhooks, Discord e Telegram para notificações de alto risco em tempo real (< 10s).
- [ ] **Fase 5 (Mainnet):** Deploy via Alchemy/Supabase para a Base Mainnet e abertura para Beta Testers.