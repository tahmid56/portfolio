# Tahmid Akter — Portfolio

A small personal portfolio site written in Go, using the standard library's
`html/template` package for templating and `net/http` for routing — no
external web framework required.

Content (name, bio, skills, and pinned projects) is sourced from
[github.com/tahmid56](https://github.com/tahmid56).

## Pages

- **Home** (`/`) — intro, GitHub stats, and featured/pinned projects
- **About** (`/about`) — bio, skills & tech stack, philosophy
- **Contact** (`/contact`) — contact details + a working contact form (POST `/contact`)

## Project layout

```
.
├── main.go              # routes, handlers, page data
├── go.mod
├── templates/
│   ├── base.html        # shared layout (header/nav/footer)
│   ├── home.html
│   ├── about.html
│   └── contact.html
├── static/
│   └── style.css
├── Dockerfile
├── docker-compose.yml
└── .dockerignore
```

## Run locally (no Docker)

Requires Go 1.22+.

```bash
go run .
```

Visit http://localhost:8080

## Run with Docker Compose

```bash
docker compose up --build
```

Visit http://localhost:8080

The compose file builds the image from the multi-stage `Dockerfile` (Go
build stage → minimal Alpine runtime stage), exposes port 8080, and includes
a container healthcheck against `/`.

## Configuration

| Env var | Default | Description          |
|---------|---------|-----------------------|
| `PORT`  | `8080`  | Port the server binds |

## Notes

- The contact form doesn't send real email (no mail provider configured);
  submissions are logged server-side and the user sees a confirmation
  message. Wire it up to an email/API service of your choice for production use.
