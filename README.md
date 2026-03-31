# GoSQL-BR

O **GoSQL-BR** é um interpretador de consultas SQL com comandos em português, desenvolvido em Go. O objetivo do projeto é processar arquivos CSV de forma simples através de um terminal interativo.

O sistema foi estruturado seguindo as etapas de um compilador real: análise léxica (Lexer), análise sintática (Parser) e um motor de execução (Engine).

---

### Como executar

1. Certifique-se de ter o Go instalado.
2. Coloque seu arquivo CSV na raiz do projeto (ex: `dados.csv`).
3. No terminal, execute:
   ```bash
   go run cmd/main.go
   ```

---

### Guia de Comandos

Abaixo estão todos os comandos suportados pelo interpretador:

| Comando | Descrição | Exemplo de Uso |
| :--- | :--- | :--- |
| **USE** | Define e carrega o arquivo CSV para a memória. | `USE dados.csv` |
| **COLUNAS** | Exibe todos os nomes de colunas (cabeçalho) do arquivo atual. | `COLUNAS` |
| **PEGUE** | Seleciona colunas específicas de cada linha. | `PEGUE nome, cidade` |
| **PEGUE \*** | Seleciona todas as colunas disponíveis no arquivo. | `PEGUE *` |
| **QUANDO** | Filtra os resultados com base em uma condição lógica. | `PEGUE nome QUANDO idade > 25` |
| **ORDENE POR** | Ordena os resultados por uma coluna (texto, número ou data). | `PEGUE * ORDENE POR data` |
| **IGNORAR VAZIO** | Exclui do denominador as linhas com campo vazio ao usar `PORCENTAGEM`. | `PEGUE PORCENTAGEM QUANDO status == 'ativo' IGNORAR VAZIO status` |
| **sair** | Encerra o terminal interativo do GoSQL-BR. | `sair` |

---

### Operadores de Comparação

Ao utilizar a cláusula **QUANDO**, você pode usar os seguintes operadores:

* `==` : Igual a
* `!=` : Diferente de
* `>`  : Maior que
* `<`  : Menor que
* `>=` : Maior ou igual a
* `<=` : Menor ou igual a
* `EM / DENTRO` : Verifica se um valor pertence a uma lista.
  * *Exemplo:* `PEGUE * QUANDO cidade DENTRO ('Campinas', 'Santos', 'Sorocaba')`
* `== VAZIO` : Verifica se o campo está vazio.
  * *Exemplo:* `PEGUE nome QUANDO email == VAZIO`
* `!= VAZIO` : Verifica se o campo está preenchido.
  * *Exemplo:* `PEGUE nome QUANDO email != VAZIO`

---

### Operadores Lógicos e Agrupamento

* `E` / `OU` : Permite combinar múltiplas condições.
* `( )` : Define a precedência das operações, garantindo que a lógica interna seja resolvida primeiro.
  * *Exemplo:* `PEGUE nome QUANDO status == 'ativo' E (data >= '01/01/2024' E data <= '31/03/2024')`

---

### Datas

O GoSQL-BR reconhece automaticamente datas no formato **`DD/MM/AAAA`** e realiza comparações cronológicas corretas — não comparações de texto.

```
PEGUE * QUANDO data_cadastro >= '01/01/2024'
PEGUE * QUANDO data_venda > '31/12/2023' E data_venda <= '30/06/2024'
PEGUE CONTE QUANDO data_entrega < '15/03/2024'
```

---

### Funções de Agregação

Funções que calculam um resultado sobre o conjunto de linhas filtradas. Utilizadas após o `PEGUE`.

