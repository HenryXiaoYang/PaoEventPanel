# PaoEventPanel

[中文文档](README_CN.md)

A real-time lap-counting leaderboard system for school events. Track laps, display live rankings, and export results.

## Features

- **Live Leaderboard** - Individual and house rankings with auto-scrolling for projection displays
- **Lap Management** - Admins add/remove laps in real time with full audit trail
- **House System** - Four houses (Spring/Summer/Autumn/Winter) with custom colors and aggregated rankings
- **Theme Engine** - Independent light/dark mode config with solid/gradient/image backgrounds, glass card effects, and 4 built-in presets (Default/Ocean/Spring/Sakura)
- **Role-Based Access** - Three tiers: viewer (read-only), admin (student & lap management), super admin (full control)
- **Data Export** - Export rankings to Excel (.xlsx)
- **Weather & Clock** - Live clock with IP-based weather display

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go, Gin, GORM, JWT |
| Frontend | React, TypeScript, Vite, Tailwind CSS, shadcn/ui, Framer Motion |
| Database | SQLite (default) / PostgreSQL |
| Deployment | Docker, multi-stage build |

## Quick Start

### Docker (Recommended)

```bash
docker compose up -d
```

Visit http://localhost:8080. Default admin credentials: `admin` / `admin123`.

Customize with environment variables:

```bash
JWT_SECRET=your-secret SUPER_ADMIN_PASS=strongpass docker compose up -d
```

### Local Development

**Prerequisites**: Go 1.25+, Node.js 22+

```bash
# Clone the project
git clone https://github.com/HenryXiaoYang/PaoEventPanel.git
cd PaoEventPanel

# Set up environment variables
cp .env.example .env

# Start the backend (default :8080)
make dev-backend

# Start the frontend dev server (in another terminal)
make dev-frontend
```

### Production Build

```bash
make build
./pao-event-panel
```

## Configuration

Configure via environment variables or a `.env` file:

| Variable | Default | Description |
|---|---|---|
| `PORT` | `8080` | Server port |
| `DB_DRIVER` | `sqlite` | Database driver (`sqlite` or `postgres`) |
| `DB_DSN` | `pao_event.db` | Database connection string |
| `JWT_SECRET` | `change-me-in-production` | JWT signing secret |
| `SUPER_ADMIN_USER` | `admin` | Initial super admin username |
| `SUPER_ADMIN_PASS` | `admin123` | Initial super admin password |

PostgreSQL example:

```
DB_DRIVER=postgres
DB_DSN=host=localhost user=postgres password=postgres dbname=pao_event port=5432 sslmode=disable
```

## Project Structure

```
PaoEventPanel/
├── main.go                  # Entry point
├── config/                  # Configuration loading
├── internal/
│   ├── database/            # DB init, migration, seeding
│   ├── handler/             # HTTP handlers
│   ├── middleware/          # JWT auth, CORS
│   ├── model/               # Data models
│   ├── router/              # Route definitions
│   └── service/             # Business logic
├── web/                     # React frontend
│   ├── src/
│   │   ├── api/             # API client
│   │   ├── components/      # UI components
│   │   ├── hooks/           # Custom hooks
│   │   ├── stores/          # Zustand state management
│   │   └── types/           # TypeScript type definitions
│   └── ...
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## API

### Public Endpoints

| Method | Path | Description |
|---|---|---|
| POST | `/api/auth/login` | Authenticate and get JWT |
| GET | `/api/rankings/houses` | House rankings |
| GET | `/api/rankings/students` | Student rankings |
| GET | `/api/rankings/stats` | Aggregate stats |
| GET | `/api/activity` | Activity settings |
| GET | `/api/students` | List students |

### Admin Endpoints (JWT required)

| Method | Path | Description |
|---|---|---|
| POST | `/api/students` | Add a student |
| PUT | `/api/students/:id` | Edit a student |
| DELETE | `/api/students/:id` | Delete a student |
| POST | `/api/students/:id/laps` | Add/remove laps |
| PUT | `/api/students/:id/laps` | Set lap count |

### Super Admin Endpoints (JWT required)

| Method | Path | Description |
|---|---|---|
| PUT | `/api/activity` | Update activity settings |
| POST | `/api/upload/logo` | Upload logo |
| POST | `/api/upload/background` | Upload background image |
| GET/POST/DELETE | `/api/users` | User management |
| GET | `/api/houses` | List houses |
| PUT | `/api/houses/:id/color` | Update house color |
| GET/POST/DELETE | `/api/theme-presets` | Theme preset management |
| POST | `/api/theme-presets/builtin/apply` | Apply built-in preset |
| POST | `/api/theme-presets/:id/apply` | Apply custom preset |
| GET | `/api/export/rankings` | Export rankings to Excel |

## License

MIT
