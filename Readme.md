
# Desafio FullCycle — Cotação do Dólar em Go

Este projeto consiste em dois programas escritos em Go que se comunicam via HTTP para buscar e registrar a cotação atual do dólar (USD → BRL) utilizando a API pública da [AwesomeAPI](https://docs.awesomeapi.com.br/api-de-moedas).

---

## Estrutura do Projeto

```
/client       → Contém o client.go
/server       → Contém o server.go
cotacao.txt   → Gerado pelo cliente com a cotação atual
cotacoes.db   → Banco de dados SQLite criado automaticamente pelo servidor
```

---

## Como Executar

###  1. Rode o servidor

No terminal, acesse a pasta `server` e execute:

#### Linux/macOS:
```bash
CGO_ENABLED=1 go run server.go
```

#### Windows (PowerShell):
```powershell
$env:CGO_ENABLED="1"; go run server.go
```

O servidor ficará ouvindo na porta `8080` com o endpoint `/cotacao`.

---

### 🔹 2. Rode o cliente

Em outro terminal, acesse a pasta `client` e execute:

```bash
go run client.go
```

---

## O que acontece

- O cliente faz uma requisição HTTP para `localhost:8080/cotacao`
- O servidor busca a cotação na [AwesomeAPI](https://economia.awesomeapi.com.br/json/last/USD-BRL)
- A cotação é gravada:
  - No banco SQLite (`cotacoes.db`) via GORM
  - No arquivo `cotacao.txt` no formato:  
    ```
    Dólar: 5.1234
    ```

---

## Detalhes Técnicos

### Timeouts (usando `context.Context`)

- **400ms**: Timeout para buscar a cotação na API externa (servidor)
- **10ms**: Timeout para gravar no banco SQLite (servidor)
- **300ms**: Timeout total para o cliente obter a resposta do servidor

Todos os timeouts geram logs e respostas apropriadas caso estoure o tempo.

---

### Banco de Dados

O banco `cotacoes.db` é criado automaticamente usando GORM + SQLite.

**Tabela: `cotacoes`**  
Campos:
- `ID`
- `Code`
- `Name`
- `Bid`
- `Data`

---

## Observações sobre ambiente

Este projeto usa o driver `github.com/mattn/go-sqlite3`, que **requer CGO** ativado.  
Em alguns ambientes, pode ser necessário ter o compilador **GCC instalado**.

### Para rodar em Windows:

Caso ocorra erro relacionado a `CGO` ou `gcc` ausente, sugerimos instalar o compilador C:

**Instalador TDM-GCC (Windows):**  

 https://jmeubank.github.io/tdm-gcc/

Durante a instalação, marque a opção **“Add to PATH”**

Após isso, reinicie o terminal e rode:

```powershell
gcc --version
```

Se exibir a versão do GCC, o ambiente está pronto para executar o projeto.

---

## Dependências

- [`gorm.io/gorm`](https://gorm.io/)
- [`gorm.io/driver/sqlite`](https://gorm.io/docs/connecting_to_the_database.html)
- [`github.com/mattn/go-sqlite3`](https://github.com/mattn/go-sqlite3)

---

## Instruções para correção

1. Rode `go run server.go` na pasta `server`
2. Rode `go run client.go` na pasta `client`
3. Verifique:
   - A cotação foi impressa no terminal
   - O arquivo `cotacao.txt` foi criado com o valor
   - O banco `cotacoes.db` foi criado com os dados da cotação
4. Pode consultar o banco com qualquer visualizador SQLite, como o [DB Browser for SQLite](https://sqlitebrowser.org/)

---

## Autor

Desenvolvido por **Felício Melloni**  
Desafio da Pós-graduação **FullCycle** | Módulo 1