| Função | Descrição | Exemplo de Uso |
| :--- | :--- | :--- |
| **CONTE** | Conta quantas linhas atendem à condição. | `PEGUE CONTE QUANDO status == 'ativo'` |
| **CONTE(DIFERENTES(coluna))** | Conta quantos valores únicos existem em uma coluna. | `PEGUE CONTE(DIFERENTES(cidade))` |
| **DIFERENTES(coluna)** | Lista os valores únicos de uma coluna em ordem alfabética. | `PEGUE DIFERENTES(categoria)` |
| **MAX(coluna)** | Retorna o maior valor numérico de uma coluna. | `PEGUE MAX(preco)` |
| **MIN(coluna)** | Retorna o menor valor numérico de uma coluna. | `PEGUE MIN(preco)` |
| **MAX_DATA(coluna)** | Retorna a data mais recente de uma coluna. | `PEGUE MAX_DATA(data_cadastro)` |
| **MIN_DATA(coluna)** | Retorna a data mais antiga de uma coluna. | `PEGUE MIN_DATA(data_cadastro)` |
| **PORCENTAGEM** | Calcula a % de linhas que atendem à condição em relação ao total. | `PEGUE PORCENTAGEM QUANDO status == 'ativo'` |

---

### Funções de Data

Funções utilizadas dentro do **QUANDO** que operam linha a linha.

| Função | Descrição | Exemplo de Uso |
| :--- | :--- | :--- |
| **DIAS_ENTRE(coluna1, coluna2)** | Calcula a diferença em dias entre duas colunas de data. A ordem das colunas não importa — o resultado é sempre positivo. Linhas com datas vazias são ignoradas. | `PEGUE * QUANDO DIAS_ENTRE(data_inicio, data_fim) >= 7` |

*Exemplo combinado:*
```
PEGUE nome, status QUANDO DIAS_ENTRE(data_abertura, data_fechamento) > 30 E status == 'aberto'
```

---

### Exemplos de Uso

**Listar registros filtrados por categoria e faixa de preço:**
```
USE produtos.csv
PEGUE nome, preco, estoque QUANDO categoria == 'eletrônicos' E preco > 100
```

**Contar registros em múltiplas cidades:**
```
USE clientes.csv
PEGUE CONTE QUANDO status == 'ativo' E cidade DENTRO ('São Paulo', 'Campinas', 'Santos')
```

**Ver valores únicos de uma coluna:**
```
USE vendas.csv
PEGUE DIFERENTES(categoria)
```

**Quantas categorias distintas têm estoque baixo:**
```
USE produtos.csv
PEGUE CONTE(DIFERENTES(categoria)) QUANDO estoque < 5
```

**Percentual de registros com um status, excluindo campos vazios do cálculo:**
```
USE pedidos.csv
PEGUE PORCENTAGEM QUANDO status == 'entregue' IGNORAR VAZIO status
```

**Registros com mais de 10 dias de intervalo entre duas datas, ordenados:**
```
USE pedidos.csv
PEGUE cliente, data_abertura, data_fechamento QUANDO DIAS_ENTRE(data_abertura, data_fechamento) > 10 ORDENE POR data_abertura
```

**Maior e menor valor de uma coluna numérica com filtro:**
```
USE produtos.csv
PEGUE MAX(preco), MIN(preco) QUANDO categoria == 'eletrônicos'
```

**Registros dentro de um intervalo de datas ordenados:**
```
USE vendas.csv
PEGUE * QUANDO data_venda >= '01/01/2024' E data_venda <= '31/03/2024' ORDENE POR data_venda
```

**Condição composta com agrupamento por parênteses:**
```
USE registros.csv
PEGUE nome QUANDO (status == 'pendente' OU status == 'aberto') E data_criacao >= '01/06/2024'
```

---

### Estrutura do Código

* **cmd/**: Contém o ponto de entrada (`main.go`) e a interface de linha de comando.
* **internal/lexer/**: Transforma a entrada de texto em tokens processáveis.
* **internal/parser/**: Valida a gramática e constrói a Árvore de Sintaxe Abstrata (AST).
* **internal/ast/**: Define as estruturas de dados que representam os comandos.
* **internal/engine/**: Coordena a execução dos comandos validados.
* **internal/storage/**: Gerencia a leitura, filtragem e execução de consultas nos arquivos CSV.
  * `csv.go` — leitura do arquivo e avaliação das condições linha a linha.
  * `functions.go` — funções de agregação e cálculo (`CONTE`, `MAX`, `DIAS_ENTRE`, etc).