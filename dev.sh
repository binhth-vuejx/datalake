#!/bin/bash
set -e

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ---------------------------------------------------------------------------
# 5.1 — Load env and clean up ports / nginx
# ---------------------------------------------------------------------------

# Load environment variables from .env and export them
export $(cat "$SCRIPT_DIR/.env" | grep -v '^#' | xargs)

# Kill processes on ports 3000, 3001, 8080, 4321, and 5000
lsof -ti:3000 | xargs kill -9 2>/dev/null || true
lsof -ti:3001 | xargs kill -9 2>/dev/null || true
lsof -ti:8080 | xargs kill -9 2>/dev/null || true
lsof -ti:4321 | xargs kill -9 2>/dev/null || true
lsof -ti:5000 | xargs kill -9 2>/dev/null || true

# Kill any existing nginx processes
pkill nginx 2>/dev/null || true

# ---------------------------------------------------------------------------
# 5.2 — Docker daemon readiness
# ---------------------------------------------------------------------------

echo "==> Starting OrbStack (for Docker daemon)..."
open -a OrbStack

echo "==> Waiting for Docker daemon to be ready..."
DOCKER_READY=false
for i in {1..30}; do
  if docker ps > /dev/null 2>&1; then
    echo "✓ Docker daemon is ready"
    DOCKER_READY=true
    break
  fi
  echo "  Waiting... ($i/30)"
  sleep 2
done

if [ "$DOCKER_READY" != "true" ]; then
  echo "ERROR: Docker daemon did not become ready within 60 seconds. Aborting."
  exit 1
fi

# ---------------------------------------------------------------------------
# 5.3 — bun install, make build, PostgreSQL, migrations
# ---------------------------------------------------------------------------

echo "==> Installing dependencies..."
bun install --cwd "$SCRIPT_DIR"

echo "==> Building server..."
make -C "$SCRIPT_DIR" build

echo "==> Starting PostgreSQL..."
docker compose -f "$SCRIPT_DIR/docker-compose.yml" up -d postgres
sleep 5

echo "==> Running migrations..."
(cd "$SCRIPT_DIR/backend_agents" && go run ./cmd/migrate up)

# ---------------------------------------------------------------------------
# 5.4 — nginx, Go backend, health-check
# ---------------------------------------------------------------------------

echo "==> Starting nginx proxy (runs in background)..."
/opt/homebrew/opt/nginx/bin/nginx -c "$SCRIPT_DIR/nginx.conf" -g "daemon off;" &
NGINX_PID=$!

echo "==> Starting backend (runs in background)..."
(cd "$SCRIPT_DIR/backend_agents" && go run ./cmd/server) &
BACKEND_PID=$!

echo "==> Waiting for backend to be ready..."
BACKEND_READY=false
for i in {1..30}; do
  if curl -sf http://localhost:8080/health > /dev/null 2>&1; then
    echo "✓ Backend is ready"
    BACKEND_READY=true
    break
  fi
  echo "  Waiting... ($i/30)"
  sleep 2
done

if [ "$BACKEND_READY" != "true" ]; then
  echo "⚠ Backend did not become ready within 60 seconds (non-fatal, continuing...)"
fi

# ---------------------------------------------------------------------------
# 5.5 — Multica daemon, Astro frontend, summary, trap
# ---------------------------------------------------------------------------

if [ "$BACKEND_READY" = "true" ]; then
  echo "==> Restarting multica daemon..."
  "$SCRIPT_DIR/backend_agents/bin/multica" daemon restart \
    && echo "✓ Daemon restarted" \
    || echo "⚠ Daemon restart failed (non-fatal, continuing...)"
fi

echo "==> Starting Astro frontend dev server (runs in background)..."
(cd "$SCRIPT_DIR/astro" && bun run dev) &
FRONTEND_PID=$!

echo ""
echo "✓ All services started!"
echo "  Frontend (via nginx): http://localhost:5000"
echo "  Backend:              http://localhost:8080"
echo ""

# Register cleanup trap for EXIT, INT, and TERM signals
trap "
  echo '==> Cleaning up...'
  pkill nginx 2>/dev/null || true
  pkill -f 'go run ./cmd/server' 2>/dev/null || true
  pkill -f 'bun run dev' 2>/dev/null || true
" EXIT INT TERM

wait $NGINX_PID
