package storage

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

const layoutBR = "02/01/2006"

func Conte(linhas [][]string) int {
	return len(linhas)
}

func Diferentes(linhas [][]string, colIdx int) []string {
	vistos := make(map[string]bool)
	var resultado []string
	for _, row := range linhas {
		if colIdx >= len(row) {
			continue
		}
		val := strings.TrimSpace(row[colIdx])
		if val == "" {
			continue
		}
		if !vistos[val] {
			vistos[val] = true
			resultado = append(resultado, val)
		}
	}
	sort.Strings(resultado)
	return resultado
}

func ConteDiferentes(linhas [][]string, colIdx int) int {
	return len(Diferentes(linhas, colIdx))
}

func Max(linhas [][]string, colIdx int) (float64, error) {
	maior := math.Inf(-1)
	encontrou := false
	for _, row := range linhas {
		if colIdx >= len(row) {
			continue
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(row[colIdx]), 64)
		if err != nil {
			continue
		}
		if val > maior {
			maior = val
			encontrou = true
		}
	}
	if !encontrou {
		return 0, fmt.Errorf("nenhum valor numérico encontrado na coluna")
	}
	return maior, nil
}

func Min(linhas [][]string, colIdx int) (float64, error) {
	menor := math.Inf(1)
	encontrou := false
	for _, row := range linhas {
		if colIdx >= len(row) {
			continue
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(row[colIdx]), 64)
		if err != nil {
			continue
		}
		if val < menor {
			menor = val
			encontrou = true
		}
	}
	if !encontrou {
		return 0, fmt.Errorf("nenhum valor numérico encontrado na coluna")
	}
	return menor, nil
}

func MaxData(linhas [][]string, colIdx int) (string, error) {
	var maior time.Time
	encontrou := false
	for _, row := range linhas {
		if colIdx >= len(row) {
			continue
		}
		data, err := time.Parse(layoutBR, strings.TrimSpace(row[colIdx]))
		if err != nil {
			continue
		}
		if !encontrou || data.After(maior) {
			maior = data
			encontrou = true
		}
	}
	if !encontrou {
		return "", fmt.Errorf("nenhuma data válida encontrada na coluna")
	}
	return maior.Format(layoutBR), nil
}

func MinData(linhas [][]string, colIdx int) (string, error) {
	var menor time.Time
	encontrou := false
	for _, row := range linhas {
		if colIdx >= len(row) {
			continue
		}
		data, err := time.Parse(layoutBR, strings.TrimSpace(row[colIdx]))
		if err != nil {
			continue
		}
		if !encontrou || data.Before(menor) {
			menor = data
			encontrou = true
		}
	}
	if !encontrou {
		return "", fmt.Errorf("nenhuma data válida encontrada na coluna")
	}
	return menor.Format(layoutBR), nil
}

func Porcentagem(filtradas [][]string, totalLinhas [][]string) float64 {
	total := len(totalLinhas)
	if total == 0 {
		return 0
	}
	return (float64(len(filtradas)) / float64(total)) * 100
}

func DiasEntre(val1, val2 string) int {
	val1 = strings.TrimSpace(val1)
	val2 = strings.TrimSpace(val2)
	if val1 == "" || val2 == "" {
		return -1
	}

	data1, err1 := time.Parse(layoutBR, val1)
	data2, err2 := time.Parse(layoutBR, val2)
	if err1 != nil || err2 != nil {
		return -1
	}

	if data1.After(data2) {
		data1, data2 = data2, data1
	}

	return int(data2.Sub(data1).Hours() / 24)
}
