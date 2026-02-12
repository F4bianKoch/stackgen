# stackgen

`stackgen` is a **CLI-first tool to generate and run lightweight, Docker-based automation stacks** for business processes.

## What problem does this solve?

Internal automation projects often repeat the same setup:

- client / API / database
- Docker Compose
- environment variables
- basic project structure
- integration points for other APIs

Setting this up repeatedly is slow and error-prone.

Think of `stackgen` as a **stack generator + runner**, optimized for internal automation projects (sales, logistics, HR, IT, etc.).
`stackgen` lets you bootstrap and run these stacks in seconds, so you can focus on **business logic**, not boilerplate.

---

## Who should use it?

It is designed for engineers who want:
- fast iteration on internal tools
- full transparency over what is generated and executed
- a Linux-first workflow
- no UI magic, no black boxes

---

## Philosophy

- **CLI first** – everything is explicit and scriptable  
- **Linux first** – development and generated stacks target Linux environments  
- **Transparent** – every file and command is visible and customizable  
- **Docker-native** – stacks are isolated and reproducible  
- **Opinionated, but minimal** – sensible defaults without locking you in  

`stackgen` does not try to be a platform, dashboard, or low-code tool.

---

## Features (current & planned)

### v0.1.x
- `stackgen version` – returns current stackgen version
- `stackgen doctor` – validate local prerequisistes (Docker, Compose, etc.)
- `stackgen init <project-name>`
- Safe project scaffolding with `--force`
- Template-based stack generation with differnet template sources
- Minimal Docker Compose stack (basic: postgres service)

### v0.2.x (latest)
- reliable template engine:
  - list available embedded templates
  - for templates from various sources
  - template rendering for user customization
  - user specifies options on `stackgen init`
- `stackgen reinit` – reinitializes project with current options in manifest
- `--defaults` flag for `stackgen init` and `stackgen reinit` to use default values
- first usable templates

### v0.3.x
- improve tests
- `--dry-run` flag
- `--verbose` flag
- cleanups

### Notes
- Linux-first development

Feedback welcome!

---

## Requirements

To use `stackgen`:

- Docker (with Docker Compose v2)
- Linux, or Windows/macOS with Docker Desktop
- (Optional) Go if building from source

Run `stackgen doctor` to check if your system satisfies all requirements.

---

## Installation

### From source (for now)

```bash
git clone https://github.com/<your-username>/stackgen.git
cd stackgen
```

Run it (with go):

```bash
go run . 
go test ./...
```

Build it (with go):

```bash
go build -o ./stackgen . 
```


Binary releases and `go install` support will be added later.

### Usage

Run it:
```bash
stackgen doctor
stackgen init my-automation
cd my-automation
docker compose up -d
```

```bash
docker compose down
```


Most commands will eventually support:

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
go run . --help
```

The source code is mounted from the host; only tooling runs inside the container.

For VSCode useres there is also a `.devcontainer` provided in the Repository that does the steps above for you.

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