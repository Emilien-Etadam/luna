#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TARGET_SCRIPT="${SCRIPT_DIR}/update.sh"

if [[ ! -f "${TARGET_SCRIPT}" ]]; then
  echo "Script update introuvable: ${TARGET_SCRIPT}" >&2
  exit 1
fi

chmod +x "${TARGET_SCRIPT}"

if [[ "${EUID}" -eq 0 ]]; then
  ln -sf "${TARGET_SCRIPT}" /usr/local/bin/update
  echo "Commande installee: /usr/local/bin/update"
else
  mkdir -p "${HOME}/.local/bin"
  ln -sf "${TARGET_SCRIPT}" "${HOME}/.local/bin/update"
  echo "Commande installee: ${HOME}/.local/bin/update"
  echo "Ajoute ~/.local/bin au PATH si necessaire:"
  echo 'export PATH="$HOME/.local/bin:$PATH"'
fi
