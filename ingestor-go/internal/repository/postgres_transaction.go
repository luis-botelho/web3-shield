package repository

import (
	"database/sql"
	"fmt"

	"web3-shield/ingestor/internal/domain"
)

// PostgresTransactionRepository persiste transações no PostgreSQL.
type PostgresTransactionRepository struct {
	db *sql.DB
}

// NewPostgresTransactionRepository cria um repositório associado à conexão fornecida.
func NewPostgresTransactionRepository(db *sql.DB) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{
		db: db,
	}
}

// Save insere a transação, ignorando hashes já persistidos.
func (r *PostgresTransactionRepository) Save(tx domain.Transaction) error {
	query := `
	INSERT INTO raw_transactions (tx_hash, from_address, to_address, function_signature, block_number, block_hash) 
	VALUES ($1, $2, $3, $4, $5, $6) 
	ON CONFLICT (tx_hash) DO NOTHING;`

	_, err := r.db.Exec(query, tx.Hash, tx.FromAddress, tx.ToAddress, tx.FunctionSignature, tx.BlockNumber, tx.BlockHash)
	if err != nil {
		return fmt.Errorf("erro ao salvar transação no banco: %w", err)
	}

	return nil
}
