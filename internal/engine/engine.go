package engine

import (
	"fmt"
	"gosql-br/internal/ast"
	"gosql-br/internal/storage"
)

type Engine struct {
	storage *storage.CSVStorage
}

func New() *Engine {
	return &Engine{}
}

func (e *Engine) Executar(program *ast.Program) error {
	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case *ast.UseStatement:
			store, err := storage.NewCSVStorage(s.File)
			if err != nil {
				return fmt.Errorf("erro no USE: %v", err)
			}
			e.storage = store
			fmt.Printf("Sucesso: Usando arquivo '%s'\n", s.File)

		case *ast.SelectStatement:
			if e.storage == nil {
				return fmt.Errorf("erro: execute 'USE <arquivo>' antes de 'PEGUE'")
			}
			e.storage.ExecuteSelect(s)

		default:
			if stmt.TokenLiteral() == "COLUNAS" {
				if e.storage == nil {
					return fmt.Errorf("use um arquivo primeiro")
				}
				fmt.Println("Colunas disponíveis:")
				for col := range e.storage.Header {
					fmt.Printf("- %s\n", col)
				}
				return nil
			}
			return fmt.Errorf("comando desconhecido")
		}
	}
	return nil
}
