package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq" // Importa o driver do Postgres silenciosamente
)

func main() {
	// 1. Conectar ao PostgreSQL (Docker)
	// Usamos a porta 5432 pois você desativou o serviço nativo do Ubuntu
	connStr := "postgres://admin:adminpassword@localhost:5432/web3shield?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("🚨 Erro ao abrir conexão com o banco: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("🚨 Erro ao pingar o banco (o container está rodando?): %v", err)
	}
	fmt.Println("🐘 Conectado ao PostgreSQL com sucesso!")

	// 2. Criar a tabela automaticamente (se não existir)
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS raw_transactions (
		id SERIAL PRIMARY KEY,
		tx_hash VARCHAR(66) UNIQUE NOT NULL,
		to_address VARCHAR(42),
		function_signature VARCHAR(10),
		timestamp TIMESTAMPTZ DEFAULT NOW()
	);`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("🚨 Erro ao criar tabela: %v", err)
	}

	// 3. Conectar ao Anvil via WebSocket
	client, err := ethclient.Dial("ws://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("🚨 Erro ao conectar no Anvil: %v", err)
	}
	fmt.Println("🟢 Conectado com sucesso ao Anvil (Base Fork)!")

	// 4. Iniciar a escuta de novos blocos
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("🎧 Escutando mempool e salvando no Banco de Dados...")

	// 5. Loop infinito de processamento
	for {
		select {
		case err := <-sub.Err():
			log.Fatal("🚨 Erro na inscrição:", err)
		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				continue
			}

			txs := block.Transactions()
			if len(txs) > 0 {
				fmt.Printf("\n🚀 Bloco %s processado. Inserindo %d transações no BD...\n", header.Number.String(), len(txs))
			}

			for _, tx := range txs {
				// Descobre o endereço de destino
				toAddress := ""
				if tx.To() != nil {
					toAddress = tx.To().Hex()
				}

				// Extrai a Function Signature (primeiros 4 bytes do payload)
				txData := tx.Data()
				funcSig := "0x00000000" // Padrão para transferências simples (sem dados)
				if len(txData) >= 4 {
					funcSig = "0x" + hex.EncodeToString(txData[:4])
				}

				// Insere a transação no PostgreSQL
				insertQuery := `
				INSERT INTO raw_transactions (tx_hash, to_address, function_signature) 
				VALUES ($1, $2, $3) 
				ON CONFLICT (tx_hash) DO NOTHING;`
				
				_, err = db.Exec(insertQuery, tx.Hash().Hex(), toAddress, funcSig)
				if err != nil {
					log.Printf("⚠️ Erro ao salvar tx %s: %v\n", tx.Hash().Hex(), err)
				} else {
					fmt.Printf("  💾 Salvo: Hash %s | Sig: %s\n", tx.Hash().Hex()[:10]+"...", funcSig)
				}
			}
		}
	}
}