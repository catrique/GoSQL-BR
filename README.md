# GoSQL-BR

O **GoSQL-BR** é um interpretador de consultas SQL com comandos em português, desenvolvido em Go. O objetivo do projeto é processar arquivos CSV de forma simples através de um terminal interativo.

O sistema foi estruturado seguindo as etapas de um compilador real: análise léxica (Lexer), análise sintática (Parser) e um motor de execução (Engine).

### Como executar

1. Certifique-se de ter o Go instalado.
2. Coloque seu arquivo CSV na raiz do projeto (ex: `dados.csv`).
3. No terminal, execute:
   ```bash
   go run cmd/main.go
   ```

### Comandos Suportados

Após iniciar o programa, você pode utilizar os seguintes comandos no prompt:

* **USE**: Define o arquivo de dados que será utilizado.
    * Exemplo: `USE dados.csv`
* **COLUNAS**: Lista todos os cabeçalhos identificados no arquivo carregado.
* **PEGUE**: Seleciona colunas específicas ou todas utilizando `*`.
    * Exemplo: `PEGUE Nome, Salario QUANDO Data == '21/10/2025'`
    * Exemplo: `PEGUE * QUANDO Salario > 5000`

### Estrutura do Código

* **cmd/**: Contém o ponto de entrada (main.go) e a interface de linha de comando.
* **internal/lexer/**: Transforma a entrada de texto em tokens processáveis.
* **internal/parser/**: Valida a gramática e constrói a Árvore de Sintaxe Abstrata (AST).
* **internal/ast/**: Define as estruturas de dados que representam os comandos.
* **internal/engine/**: Coordena a execução dos comandos validados.
* **internal/storage/**: Gerencia a leitura física e filtragem dos arquivos CSV.
```