# Projeto Enube

## Estrutura do Projeto

```
enube_projeto/
├── backend/                 # API em Go
│   ├── cmd/                # Ponto de entrada da aplicação
│   │   ├── main.go        # Arquivo principal
│   │   └── importer/      # Importador de dados
│   ├── internal/          # Código interno da aplicação
│   │   ├── auth/         # Autenticação e JWT
│   │   ├── db/           # Configuração do banco
│   │   ├── handlers/     # Handlers HTTP
│   │   ├── importer/     # Lógica de importação
│   │   ├── middleware/   # Middlewares
│   │   ├── models/       # Modelos do banco
│   │   └── routes/       # Rotas da API
│   ├── .env              # Variáveis de ambiente
│   ├── .air.toml         # Configuração do Air (hot reload)
│   ├── Dockerfile        # Configuração do container
│   ├── go.mod            # Dependências Go
│   └── init.sql          # Script de inicialização do banco
│
├── frontend/              # Frontend em React (opcional)
│   ├── src/              # Código fonte
│   ├── public/           # Arquivos estáticos
│   └── package.json      # Dependências
│
├── base/                  # Arquivos base para importação
│   └── Reconfile fornecedores.xlsx
│
├── .env                   # Variáveis de ambiente do projeto
├── .gitignore            # Arquivos ignorados pelo git
├── docker-compose.yml    # Configuração dos containers
└── README.md             # Este arquivo
```

## Configuração do Ambiente

1. Clone o repositório:
```bash
git clone https://github.com/seu-usuario/enube_projeto.git
cd enube_projeto
```

2. Configure as variáveis de ambiente:
   - Copie `.env.example` para `.env` na raiz do projeto
   - Copie `.env.example` para `backend/.env`
   - Ajuste as variáveis conforme necessário

3. Inicie os containers:
```bash
docker-compose up -d
```

4. Acesse a API:
   - http://localhost:8080

## Endpoints da API

### Autenticação
- `POST /auth/register` - Registro de usuário
- `POST /auth/login` - Login e obtenção do token JWT

### Usuários
- `GET /api/users?page=1&limit=10` - Listar usuários com paginação

### Importação
- `POST /api/import` - Importar arquivo Excel

### Consultas
- `GET /api/clients?page=1&limit=10` - Listar clientes com paginação
- `GET /api/categories?page=1&limit=10` - Listar categorias com paginação
- `GET /api/resources?page=1&limit=10` - Listar recursos com paginação
- `GET /api/billing?page=1&limit=10` - Listar faturamentos com paginação

### Agrupamentos
- `GET /api/billing/summary/categories` - Resumo por categoria
  - Retorna: categoria, total e quantidade
- `GET /api/billing/summary/resources` - Resumo por recurso
  - Retorna: URI do recurso, total e quantidade
- `GET /api/billing/summary/clients` - Resumo por cliente
  - Retorna: nome do cliente, total e quantidade
- `GET /api/billing/summary/months` - Resumo por mês
  - Retorna: mês, total e quantidade

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

- O importador de dados está otimizado para processar arquivos Excel grandes
- A API usa JWT para autenticação
- O banco de dados é PostgreSQL
- A aplicação usa Docker para containerização
- Hot reload configurado com Air para desenvolvimento
- Todos os endpoints de listagem suportam paginação
- Os agrupamentos retornam totais e quantidades

## TODO

- [ ] Otimizar performance do importador
- [ ] Adicionar mais testes
- [ ] Implementar frontend em React
- [ ] Adicionar documentação Swagger
- [ ] Implementar cache
- [ ] Adicionar mais métricas

## Autor

Erasmo Cardoso da Silva
Desenvolvedor Full Stack 