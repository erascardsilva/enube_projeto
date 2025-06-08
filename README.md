# Enube - Projeto de Teste

## Descrição
Projeto de teste para enube (importação e análise de dados arquivo excel). Uma aplicação moderna e robusta para gerenciamento e análise de dados de faturamento.

## Tecnologias

### Backend
- Go 1.23.0
  - Gin Framework 
  - GORM (ORM)
  - CORS (Cross-Origin Resource Sharing)
  - JWT (Autenticação)
  - Bcrypt (Criptografia)
  - Excelize (Manipulação de arquivos Excel)

#### Rotas da API
- **Autenticação**
  - `POST /auth/register` - Registro de novo usuário
  - `POST /auth/login` - Login e obtenção do token JWT

- **Usuários**
  - `GET /api/users` - Listar todos os usuários (com paginação)

- **Importação**
  - `POST /api/import` - Importar arquivo Excel

- **Clientes**
  - `GET /api/clients` - Listar todos os clientes (com paginação)

- **Categorias**
  - `GET /api/categories` - Listar todas as categorias (com paginação)

- **Recursos**
  - `GET /api/resources` - Listar todos os recursos (com paginação)

- **Faturamento**
  - `GET /api/billing` - Listar todos os faturamentos (com paginação)
  - `GET /api/billing/summary/categories` - Resumo de faturamento por categoria
  - `GET /api/billing/summary/resources` - Resumo de faturamento por recursos
  - `GET /api/billing/summary/clients` - Resumo de faturamento por clientes
  - `GET /api/billing/summary/months` - Resumo de faturamento por meses

### Frontend
- React 18.2.0
- Material-UI (MUI) 5.12.3
- Material Dashboard 2 React 2.2.0
- React Router DOM 6.11.0
- Fetch API (Requisições HTTP)
- JWT Authentication
- React Icons
- React Dropzone

### Banco de Dados
- PostgreSQL 17

### Containerização
- Docker
- Docker Compose

## Requisitos
- Docker
- Docker Compose

## Instalação

1. Clone o repositório:
```bash
git clone https://github.com/erascardsilva/enube_projeto.git
cd enube_projeto
na raiz e na pasta backend renomear arquivos .env.example para .env
```

2. Inicie os containers:
```bash
docker-compose build
docker-compose up -d
```

3. Acesse a aplicação:
```
Frontend: http://localhost:3000
Backend: http://localhost:8080
```

## Documentação da API

O arquivo `api.rest` na raiz do projeto contém todas as rotas da API com exemplos de uso. Este arquivo pode ser usado com a extensão REST Client do VS Code para testar as requisições diretamente do editor.

### Como usar o api.rest
1. Instale a extensão "REST Client" no VS Code
2. Abra o arquivo `api.rest`
3. Clique em "Send Request" acima de cada requisição para testá-la
4. As variáveis de ambiente (como token) são gerenciadas automaticamente

O arquivo inclui exemplos de:
- Autenticação (login/registro)
- Gerenciamento de usuários
- Importação de dados
- Consultas de categorias, recursos e faturamentos
- Resumos e relatórios

## Frontend - Funcionalidades

### Autenticação
- **Login**: Interface moderna com validação de campos e feedback visual
- **Registro**: Formulário de cadastro com validação em tempo real
- **Proteção de Rotas**: Sistema de autenticação JWT com redirecionamento automático

### Dashboard
- **Visão Geral**: Cards informativos com estatísticas do sistema
- **Resumo do Projeto**: Informações sobre tecnologias utilizadas e timeline
- **Interface Responsiva**: Adaptação perfeita para diferentes tamanhos de tela

### Importação de Dados
- **Upload de Excel**: Interface drag-and-drop para arquivos
- **Validação de Arquivos**: Verificação de formato e estrutura
- **Feedback Visual**: Progresso de upload e status de processamento

### Busca e Visualização
- **Categorias**: Listagem paginada com filtros
- **Recursos**: Visualização detalhada com paginação
- **Faturamentos**: Tabela interativa com ordenação
- **Resumos**:
  - Por Categorias
  - Por Recursos
  - Por Clientes
  - Por Meses

### Interface
- **Design Moderno**: Baseado no Material Design
- **Temas**: Suporte a tema claro/escuro
- **Componentes**: Cards, Tabelas, Gráficos e Formulários personalizados
- **Navegação**: Menu lateral responsivo
- **Feedback**: Alertas e notificações contextuais

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
│   ├── src/
│   │   ├── assets/
│   │   ├── components/
│   │   ├── contexts/
│   │   ├── layouts/
│   │   ├── pages/
│   │   ├── api/
│   │   └── utils/
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

### Frontend
```bash
cd frontend
npm install
npm start
```

## Screenshots

### Login
![Login](base/login.png)

### Dashboard
![Dashboard](base/dashboard.png)

### Importação
![Importação](base/importacao.png)

### Busca
![Busca](base/buscascategorias.png)

## Uso

1. **Acesso ao Sistema**
   - Acesse http://localhost:3000
   - Faça login com suas credenciais ou crie uma nova conta

2. **Importação de Dados**
   - Navegue até a seção "Importar"
   - Busque seu arquivo Excel ou clique para selecionar ( importante ser Reconfile fornecedores.xlsx que se encontra na pasta Base na raiz do projeto  )
   - Aguarde o processamento

3. **Análise de Dados**
   - Utilize o dashboard para visão geral
   - Explore os diferentes resumos disponíveis
   - Filtre e pesquise dados específicos

4. **Gerenciamento**
   - Visualize categorias, recursos e faturamentos
   - Utilize a paginação para navegar entre grandes conjuntos de dados
   - Exporte relatórios quando necessário

Desenvolvido por Erasmo Cardoso da Silva ( whats 11949224355 )