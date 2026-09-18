package domain

import (
	"testing"
)

func TestExtractFunctionSignature(t *testing.T) {
	// A nossa "tabela" de cenários de teste
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Transferencia de ETH simples (Sem payload)",
			input:    []byte{}, // Array de bytes vazio
			expected: "0x00000000",
		},
		{
			name:     "Chamada de Contrato Valida (Ex: Approve ou Transfer)",
			input:    []byte{0xa9, 0x05, 0x9c, 0xbb, 0x01, 0x02, 0x03}, // Payload com mais de 4 bytes
			expected: "0xa9059cbb", // Deve pegar apenas os 4 primeiros
		},
		{
			name:     "Payload Invalido (Menos de 4 bytes)",
			input:    []byte{0xaa, 0xbb},
			expected: "0x00000000",
		},
	}

	// O loop que roda cada cenário
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractFunctionSignature(tt.input)
			if result != tt.expected {
				t.Errorf("Erro no cenário '%s': esperava %s, mas recebeu %s", tt.name, tt.expected, result)
			}
		})
	}
}