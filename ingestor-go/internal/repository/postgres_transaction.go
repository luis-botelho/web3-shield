package repository

import (
	"database/sql"
	"fmt"

	// Importa o nosso domínio para conhecer a entidade Transaction
	"web3-shield/ingestor/internal/domain"
)

// PostgresTransactionRepository é o nosso "Arquivista" focado no PostgreSQL
type PostgresTransactionRepository struct {
	db *sql.DB
}

// NewPostgresTransactionRepository é o "construtor" que recebe a conexão do banco
func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{
		db: db,
	}
}

// Save é a função que cumpre o contrato exigido pelo domain.TransactionRepository
func (r *PostgresTransactionRepository) Save(tx domain.Transaction) error {
	query := `
	INSERT INTO raw_transactions (tx_hash, to_address, function_signature, block_number, block_hash) 
	VALUES ($1, $2, $3, $4, $5) 
	ON CONFLICT (tx_hash) DO NOTHING;`

	// O repositório SÓ sabe executar SQL. Ele não sabe o que é Web3, Anvil ou Blockchain.
	_, err := r.db.Exec(query, tx.Hash, tx.ToAddress, tx.FunctionSignature, tx.BlockNumber, tx.BlockHash)
	if err != nil {
		return fmt.Errorf("erro ao salvar transação no banco: %w", err)
	}

	return nil
}
