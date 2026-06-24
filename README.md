# Luna

**Français** — Agrégateur de calendriers auto-hébergé : une interface unique pour CalDAV, iCal et Google Calendar.

**English** — Self-hosted calendar aggregator with a single UI for CalDAV, iCal, and Google Calendar.

> **FR** — Fork personnel de [Opisek/luna](https://github.com/Opisek/luna), maintenu par [Emilien-Etadam](https://github.com/Emilien-Etadam).
>
> **EN** — Personal fork of [Opisek/luna](https://github.com/Opisek/luna), maintained by [Emilien-Etadam](https://github.com/Emilien-Etadam).

## Démarrage rapide · Quick start (Docker)

**Prérequis / Requirements :** `make`, `docker`, `docker compose`.

```bash
git clone https://github.com/Emilien-Etadam/luna.git
cd luna
# FR : adapter PUBLIC_URL et les mots de passe dans docker-compose.yml
# EN : set PUBLIC_URL and passwords in docker-compose.yml
make run
```

**FR** — Interface sur le port **8080**. Données persistantes : `/srv/luna/data` (backend), `/srv/luna/postgres` (base).

**EN** — UI on port **8080**. Persistent data: `/srv/luna/data` (backend), `/srv/luna/postgres` (database).

| Commande / Command | Action |
|--------------------|--------|
| `make run` | Build et démarrage / Build and start |
| `make up` / `make down` | Démarrer / arrêter · Start / stop |
| `make restart` | Redémarrer · Restart |
| `make build` | Rebuild sans cache · Rebuild without cache |
| `make purge` | Arrêter et **effacer** le volume Postgres · Stop and **wipe** the Postgres volume |

## Mise à jour · Update (bare metal)

**FR** — Pour un déploiement hors Docker avec services systemd :

**EN** — For non-Docker deployments using systemd services:

```bash
./scripts/install-update-command.sh   # installe / installs the `update` command
update                                # pull, build backend + frontend, redémarrage / restart
```

Variables utiles / Useful variables : `LUNA_REPO_DIR`, `API_HEALTH_URL`, `RESTART_PROXY=1`.

## Documentation

- [Déploiement / Deployment](./documentation/deployment.md)
- [API](./documentation/api.md)
- [Sécurité / Security](./documentation/security.md)

## Licence / License

**FR** — Projet upstream : copyright © 2026 Kacper Darowski (Opisek) et contributeurs. Licence à définir côté upstream.

**EN** — Upstream project: copyright © 2026 Kacper Darowski (Opisek) and contributors. License to be defined upstream.
