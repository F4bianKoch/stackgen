# stackgen

`stackgen` is a **CLI-first tool to generate and run lightweight, Docker-based automation stacks** for business processes.

It is designed for engineers who want:
- fast iteration on internal tools
- full transparency over what is generated and executed
- a Linux-first workflow
- no UI magic, no black boxes

Think of `stackgen` as a **stack generator + runner**, optimized for internal automation projects (sales, logistics, HR, IT, etc.).

---

## Philosophy

- **CLI first** – everything is explicit and scriptable  
- **Linux first** – development and generated stacks target Linux environments  
- **Transparent** – every file and command is visible and customizable  
- **Docker-native** – stacks are isolated and reproducible  
- **Opinionated, but minimal** – sensible defaults without locking you in  

`stackgen` does not try to be a platform, dashboard, or low-code tool.

---

## What problem does this solve?

Internal automation projects often repeat the same setup:

- client / API / database
- Docker Compose
- environment variables
- basic project structure
- integration points for other APIs

Setting this up repeatedly is slow and error-prone.

`stackgen` lets you bootstrap and run these stacks in seconds, so you can focus on **business logic**, not boilerplate.

---

## Features (current & planned)

### Current
- CLI skeleton with standard commands
- Linux-first development workflow
- Docker Compose–based execution
- Clear and explicit execution model

### Planned (v0.1.0)
- `stackgen init` – generate a runnable stack from templates
- `stackgen up / down` – run and stop the stack
- `stackgen doctor` – validate local prerequisites (Docker, Compose, etc.)
- `--dry-run` and `--verbose` for full transparency
- Opinionated starter stacks (e.g. API + Postgres)
- Simple integration patterns for external APIs

Authentication, advanced workflows, and plugins are intentionally **out of scope for the first versions**.

---

## Requirements

To use `stackgen`:

- Docker (with Docker Compose v2)
- Linux, or Windows/macOS with Docker Desktop
- (Optional) Go if building from source

---

## Installation

### From source (for now)

```bash
git clone https://github.com/<your-username>/stackgen.git
cd stackgen
go build -o bin/stackgen .
```

Run it:

```bash
./bin/stackgen --help
```

Binary releases and `go install` support will be added later.

### Usage (early preview)
```bash
stackgen init my-automation
cd my-automation
stackgen up
```

```bash
stackgen doctor
stackgen down
```


Most commands support:

- `--verbose` → print all executed steps

- `--dry-run` → show what would happen without changing anything

---

## Development
`stackgen` is developed **inside a Linux container** using Docker Compose.

Start the dev environment:

```bash
docker compose -f compose.dev.yml up -d
docker compose -f compose.dev.yml exec dev bash
```

Inside the container:

```bash
go test ./...
go run . --help
```

The source code is mounted from the host; only tooling runs inside the container.

## Status

**⚠️ Early development** 

The API, command set, and templates are expected to evolve.
Breaking changes may happen until a stable v1.0 is released.

## License

This project is licensed under the **MIT License**.

This means:

- you can use it for any purpose
- you can modify it
- you can redistribute it
- you can use it commercially

You only need to keep the original copyright notice.

See the `LICENSE` file for full details.