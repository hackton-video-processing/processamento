Hackaton - Sistema de Processamento do video
O sistema tem como intuito criar um registro de um novo processamento e consulta do status do processamento
além de receber 

Integrantes do Grupo
RM354032 - Alysson Gustavo Rodrigues Maciel
RM355969 - Vinicius Duarte Mendes Nepomuceno
RM354090 - Lucas Pugliese de Morais Barros
RM353273 - Felipe Pinheiro Dantas
Para acessar o swagger e realizar os testes
Rota para acessar Swagger

Dentro do Projeto no diretório "raiz" há um arquivo com uma collection postman com todas as rotas mapeadas para teste
```
processamento.postman_collection.json
```

## Criar o processamento do video

Cria o registro inicial do processamento que será usado posteriormente para consumir as mensagens do Kafka

```url
POST http://localhost:8080/api/process:id
  {
    "files": [
        "teste.mov"
    ]
  }
```

## Buscar processamento pelo ID

Faz a busca de um processamento já inicado para verificar status do processo.

```url
GET http://localhost:8080/api/process:id
```

## Fazer o processamento do video

Faz o processamento das mensagens do Kafka

```url
POST http://localhost:8080/sink/process/video
  {
    "email": "vi.nepomuceno@outlook.com",
    "processId": "c17cfb41-3219-45fa-b394-261d78f4c8fb"
  }
```
