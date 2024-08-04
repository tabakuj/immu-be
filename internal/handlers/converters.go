package handlers

import "immudb/internal/models"

func convertAccountInfoToDTO(sc *models.AccountInfo) *AccountInfoDto {
	if sc == nil {
		return nil
	}
	return &AccountInfoDto{
		AccountNumber: sc.Id,
		AccountName:   sc.Name,
		Iban:          sc.Iban,
		Address:       sc.Address,
		Amount:        sc.Amount,
		Type:          sc.Type,
	}
}

func convertDtoToAccountInfo(input *AccountInfoDto) *models.AccountInfo {
	if input == nil {
		return nil
	}
	return &models.AccountInfo{
		Id:      input.AccountNumber,
		Name:    input.AccountName,
		Iban:    input.Iban,
		Address: input.Address,
		Amount:  input.Amount,
		Type:    input.Type,
	}
}
