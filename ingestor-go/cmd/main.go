package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/lib/pq"
	
	"web3-shield/ingestor/internal/repository"
	"web3-shield/ingestor/internal/worker"
)

func main() {
	fmt.Println("🛡️ Iniciando Web3 Shield Ingestor...")

	// 1. Inicializa Conexão com o Banco
	connStr := "postgres://admin:adminpassword@localhost:5432/web3shield?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("🚨 Erro fatal no banco: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("🚨 Postgres indisponível: %v", err)
	}
	
	// 2. Inicializa o Repositório (Injetando o banco)
	repo := repository.NewPostgresTransactionRepository(db)

	// 3. Loop de Resiliência Principal
	for {
		err := run(repo)
		if err != nil {
			log.Printf("⚠️ Conexão perdida: %v", err)
			log.Println("🔄 Tentando reconectar em 5 segundos...")
			time.Sleep(5 * time.Second)
		}
	}
}

// run centraliza a inicialização do cliente Web3 para facilitar o retry
func run(repo *repository.PostgresTransactionRepository) error {
	// Inicializa Conexão com a Blockchain
	client, err := ethclient.Dial("ws://127.0.0.1:8545")
	if err != nil {
		return fmt.Errorf("falha ao conectar no RPC: %w", err)
	}
	defer client.Close()

	// 4. Inicializa o Worker (Injetando o cliente Web3 e o Repositório)
	blockchainWorker := worker.NewBlockchainWorker(client, repo)

	// 5. Dá o Play!
	return blockchainWorker.Start()
}