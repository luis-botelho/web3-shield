package domain

import "encoding/hex"

// 1. A Entidade: Define exatamente quais dados importam para o nosso negócio
type Transaction struct {
	Hash              string
	ToAddress         string
	FunctionSignature string
	BlockNumber       uint64
	BlockHash         string
}

// 2. O Contrato (Interface): Quem quiser salvar transações, tem que ter essa função "Save"
type TransactionRepository interface {
	Save(tx Transaction) error
}

// ExtractFunctionSignature pega o payload bruto da transação e retorna os 4 primeiros bytes em Hexadecimal.
func ExtractFunctionSignature(data []byte) string {
	if len(data) < 4 {
		return "0x00000000"
	}
	return "0x" + hex.EncodeToString(data[:4])
}