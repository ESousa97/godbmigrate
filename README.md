# godbmigrate

> A fast, flexible, and database-agnostic migration tool for Go projects.

![Go Report Card](https://goreportcard.com/badge/github.com/lucassousa/godbmigrate)
![Go Reference](https://pkg.go.dev/badge/github.com/lucassousa/godbmigrate.svg)
![License](https://img.shields.io/github/license/lucassousa/godbmigrate)
![Go Version](https://img.shields.io/github/go-mod/go-version/lucassousa/godbmigrate)
![Last Commit](https://img.shields.io/github/last-commit/lucassousa/godbmigrate)

---

godbmigrate is a lightweight CLI tool and Go library designed to handle database migrations with ease. It supports PostgreSQL out-of-the-box and focuses on simplicity, speed, and safety through advisory locks.

## Demonstração

### CLI Usage

```bash
# Create a new migration
godbmigrate new add_users_table

# Apply all pending migrations
godbmigrate up --dsn "postgres://user:pass@localhost:5432/dbname?sslmode=disable"

# Revert the last migration
godbmigrate down --dsn "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
```

### Library Usage

```go
import "github.com/lucassousa/godbmigrate/internal/db"

// Connect to the database
store, err := db.Connect(dsn)
if err != nil {
    log.Fatal(err)
}
defer store.Close()

// Apply pending migrations
if err := store.ApplyMigration(version, sqlContent); err != nil {
    log.Fatal(err)
}
```

## Stack Tecnológico

| Tecnologia | Papel |
|---|---|
| Go | Linguagem de programação principal |
| Cobra | Framework para criação de CLI |
| PostgreSQL | Banco de dados alvo (suporte inicial) |
| Slog | Logging estruturado nativo |

## Pré-requisitos

- Go >= 1.25.0
- PostgreSQL (ou Docker para rodar via Makefile)

## Instalação e Uso

### Como binário

```bash
go install github.com/lucassousa/godbmigrate@latest
```

### A partir do source

```bash
git clone https://github.com/lucassousa/godbmigrate.git
cd godbmigrate
make build
# Configure suas variáveis no Makefile ou via flags
make test-full
```

## Makefile Targets

| Target | Descrição |
|---|---|
| `build` | Compila o binário `godbmigrate.exe` |
| `db-up` | Inicia um container PostgreSQL via Docker |
| `db-down` | Para e remove o container PostgreSQL |
| `test-full` | Executa um ciclo completo de teste (build, new, up, status, down) |
| `clean` | Remove binários e diretório de migrations temporárias |

## Arquitetura

O projeto segue uma estrutura modular simples:
- `cmd/`: Define a interface CLI usando Cobra.
- `internal/db/`: Contém a lógica de persistência e execução de SQL.
- `migrations/`: Diretório padrão para os arquivos `.up.sql` e `.down.sql`.

Utiliza **Advisory Locks** do PostgreSQL para garantir que apenas um processo de migração execute por vez, evitando condições de corrida em ambientes distribuídos.

## API Reference

Veja a documentação completa em [pkg.go.dev](https://pkg.go.dev/github.com/lucassousa/godbmigrate).

## Configuração

| Variável | Descrição | Tipo | Padrão |
|---|---|---|---|
| `--dsn` | String de conexão PostgreSQL | string | `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable` |
| `--debug` | Habilita logs em nível DEBUG | bool | `false` |

## Roadmap

- [x] Suporte básico para PostgreSQL
- [x] Advisory Locks para concorrência
- [ ] Suporte para MySQL e SQLite
- [ ] Migrações programáticas em Go (além de SQL)
- [ ] Integração com CI/CD (GitHub Actions)

## Contribuindo

Veja [CONTRIBUTING.md](CONTRIBUTING.md) para detalhes sobre como abrir PRs e seguir padrões de código.

## Licença

Distribuído sob a licença MIT. Veja [LICENSE](LICENSE) para mais informações.

## Autor

Enoque Sousa - [Portfólio](https://enoquesousa.vercel.app) - [GitHub](https://github.com/lucassousa)
