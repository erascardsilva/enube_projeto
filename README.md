# Enube - Projeto de Teste

## Descrição
Projeto de teste para enube (importação e análise de dados arquivo excel).

## Tecnologias
- Backend: Go (Gin ORM GORM)
- Frontend: React (em desenvolvimento)
- Banco de Dados: PostgreSQL
- Docker

## Requisitos
- Docker
- Docker Compose

## Instalação

1. Clone o repositório:
```bash
git clone https://github.com/erascardsilva/enube_projeto.git
cd enube_projeto
```

2. Inicie os containers:
```bash
docker-compose build
docker-compose up -d
```

3. Acesse a API:
```
http://localhost:8080
```

## Rotas da API

### Autenticação
- `POST /auth/register` - Registro de novo usuário
- `POST /auth/login` - Login e obtenção do token JWT

### Usuários
- `GET /api/users` - Listar todos os usuários (com paginação)

### Importação
- `POST /api/import` - Importar arquivo Excel

### Clientes
- `GET /api/clients` - Listar todos os clientes (com paginação)

### Categorias
- `GET /api/categories` - Listar todas as categorias (com paginação)

### Recursos
- `GET /api/resources` - Listar todos os recursos (com paginação)

### Faturamento
- `GET /api/billing` - Listar todos os faturamentos (com paginação)
- `GET /api/billing/summary/categories` - Resumo de faturamento por categoria
- `GET /api/billing/summary/resources` - Resumo de faturamento por recursos
- `GET /api/billing/summary/clients` - Resumo de faturamento por clientes
- `GET /api/billing/summary/months` - Resumo de faturamento por meses

## Estrutura do Projeto
```
.
├── backend/
│   ├── cmd/
│   ├── internal/
│   │   ├── auth/
│   │   ├── handlers/
│   │   ├── importer/
│   │   ├── middleware/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── routes/
│   │   └── service/
│   ├── init.sql
│   └── Dockerfile
├── frontend/
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

## Desenvolvimento

### Backend
```bash
cd backend
go mod tidy
go run main.go
```

### Frontend (em desenvolvimento)
```bash
cd frontend
npm install
npm start
```

## Testes
```bash
cd backend
go test ./...
```

## Licença
MIT

Erasmo Cardoso da Silva 