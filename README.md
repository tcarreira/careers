# superhero
SuperHero API - Go (inspired by superheroapi.com)

This is being made in the context of https://github.com/levpay/careers#desafio


# Running


## Setup Environment (first time: create schema)

With docker-compose (this way, docker is the only dependency)
```
docker-compose up --build -d # first time takes longer
docker-compose exec api /superhero admin schema  # first time only
```

Without docker, you need to have Go installed and a PostgreSQL database ready
```
go build
./superhero admin schema # first time only
./superhero serve 
# ./superhero serve swagger # in order to turn on the endpoint http://localhost:8080/index.html

```

if running with docker-compose or serve swagger, access the Swagger UI for testing the REST API: http://localhost:8080/swagger/index.html



## Testing

Run unit tests without a PostgeSQL database (on every save): 
```
go test -v ./...
```

## testing with a database

For performing full tests (with a real database), run it with docker:
```
 docker run --rm --name pgsql -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres:12-alpine

 go test -v ./... -tags sql -cover

 docker stop pgsql
```

If you have a PostgtreSQL instance you want to run test against:

```
DB_HOST=db DB_USER=user DB_PASS=pass go test -v ./... -tags sql
``` 

------

# Features

- [X] Create new Super(Hero/Vilan)
- [X] Get Super list
- [X] Get Super(Heroes) list
- [X] Get Super(Vilans) list
- [X] Search by name
- [X] Search by uuid
- [X] Delete Super
- [X] Super Groups


------------------------------------

== desafio ==

------------------------------------




### Gerais
Através da API deve ser possível:
- Cadastrar um Super/Vilão
```
curl -X POST "http://localhost:8080/api/v1/supers" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"name\": \"name1\", \"type\": \"HERO\"}"
```

- Listar todos os Super's cadastrados
```
curl -X GET "http://localhost:8080/api/v1/supers" -H "accept: application/json"
```

- Listar apenas os Super Heróis
```
curl -X GET "http://localhost:8080/api/v1/supers?type=hero" -H "accept: application/json"
```

- Listar apenas os Super Vilões
```
curl -X GET "http://localhost:8080/api/v1/supers?type=vilan" -H "accept: application/json"
```

- Buscar por nome
```
curl -X GET "http://localhost:8080/api/v1/supers/name1" -H "accept: application/json"
```

- Buscar por 'uuid'
```
curl -X GET "http://localhost:8080/api/v1/supers/38ccc4e8-2501-4f45-87a1-817bc452c393" -H "accept: application/json"
```

- Remover o Super
```
curl -X DELETE "http://localhost:8080/api/v1/supers/name1" -H "accept: application/json"
```


### Específicos
- API deve seguir a arquitetura [REST](https://restfulapi.net/)
```
/
├── api
│   └── v1
│       ├── groups
│       │   └── name
│       ├── super-hero
│       ├── supers
│       │   └── name
│       └── super-vilan
└── swagger
    ├── index.html
    └── doc.json

```

- API deve seguir os principios do [12 factor app](https://12factor.net/pt_br/)
   
    `OK`

- Cada super deve ser cadastrado somente a partir do seu `name`.

```
curl -X POST "http://localhost:8080/api/v1/super-hero" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"name\": \"hero1\"}"


curl -X POST "http://localhost:8080/api/v1/super-vilan" -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"name\": \"vilan1\"}"
```


- A pesquisa por um super deve conter os seguintes campos:
    - uuid - `OK`
    - name - `OK`
    - full name - `OK`
    - intelligence - `OK - int as string`
    - power - `OK - int as string`
    - occupation - `OK`
    - image - `OK (url)`
    - type - `ADDED for coherence`
- A pesquisa por um super também precisa conter:
    - lista de grupos em que tal super está associado - `OK - list of groups names`
    - número de parentes - `OK - relatives are defined as other unique Super(Hero/Vilan) that are part of the same groups than I'm part of`


## Como será a avaliação

A ideia aqui é entender como você toma suas decisões frente a certas adversidades e como você desenvolve através de multiplas funcionalidades.

Pontos que vamos avaliar:
- Commits
    - como você evoluiu seu pensamento durante o projeto, pontualidade e clareza.

    `cometi alguns erros no início do projecto que foram muito mais complexos de resolver no fim (eg: organização por packages, injeção de dependências)`
    
- Testes
    - Quanto mais testes melhor! Vide https://code.tutsplus.com/pt/tutorials/lets-go-testing-golang-programs--cms-26499 .

    `poderia ter feito mais TDD. `

- Complexidade
    - Código bom é código legivel e simples (https://medium.com/trainingcenter/golang-d94e16d4b383).

    `gostaria de obter feedback sobre este ponto. Muito mais fácil de saber a verdade com alguma versão de pair-programming`

- Dependências
    - O ecosistema (https://github.com/avelino/awesome-go) da linguagem possui diversas ferramentas para serem usadas, use-as bem!

    ```
    Dependências principais:
    - GO >= 1.13 (swagger) ; GO >=1.11 (go mod)
    - gin-gonic : Framework HTTP para criar uma API REST decente
    - go-pg/v9 : ORM para PostgreSQL

    - testify : conjunto de ferramentas para facilitar testes unitários (incluindo Mocks)
    - swaggo : framework para documentação e testes da REST API (compatível com a especificação open-api)
    ```


- Documentação `ver acima`

    - Qual versão de Go você usou?
    - Quais bibliotecas e ferramentas usou?
    - Como se utiliza a sua aplicação?
    - Como executamos os testes?
    

- Considerações
    - as regras de negócio não foram definidas intencionalmente
    - cabe a você decidir como vai manter os cadastros no banco da aplicação

    ```
    - não foram consideradas relações fortes em BD (não existem Foreign Keys)
    - não foram consideradas colunas que todas as entidades deveriam ter (active, created_at, updated_at).
    - não foi implementado nenhum método de migração de dados/schema
    ``` 

    - cabe a você decidir como vai tratar cadastros repetidos

    `não são permitidos. A edição deverá usar um endpoint diferente: PUT /supers`