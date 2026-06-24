# Luna

Agrégateur de calendriers auto-hébergé — une interface unique pour CalDAV, iCal et Google Calendar.

> Fork personnel de [Opisek/luna](https://github.com/Opisek/luna), maintenu par [Emilien-Etadam](https://github.com/Emilien-Etadam).

## Démarrage rapide (Docker)

Prérequis : `make`, `docker`, `docker compose`.

```bash
git clone https://github.com/Emilien-Etadam/luna.git
cd luna
# Adapter PUBLIC_URL et les mots de passe dans docker-compose.yml
make run
```

Interface disponible sur le port **8080**. Données persistantes : `/srv/luna/data` (backend) et `/srv/luna/postgres` (base).

| Commande | Action |
|----------|--------|
| `make run` | Build et démarrage |
| `make up` / `make down` | Démarrer / arrêter |
| `make restart` | Redémarrer |
| `make build` | Rebuild sans cache |
| `make purge` | Arrêter et **effacer** le volume Postgres |

## Mise à jour (bare metal)

Pour un déploiement hors Docker avec services systemd :

```bash
./scripts/install-update-command.sh   # installe la commande `update`
update                                # pull, build backend + frontend, redémarrage
```

Variables utiles : `LUNA_REPO_DIR`, `API_HEALTH_URL`, `RESTART_PROXY=1`.

## Documentation

- [Déploiement](./documentation/deployment.md)
- [API](./documentation/api.md)
- [Sécurité](./documentation/security.md)

## Licence

Projet upstream : copyright © 2026 Kacper Darowski (Opisek) et contributeurs. Licence à définir côté upstream.
