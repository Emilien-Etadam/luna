#!/usr/bin/env bash
set -euo pipefail

log() {
  printf "\n[%s] %s\n" "$(date '+%Y-%m-%d %H:%M:%S')" "$*"
}

fail() {
  printf "\n[ERREUR] %s\n" "$*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "Commande manquante: $1"
}

SCRIPT_PATH="$(readlink -f "${BASH_SOURCE[0]}")"
SCRIPT_DIR="$(cd "$(dirname "${SCRIPT_PATH}")" && pwd)"

detect_repo_root() {
  if [[ -n "${LUNA_REPO_DIR:-}" ]]; then
    printf "%s\n" "${LUNA_REPO_DIR}"
    return
  fi

  if git -C "${PWD}" rev-parse --show-toplevel >/dev/null 2>&1; then
    git -C "${PWD}" rev-parse --show-toplevel
    return
  fi

  if git -C "${SCRIPT_DIR}/.." rev-parse --show-toplevel >/dev/null 2>&1; then
    git -C "${SCRIPT_DIR}/.." rev-parse --show-toplevel
    return
  fi

  if [[ -d "${HOME}/luna/.git" ]]; then
    printf "%s\n" "${HOME}/luna"
    return
  fi

  fail "Impossible de trouver le repo Luna. Definis LUNA_REPO_DIR=/chemin/vers/luna."
}

REPO_ROOT="$(detect_repo_root)"

BACKEND_SERVICE="${BACKEND_SERVICE:-luna-backend}"
FRONTEND_SERVICE="${FRONTEND_SERVICE:-luna-frontend}"
PROXY_SERVICE="${PROXY_SERVICE:-nginx}"
RESTART_PROXY="${RESTART_PROXY:-0}"
ALLOW_DIRTY="${ALLOW_DIRTY:-0}"

require_cmd git
require_cmd go
require_cmd bun
require_cmd systemctl

cd "${REPO_ROOT}"

log "Verification de l'etat git"
if [[ -n "$(git status --porcelain)" && "${ALLOW_DIRTY}" != "1" ]]; then
  fail "Le depot contient des changements locaux. Commit/stash ou relance avec ALLOW_DIRTY=1."
fi

CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
log "Branche courante: ${CURRENT_BRANCH}"

log "Mise a jour du code"
git fetch --all --tags
git pull --ff-only

log "Build backend"
cd "${REPO_ROOT}/backend"
make build

log "Mise a jour dependances frontend (bun.lock fige)"
cd "${REPO_ROOT}/frontend"
bun install --frozen-lockfile

log "Build frontend"
make build

restart_if_present() {
  local service_name="$1"
  if systemctl list-unit-files --type=service --no-legend 2>/dev/null | awk '{print $1}' | grep -Fxq "${service_name}.service"; then
    log "Redemarrage ${service_name}.service"
    systemctl restart "${service_name}"
    systemctl is-active --quiet "${service_name}" || fail "Le service ${service_name} n'est pas actif apres redemarrage."
  else
    log "Service ${service_name}.service non trouve, ignore."
  fi
}

log "Redemarrage des services applicatifs"
restart_if_present "${BACKEND_SERVICE}"
restart_if_present "${FRONTEND_SERVICE}"

if [[ "${RESTART_PROXY}" == "1" ]]; then
  restart_if_present "${PROXY_SERVICE}"
fi

log "Verification rapide API"
if command -v curl >/dev/null 2>&1; then
  if [[ -n "${API_HEALTH_URL:-}" ]]; then
    curl -fsS "${API_HEALTH_URL}" >/dev/null && log "Healthcheck OK: ${API_HEALTH_URL}" || fail "Healthcheck KO: ${API_HEALTH_URL}"
  else
    log "API_HEALTH_URL non defini, healthcheck HTTP saute."
  fi
fi

log "Update terminee avec succes"
