## Lab Auction Goexpert

## Descrição
Projeto desenvolvido no lab de leilão do goexpert, iniciando um server http usando o router [gin](github.com/gin-gonic/gin) com endpoints de registros de leilão(auctions) e lances(bids), sua arquitetura abrange conceitos da Clean Architecture e temos algumas funcionalidades interessantes do Go, como channels e go routines para encerramento dos leilões e criação dos lances.

Foi proposto para adicionar o encerramento dos leilões/auctions criados pela aplicação.

### Observação
A adição do encerramento das auctions foi inserida no arquivo `creation_auction_usecase.go` ao invés do `database/creation_auction.go`, como proposto, mas por quê?
Porque isso é uma regra de negócio e o repository não deve conter tal responsabilidade.

## Testes local
### Requisitos
- Docker
- Docker Compose
- Go (para execução dos testes)

### Comandos
Os comandos são executados via Makefile para simplificar a chamada dos mesmos.

- Build da aplicação
```sh
make build
```

- Iniciar a aplicação
```sh
make up
```

- Parar a aplicação
```sh
make down
```

- Executar testes
```sh
make test
```