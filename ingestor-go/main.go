package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
)

func main() {
	// 1. Conexão com o Banco de Dados (PostgreSQL)
	connStr := "postgres://admin:adminpassword@localhost:5432/web3shield?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("🚨 Erro ao abrir conexão com o banco: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("🚨 Erro ao pingar o banco: %v", err)
	}
	fmt.Println("🐘 Conectado ao PostgreSQL com sucesso!")

	// 2. Atualizamos a tabela para incluir block_number e block_hash (Prevenção de Reorgs)
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS raw_transactions (
		id SERIAL PRIMARY KEY,
		tx_hash VARCHAR(66) UNIQUE NOT NULL,
		to_address VARCHAR(42),
		function_signature VARCHAR(10),
		block_number BIGINT,
		block_hash VARCHAR(66),
		timestamp TIMESTAMPTZ DEFAULT NOW()
	);`
	if _, err = db.Exec(createTableQuery); err != nil {
		log.Fatalf("🚨 Erro ao criar tabela: %v", err)
	}

	// 3. Loop principal de Resiliência (Reconexão Automática)
	for {
		err := runIngestor(db)
		if err != nil {
			log.Printf("⚠️ Conexão perdida ou erro crítico: %v", err)
			log.Println("🔄 Tentando reconectar em 5 segundos...")
			time.Sleep(5 * time.Second) // Aguarda antes de tentar de novo
		}
	}
}

// runIngestor concentra a lógica de escuta. Se ela retornar um erro, o loop do main() a reinicia.
func runIngestor(db *sql.DB) error {
	client, err := ethclient.Dial("ws://127.0.0.1:8545")
	if err != nil {
		return fmt.Errorf("falha ao conectar no RPC: %w", err)
	}
	defer client.Close()
	fmt.Println("🟢 Conectado com sucesso ao Anvil via WebSocket!")

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return fmt.Errorf("falha ao assinar novos blocos: %w", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("🎧 Escutando mempool de forma resiliente...")

	for {
		select {
		case err := <-sub.Err():
			// Em vez de log.Fatal, nós retornamos o erro para o main() acionar a reconexão
			return fmt.Errorf("erro no websocket: %w", err)

		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Printf("⚠️ Erro ao buscar bloco %s: %v\n", header.Number.String(), err)
				continue
			}

			txs := block.Transactions()
			if len(txs) > 0 {
				fmt.Printf("\n🚀 Bloco %s processado. Inserindo %d transações...\n", header.Number.String(), len(txs))
			}

			for _, tx := range txs {
				toAddress := ""
				if tx.To() != nil {
					toAddress = tx.To().Hex()
				}

				txData := tx.Data()
				funcSig := "0x00000000"
				if len(txData) >= 4 {
					funcSig = "0x" + hex.EncodeToString(txData[:4])
				}

				// Agora salvamos também o contexto do bloco para lidar com Reorgs no futuro
				insertQuery := `
				INSERT INTO raw_transactions (tx_hash, to_address, function_signature, block_number, block_hash) 
				VALUES ($1, $2, $3, $4, $5) 
				ON CONFLICT (tx_hash) DO NOTHING;`

				_, err = db.Exec(insertQuery, tx.Hash().Hex(), toAddress, funcSig, block.NumberU64(), block.Hash().Hex())
				if err != nil {
					log.Printf("⚠️ Erro ao salvar tx %s: %v\n", tx.Hash().Hex(), err)
				} else {
					fmt.Printf("  💾 Salvo: Hash %s... | Bloco: %d\n", tx.Hash().Hex()[:10], block.NumberU64())
				}
			}
		}
	}
}