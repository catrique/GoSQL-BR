package storage

import (
	"encoding/csv"
	"fmt"
	"gosql-br/internal/ast"
	"os"
	"sort"
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
	var filtradas [][]string
	for _, row := range s.Data {
		if stmt.Condition == nil || s.evaluate(row, stmt.Condition) {
			filtradas = append(filtradas, row)
		}
	}

	if len(stmt.Functions) > 0 {
		total := s.Data
		if stmt.IgnoreEmpty != "" {
			colIdx, ok := s.Header[strings.ToLower(stmt.IgnoreEmpty)]
			if ok {
				var semVazio [][]string
				for _, row := range s.Data {
					if colIdx < len(row) && strings.TrimSpace(row[colIdx]) != "" {
						semVazio = append(semVazio, row)
					}
				}
				total = semVazio
			}
		}

		fmt.Printf("\n--- Resultado ---\n")
		for _, fn := range stmt.Functions {
			s.executarFuncao(fn, filtradas, total)
		}
		return
	}

	columnIndices := s.resolverColunas(stmt.Columns)

	if stmt.OrderBy != "" {
		colIdx, ok := s.Header[strings.ToLower(stmt.OrderBy)]
		if ok {
			sort.SliceStable(filtradas, func(i, j int) bool {
				a := strings.TrimSpace(filtradas[i][colIdx])
				b := strings.TrimSpace(filtradas[j][colIdx])

				fa, errA := strconv.ParseFloat(a, 64)
				fb, errB := strconv.ParseFloat(b, 64)
				if errA == nil && errB == nil {
					return fa < fb
				}

				da, errDA := time.Parse(layoutBR, a)
				db, errDB := time.Parse(layoutBR, b)
				if errDA == nil && errDB == nil {
					return da.Before(db)
				}

				return a < b
			})
		}
	}

	fmt.Printf("\n--- Resultado da Busca ---\n")
	for _, row := range filtradas {
		for _, idx := range columnIndices {
			fmt.Printf("%s \t", row[idx])
		}
		fmt.Println()
	}
	fmt.Printf("Total: %d linha(s)\n", len(filtradas))
}

func (s *CSVStorage) executarFuncao(fn ast.FunctionCall, filtradas [][]string, total [][]string) {
	switch fn.Name {

	case "CONTE":
		if fn.Inner != nil && fn.Inner.Name == "DIFERENTES" {
			colIdx, ok := s.Header[strings.ToLower(fn.Inner.Column)]
			if !ok {
				fmt.Printf("CONTE(DIFERENTES): coluna '%s' não encontrada\n", fn.Inner.Column)
				return
			}
			fmt.Printf("CONTE(DIFERENTES(%s)): %d\n", fn.Inner.Column, ConteDiferentes(filtradas, colIdx))
		} else {
			fmt.Printf("CONTE: %d\n", Conte(filtradas))
		}

	case "DIFERENTES":
		colIdx, ok := s.Header[strings.ToLower(fn.Column)]
		if !ok {
			fmt.Printf("DIFERENTES: coluna '%s' não encontrada\n", fn.Column)
			return
		}
		valores := Diferentes(filtradas, colIdx)
		fmt.Printf("DIFERENTES(%s): %d valor(es)\n", fn.Column, len(valores))
		for _, v := range valores {
			fmt.Printf("  - %s\n", v)
		}

	case "MAX":
		colIdx, ok := s.Header[strings.ToLower(fn.Column)]
		if !ok {
			fmt.Printf("MAX: coluna '%s' não encontrada\n", fn.Column)
			return
		}
		val, err := Max(filtradas, colIdx)
		if err != nil {
			fmt.Printf("MAX(%s): %v\n", fn.Column, err)
			return
		}
		fmt.Printf("MAX(%s): %g\n", fn.Column, val)

	case "MIN":
		colIdx, ok := s.Header[strings.ToLower(fn.Column)]
		if !ok {
			fmt.Printf("MIN: coluna '%s' não encontrada\n", fn.Column)
			return
		}
		val, err := Min(filtradas, colIdx)
		if err != nil {
			fmt.Printf("MIN(%s): %v\n", fn.Column, err)
			return
		}
		fmt.Printf("MIN(%s): %g\n", fn.Column, val)

	case "MAX_DATA":
		colIdx, ok := s.Header[strings.ToLower(fn.Column)]
		if !ok {
			fmt.Printf("MAX_DATA: coluna '%s' não encontrada\n", fn.Column)
			return
		}
		val, err := MaxData(filtradas, colIdx)
		if err != nil {
			fmt.Printf("MAX_DATA(%s): %v\n", fn.Column, err)
			return
		}
		fmt.Printf("MAX_DATA(%s): %s\n", fn.Column, val)

	case "MIN_DATA":
		colIdx, ok := s.Header[strings.ToLower(fn.Column)]
		if !ok {
			fmt.Printf("MIN_DATA: coluna '%s' não encontrada\n", fn.Column)
			return
		}
		val, err := MinData(filtradas, colIdx)
		if err != nil {
			fmt.Printf("MIN_DATA(%s): %v\n", fn.Column, err)
			return
		}
		fmt.Printf("MIN_DATA(%s): %s\n", fn.Column, val)

	case "PORCENTAGEM":
		pct := Porcentagem(filtradas, total)
		fmt.Printf("PORCENTAGEM: %.2f%% (%d de %d)\n", pct, len(filtradas), len(total))
	}
}

func (s *CSVStorage) resolverColunas(columns []string) []int {
	var indices []int
	if len(columns) == 1 && columns[0] == "*" {
		ordenados := make([]string, len(s.Header))
		for nome, idx := range s.Header {
			ordenados[idx] = nome
		}
		for i := range ordenados {
			indices = append(indices, i)
		}
		return indices
	}
	for _, colName := range columns {
		idx, ok := s.Header[strings.ToLower(colName)]
		if ok {
			indices = append(indices, idx)
		} else {
			fmt.Printf("Aviso: coluna '%s' não encontrada\n", colName)
		}
	}
	return indices
}

func (s *CSVStorage) evaluate(row []string, expr ast.Expression) bool {
	switch e := expr.(type) {
	case *ast.LogicalExpression:
		left := s.evaluate(row, e.Left)
		right := s.evaluate(row, e.Right)
		if e.Operator == "E" {
			return left && right
		}
		if e.Operator == "OU" {
			return left || right
		}
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

	if cond.Operator == "VAZIO" {
		return valorCSV == ""
	}
	if cond.Operator == "NAO_VAZIO" {
		return valorCSV != ""
	}

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
