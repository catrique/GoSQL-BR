package main

import (
	"bufio"
	"fmt"
	"gosql-br/internal/engine"
	"gosql-br/internal/lexer"
	"gosql-br/internal/parser"
	"os"
)

func main() {
	// Iniciamos o motor
	motor := engine.New()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("GoSQL-BR Terminal v1.0")
	fmt.Println("Digite 'sair' para encerrar.")

	for {
		fmt.Print("GoSQL > ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if input == "sair" {
			break
		}

		// 1. Analisa
		l := lexer.New(input)
		p := parser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) > 0 {
			fmt.Println("Erro de sintaxe:", p.Errors())
			continue
		}

		// 2. Executa
		err := motor.Executar(program)
		if err != nil {
			fmt.Println("Erro de execução:", err)
		}
	}
}
