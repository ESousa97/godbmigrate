# godbmigrate

Um motor de engenharia simples e eficiente para gerenciar migrações de banco de dados em Go.

## Instalação

```bash
go get github.com/lucassousa/godbmigrate
```

## Como usar

### Criar uma nova migração
```bash
godbmigrate new <nome_da_migracao>
```

Isso gerará dois arquivos na pasta `migrations/`:
- `YYYYMMDDHHMMSS_<nome>.up.sql`
- `YYYYMMDDHHMMSS_<nome>.down.sql`

### Listar migrações
```bash
godbmigrate list
```

## Tecnologias
- Go (Golang)
- Cobra CLI
