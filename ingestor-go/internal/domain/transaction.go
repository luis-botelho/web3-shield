package domain

import "encoding/hex"

// Transaction representa uma transação confirmada e seu contexto de bloco.
type Transaction struct {
	Hash              string
	FromAddress       string
	ToAddress         string
	FunctionSignature string
	BlockNumber       uint64
	BlockHash         string
}

// TransactionRepository define a persistência de transações.
type TransactionRepository interface {
	Save(tx Transaction) error
}

// ExtractFunctionSignature retorna o seletor ABI de quatro bytes do payload.
func ExtractFunctionSignature(data []byte) string {
	if len(data) < 4 {
		return "0x00000000"
	}
	return "0x" + hex.EncodeToString(data[:4])
}
