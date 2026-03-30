package storage

import (
	"encoding/csv"
	"fmt"
	"gosql-br/internal/ast"
	"os"
	"strconv"
	"strings"
	"time"
)

type CSVStorage struct {
	FileName string
	Data     [][]string
	Header   map[string]int
}

func NewCSVStorage(fileName string) (*CSVStorage, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("não foi possível abrir o arquivo: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	headerMap := make(map[string]int)
	for i, name := range records[0] {
		headerMap[strings.ToLower(strings.TrimSpace(name))] = i
	}

	return &CSVStorage{
		FileName: fileName,
		Data:     records[1:],
		Header:   headerMap,
	}, nil
}

func (s *CSVStorage) ExecuteSelect(stmt *ast.SelectStatement) {
	fmt.Printf("\n--- Resultado da Busca ---\n")

	columnIndices := []int{}
	if len(stmt.Columns) == 1 && stmt.Columns[0] == "*" {
		for i := 0; i < len(s.Header); i++ {
			columnIndices = append(columnIndices, i)
		}
	} else {
		for _, colName := range stmt.Columns {
			idx, ok := s.Header[strings.ToLower(colName)]
			if ok {
				columnIndices = append(columnIndices, idx)
			}
		}
	}

	for _, row := range s.Data {
		// Agora passamos a Condition (que é uma Expression) para o novo avaliador
		if stmt.Condition != nil {
			if !s.evaluate(row, stmt.Condition) {
				continue
			}
		}

		for _, idx := range columnIndices {
			fmt.Printf("%s \t", row[idx])
		}
		fmt.Println()
	}
}

// evaluate é a função recursiva que "navega" pela árvore do AST
func (s *CSVStorage) evaluate(row []string, expr ast.Expression) bool {
	switch e := expr.(type) {

	// Caso seja uma operação lógica (E / OU)
	case *ast.LogicalExpression:
		left := s.evaluate(row, e.Left)
		right := s.evaluate(row, e.Right)

		if e.Operator == "E" {
			return left && right
		}
		if e.Operator == "OU" {
			return left || right
		}

	// Caso seja uma comparação simples (idade > 20, municipio DENTRO ...)
	case *ast.ComparisonExpression:
		return s.evaluateComparison(row, e)
	}

	return false
}

func (s *CSVStorage) evaluateComparison(row []string, cond *ast.ComparisonExpression) bool {
	colIdx, ok := s.Header[strings.ToLower(cond.Left)]
	if !ok {
		return false
	}

	valorCSV := strings.TrimSpace(row[colIdx])

	if cond.Operator == "DENTRO" || cond.Operator == "EM" {
		lista, ok := cond.Right.([]string)
		if !ok {
			return false
		}
		for _, v := range lista {
			if strings.TrimSpace(v) == valorCSV {
				return true
			}
		}
		return false
	}

	valP, errP := strconv.ParseFloat(valorCSV, 64)
	valC, errC := strconv.ParseFloat(fmt.Sprintf("%v", cond.Right), 64)

	if errP == nil && errC == nil {
		switch cond.Operator {
		case ">":
			return valP > valC
		case ">=":
			return valP >= valC
		case "<":
			return valP < valC
		case "<=":
			return valP <= valC
		case "==":
			return valP == valC
		case "!=":
			return valP != valC
		}
	}

	// Tenta comparação de data (formato DD/MM/AAAA)
	layoutBR := "02/01/2006"
	dataP, errDP := time.Parse(layoutBR, valorCSV)
	dataC, errDC := time.Parse(layoutBR, fmt.Sprintf("%v", cond.Right))

	if errDP == nil && errDC == nil {
		switch cond.Operator {
		case ">":
			return dataP.After(dataC)
		case ">=":
			return !dataP.Before(dataC)
		case "<":
			return dataP.Before(dataC)
		case "<=":
			return !dataP.After(dataC)
		case "==":
			return dataP.Equal(dataC)
		case "!=":
			return !dataP.Equal(dataC)
		}
	}

	valorFiltro := fmt.Sprintf("%v", cond.Right)
	switch cond.Operator {
	case "==":
		return valorCSV == valorFiltro
	case "!=":
		return valorCSV != valorFiltro
	}
	return false
}
