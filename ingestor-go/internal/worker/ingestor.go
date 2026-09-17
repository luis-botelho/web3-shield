package worker

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"web3-shield/ingestor/internal/domain"
)

// BlockchainWorker gerencia a escuta da rede
type BlockchainWorker struct {
	client *ethclient.Client
	repo   domain.TransactionRepository // Injeção de Dependência! Ele só conhece o contrato.
}

// NewBlockchainWorker é o construtor. Ele RECEBE o cliente e o repositório prontos.
func NewBlockchainWorker(client *ethclient.Client, repo domain.TransactionRepository) *BlockchainWorker {
	return &BlockchainWorker{
		client: client,
		repo:   repo,
	}
}

// Start inicia o loop infinito de escuta e retorna erro se o WebSocket cair
func (w *BlockchainWorker) Start() error {
	headers := make(chan *types.Header)
	sub, err := w.client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		return fmt.Errorf("falha ao assinar novos blocos: %w", err)
	}
	defer sub.Unsubscribe()

	fmt.Println("🎧 Ouvinte Web3 ativo. Escutando mempool...")

	for {
		select {
		case err := <-sub.Err():
			return fmt.Errorf("erro no websocket: %w", err)

		case header := <-headers:
			block, err := w.client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Printf("⚠️ Erro ao buscar bloco %s: %v\n", header.Number.String(), err)
				continue
			}

			txs := block.Transactions()
			if len(txs) > 0 {
				fmt.Printf("\n🚀 Bloco %s recebido. Processando %d transações...\n", header.Number.String(), len(txs))
			}

			for _, tx := range txs {
				toAddress := ""
				if tx.To() != nil {
					toAddress = tx.To().Hex()
				}

				// 1. Usa a regra de negócio do DOMAIN para traduzir a assinatura
				funcSig := domain.ExtractFunctionSignature(tx.Data())

				// 2. Monta a Entidade
				transaction := domain.Transaction{
					Hash:              tx.Hash().Hex(),
					ToAddress:         toAddress,
					FunctionSignature: funcSig,
					BlockNumber:       block.NumberU64(),
					BlockHash:         block.Hash().Hex(),
				}

				// 3. Usa o REPOSITORY (sem saber que é Postgres) para salvar
				err = w.repo.Save(transaction)
				if err != nil {
					log.Printf("⚠️ Erro ao salvar tx %s: %v\n", transaction.Hash, err)
				} else {
					fmt.Printf("  💾 Salvo: Hash %s... | Sig: %s\n", transaction.Hash[:10], transaction.FunctionSignature)
				}
			}
		}
	}
}