// Erasmo Cardoso da Silva
// Desenvolvedor Full Stack
//
// TODO: Adicionar mais tipos de erro específicos
// TODO: Implementar error wrapping
// TODO: Adicionar códigos de erro
// TODO: Melhorar mensagens de erro
// TODO: Adicionar stack trace
// TODO: Implementar error logging
// TODO: Adicionar error recovery
// TODO: Considerar usar error types do Go 1.13+
// TODO: Implementar error translation
// TODO: Adicionar error metrics

package service

import "errors"

var (
	// ErrNoSheets is returned when the Excel file has no sheets
	ErrNoSheets = errors.New("no sheets found in Excel file")

	// ErrEmptyFile is returned when the Excel file has no data
	ErrEmptyFile = errors.New("Excel file is empty")
)
