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

	// Carrega configurações locais sem sobrescrever variáveis já definidas no ambiente.
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Nenhum arquivo .env encontrado. Usando variáveis do sistema.")
	}

	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("🚨 Erro fatal no banco: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("🚨 Postgres indisponível: %v", err)
	}

	repo := repository.NewPostgresTransactionRepository(db)

	// Reinicia a assinatura após falhas de conexão com o nó.
	for {
		err := run(repo)
		if err != nil {
			log.Printf("⚠️ Conexão perdida: %v", err)
			log.Println("🔄 Tentando reconectar em 5 segundos...")
			time.Sleep(5 * time.Second)
		}
	}
}

// run cria os recursos associados a uma sessão de escuta do nó.
func run(repo *repository.PostgresTransactionRepository) error {
	wsUrl := os.Getenv("WS_URL")
	client, err := ethclient.Dial(wsUrl)
	if err != nil {
		return fmt.Errorf("falha ao conectar no RPC: %w", err)
	}
	defer client.Close()

	blockchainWorker := worker.NewBlockchainWorker(client, repo)

	return blockchainWorker.Start()
}
