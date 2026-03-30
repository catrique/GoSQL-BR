package storage

import (
	"encoding/csv"
	"fmt"
	"gosql-br/internal/ast"
	"os"
	"strconv"
	"strings"
)

type CSVStorage struct {
	FileName string
	Data     [][]string
	Header   map[string]int // Mapeia "Nome" -> Coluna 0
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

	// Mapeia o cabeçalho (primeira linha)
	headerMap := make(map[string]int)
	for i, name := range records[0] {
		headerMap[strings.ToLower(name)] = i
	}

	return &CSVStorage{
		FileName: fileName,
		Data:     records[1:], // Dados sem o cabeçalho
		Header:   headerMap,
	}, nil
}

// ExecuteSelect recebe o comando do AST e filtra os dados do CSV
func (s *CSVStorage) ExecuteSelect(stmt *ast.SelectStatement) {
	fmt.Printf("\n--- Resultado da Busca ---\n")

	// 1. Identifica quais colunas o usuário quer (Nome, Salario, etc.)
	columnIndices := []int{}
	if len(stmt.Columns) == 1 && stmt.Columns[0] == "*" {
		// Se for *, pega todas as colunas mapeadas no Header
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

	// 2. Percorre cada linha dos dados
	for _, row := range s.Data {
		keepRow := true

		// 3. Se houver uma condição (QUANDO), verifica se a linha passa no teste
		if stmt.Condition != nil {
			keepRow = s.evaluateCondition(row, stmt.Condition)
		}

		// 4. Se a linha passou no filtro, imprime as colunas selecionadas
		if keepRow {
			for _, idx := range columnIndices {
				fmt.Printf("%s \t", row[idx])
			}
			fmt.Println()
		}
	}
}

func (s *CSVStorage) evaluateCondition(row []string, cond *ast.Condition) bool {
	colIdx, ok := s.Header[strings.ToLower(cond.Left)]
	if !ok {
		fmt.Printf("Aviso: Coluna '%s' não encontrada!\n", cond.Left)
		return false
	}

	valorNaPlanilha := row[colIdx]

	// Tenta comparar como número primeiro (para Salário, Idade, etc.)
	valP, errP := strconv.ParseFloat(valorNaPlanilha, 64)
	valC, errC := strconv.ParseFloat(cond.Right, 64)

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

	// Se falhar como número, compara como TEXTO (serve para DATAS e NOMES)
	// Importante: cond.Right virá do Lexer. Se for string, já vem limpa.
	switch cond.Operator {
	case "==":
		return strings.TrimSpace(valorNaPlanilha) == strings.TrimSpace(cond.Right)
	case "!=":
		return strings.TrimSpace(valorNaPlanilha) != strings.TrimSpace(cond.Right)
	}

	return false
}
