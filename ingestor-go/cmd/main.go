package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	
	"web3-shield/ingestor/internal/repository"
	"web3-shield/ingestor/internal/worker"
)

func main() {
	fmt.Println("🛡️ Iniciando Web3 Shield Ingestor...")

	// Carrega o .env (opcional: usa variáveis do ambiente do sistema quando não existir)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Nenhum arquivo .env encontrado. Usando variáveis do sistema.")
	}

	// 1. Inicializa Conexão com o Banco
	connStr := os.Getenv("DATABASE_URL")
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
	wsUrl := os.Getenv("WS_URL")
	client, err := ethclient.Dial(wsUrl)
	if err != nil {
		return fmt.Errorf("falha ao conectar no RPC: %w", err)
	}
	defer client.Close()

	// 4. Inicializa o Worker (Injetando o cliente Web3 e o Repositório)
	blockchainWorker := worker.NewBlockchainWorker(client, repo)

	// 5. Dá o Play!
	return blockchainWorker.Start()
}