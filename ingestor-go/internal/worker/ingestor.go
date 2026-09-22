package worker

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"web3-shield/ingestor/internal/domain"
)

// BlockchainWorker processa blocos recebidos pela assinatura WebSocket.
type BlockchainWorker struct {
	client *ethclient.Client
	repo   domain.TransactionRepository
}

// NewBlockchainWorker cria um worker com as dependências de rede e persistência.
func NewBlockchainWorker(client *ethclient.Client, repo domain.TransactionRepository) *BlockchainWorker {
	return &BlockchainWorker{
		client: client,
		repo:   repo,
	}
}

// Start consome blocos até a assinatura falhar ou ser cancelada.
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

				// O remetente é derivado da assinatura ECDSA da transação; o Signer
				// conhece as regras da rede (chain ID) usadas para recuperá-lo.
				fromAddress := "0x0000000000000000000000000000000000000000"
				if sender, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx); err == nil {
					fromAddress = strings.ToLower(sender.Hex())
				}

				// O seletor ABI identifica a função chamada pelo contrato.
				funcSig := domain.ExtractFunctionSignature(tx.Data())

				transaction := domain.Transaction{
					Hash:              tx.Hash().Hex(),
					FromAddress:       fromAddress,
					ToAddress:         toAddress,
					FunctionSignature: funcSig,
					BlockNumber:       block.NumberU64(),
					BlockHash:         block.Hash().Hex(),
				}

				err = w.repo.Save(transaction)
				if err != nil {
					log.Printf("⚠️ Erro ao salvar tx %s: %v\n", transaction.Hash, err)
				} else {
					fmt.Printf("  💾 Salvo: Hash %s... | From: %s | Sig: %s\n", transaction.Hash[:10], transaction.FromAddress, transaction.FunctionSignature)
				}
			}
		}
	}
}
