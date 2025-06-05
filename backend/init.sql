-- Conectar ao banco de dados
\c enube;

-- Criar extensão para UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Garantir que o usuário postgres tenha todas as permissões
GRANT ALL PRIVILEGES ON DATABASE enube TO postgres;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO postgres;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO postgres; 