package handlers

import (
	"github.com/stretchr/testify/assert"
	"immudb/internal/models"
	"immudb/internal/persistance"
	"testing"
)

func TestHandler_ModelToDtoConverters(t *testing.T) {
	tests := []struct {
		name     string
		input    *models.AccountInfo
		expected *AccountInfoDto
	}{
		{
			name: "validity_all_fields",
			input: &models.AccountInfo{
				Id:      1,
				Name:    "tt",
				Iban:    "iban",
				Amount:  11.2,
				Address: persistance.AddPointer("addr"),
				Type:    persistance.AddPointer(models.Receiving),
			},
			expected: &AccountInfoDto{
				AccountNumber: 1,
				AccountName:   "tt",
				Iban:          "iban",
				Amount:        11.2,
				Address:       persistance.AddPointer("addr"),
				Type:          persistance.AddPointer(models.Receiving),
			},
		},
		{
			name: "validity_all_required_fields",
			input: &models.AccountInfo{
				Id:     1,
				Name:   "tt",
				Iban:   "iban",
				Amount: 11.2,
			},
			expected: &AccountInfoDto{
				AccountNumber: 1,
				AccountName:   "tt",
				Iban:          "iban",
				Amount:        11.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// action
			result := convertAccountInfoToDTO(tt.input)

			//assert
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHandler_DtoToModelConverters(t *testing.T) {
	tests := []struct {
		name     string
		input    *AccountInfoDto
		expected *models.AccountInfo
	}{
		{
			name: "validity_all_fields",
			input: &AccountInfoDto{
				AccountNumber: 1,
				AccountName:   "tt",
				Iban:          "iban",
				Amount:        11.2,
				Address:       persistance.AddPointer("addr"),
				Type:          persistance.AddPointer(models.Receiving),
			},
			expected: &models.AccountInfo{
				Id:      1,
				Name:    "tt",
				Iban:    "iban",
				Amount:  11.2,
				Address: persistance.AddPointer("addr"),
				Type:    persistance.AddPointer(models.Receiving),
			},
		},
		{
			name: "validity_all_required_fields",
			input: &AccountInfoDto{
				AccountNumber: 1,
				AccountName:   "tt",
				Iban:          "iban",
				Amount:        11.2,
			},
			expected: &models.AccountInfo{
				Id:     1,
				Name:   "tt",
				Iban:   "iban",
				Amount: 11.2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// action
			result := convertDtoToAccountInfo(tt.input)

			//assert
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}
