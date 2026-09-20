package domain

import (
	"testing"
)

func TestExtractFunctionSignature(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
	}{
		{
			name:     "Transferencia de ETH simples (Sem payload)",
			input:    []byte{},
			expected: "0x00000000",
		},
		{
			name:     "Chamada de Contrato Valida (Ex: Approve ou Transfer)",
			input:    []byte{0xa9, 0x05, 0x9c, 0xbb, 0x01, 0x02, 0x03},
			expected: "0xa9059cbb",
		},
		{
			name:     "Payload Invalido (Menos de 4 bytes)",
			input:    []byte{0xaa, 0xbb},
			expected: "0x00000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractFunctionSignature(tt.input)
			if result != tt.expected {
				t.Errorf("Erro no cenário '%s': esperava %s, mas recebeu %s", tt.name, tt.expected, result)
			}
		})
	}
}
