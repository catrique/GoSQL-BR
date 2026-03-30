```markdown
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

### Guia de Comandos

Abaixo estão todos os comandos suportados pelo interpretador:

| Comando | Descrição | Exemplo de Uso |
| :--- | :--- | :--- |
| **USE** | Define e carrega o arquivo CSV para a memória. | `USE dados.csv` |
| **COLUNAS** | Exibe todos os nomes de colunas (cabeçalho) do arquivo atual. | `COLUNAS` |
| **PEGUE** | Seleciona colunas específicas de cada linha. | `PEGUE Nome, Cargo` |
| **PEGUE \*** | Seleciona todas as colunas disponíveis no arquivo. | `PEGUE *` |
| **QUANDO** | Filtra os resultados baseados em uma condição lógica. | `PEGUE Nome QUANDO Idade > 25` |
| **sair** | Encerra o terminal interativo do GoSQL-BR. | `sair` |

### Operadores de Comparação

Ao utilizar a cláusula **QUANDO**, você pode usar os seguintes operadores:

* `==` : Igual a
* `!=` : Diferente de
* `>`  : Maior que
* `<`  : Menor que
* `>=` : Maior ou igual a
* `<=` : Menor ou igual a

### Estrutura do Código

* **cmd/**: Contém o ponto de entrada (main.go) e a interface de linha de comando.
* **internal/lexer/**: Transforma a entrada de texto em tokens processáveis.
* **internal/parser/**: Valida a gramática e constrói a Árvore de Sintaxe Abstrata (AST).
* **internal/ast/**: Define as estruturas de dados que representam os comandos.
* **internal/engine/**: Coordena a execução dos comandos validados.
* **internal/storage/**: Gerencia a leitura física e filtragem dos arquivos CSV.
```