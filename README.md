# Projeto Enube

## Estrutura do Projeto

```
enube_projeto/
├── backend/                 # API em Go
│   ├── cmd/                # Aplicação
│   │   ├── main.go        
│   │   └── importer/      # Importador de dados
│   ├── internal/          # Código interno da aplicação
│   │   ├── auth/         # Autenticação e JWT
│   │   ├── db/           # Config Banco
│   │   ├── handlers/     # Handlers HTTP
│   │   ├── importer/     # Lógica 
│   │   ├── middleware/   # Middlewares
│   │   ├── models/       # Modelos 
│   │   └── routes/       # Rotas 
│   ├── .env              # Variáveis
│   ├── .air.toml         # Configuração do Air (reload)
│   ├── Dockerfile        # Coontainer
│   ├── go.mod            # Dependências
│   └── init.sql          # Script 
│
├── frontend/              # Frontend em React (opcional)
│   ├── src/              # Código fonte
│   ├── public/           # Arquivos estáticos
│   └── package.json      # Dependências
│
├── base/                  # Arquivos base para importação
│   └── Reconfile fornecedores.xlsx
│
├── .env                   # Variáveis 
├── .gitignore            
├── docker-compose.yml    # Containers
└── README.md            
```

## Configuração do Ambiente

1. Clone o repositório:
```bash
git clone git@github.com:erascardsilva/enube_projeto.git
cd enube_projeto
```

2. Configure as variáveis de ambiente:
   - Copie `.env.example` para `.env` na raiz do projeto
   - Copie `.env.example` para `backend/.env`
   - Ajuste as variáveis conforme necessário

3. Inicie os containers:
```bash
docker-compose build
docker-compose up -d
```

4. Acesse a API:
   - http://localhost:8080

## Endpoints da API

### Autenticação
- `POST /auth/register` - Registro de usuário
- `POST /auth/login` - Login e obtenção do token JWT

### Usuários
- `GET /api/users` - Listar usuários
- `GET /api/users?page=1&limit=10` - Listar com paginação

### Importação
- `POST /api/import` - Importar arquivo Excel

### Consultas
- `GET /api/clients` - Listar clientes
- `GET /api/categories` - Listar categorias
- `GET /api/resources` - Listar recursos
- `GET /api/billing` - Listar faturamentos
- `GET /api/billing/summary/categories` - Resumo por categoria

## Desenvolvimento

### Backend
```bash
cd backend
go mod tidy
go run cmd/main.go
```

### Frontend (opcional)
```bash
cd frontend
npm install
npm start
```

## Notas de Implementação

- O importador de dados está otimizado para processar arquivos Excel grandes mesmo assim ainda demora
- A API usa JWT para autenticação
- O banco de dados é PostgreSQL
- A aplicação usa Docker para containerização
- Hot reload com Air para desenvolvimento

## Autor

Erasmo Cardoso da Silva
Desenvolvedor Full Stack 