package cli

import (
	"errors"
	"fmt"
	"strings"
)

const newPasswordEnv = "LUNA_NEW_PASSWORD"

var (
	ErrNotCLI = errors.New("not a CLI command")
)

const (
	cmdHelp          = "help"
	cmdUsers         = "users"
	cmdResetPassword = "reset-password"
	cmdPromoteAdmin  = "promote-admin"
)

type Command struct {
	Name     string
	Username string
}

func Usage() string {
	return `luna-backend — administration

Usage :
  luna-backend                         Démarre le serveur API
  luna-backend users                   Liste les comptes (repère l'administrateur)
  luna-backend reset-password <user>   Définit un nouveau mot de passe
  luna-backend promote-admin <user>    Donne le rôle administrateur
  luna-backend help                    Affiche cette aide

Le mot de passe d'origine ne peut pas être retrouvé : Luna le stocke uniquement
sous forme de hash. reset-password en définit un nouveau.

Saisie du mot de passe :
  - interactive (TTY) : saisie masquée, demandée deux fois
  - sinon : variable d'environnement ` + newPasswordEnv + `
  - ne passe jamais le mot de passe en argument (historique shell)

Ces commandes doivent utiliser le même environnement que le backend
(Postgres et DATA_PATH, pour le pepper de hash).

Exemples (Docker) :
  docker exec -it luna-backend ./luna-backend users
  docker exec -it luna-backend ./luna-backend reset-password monidentifiant
`
}

func ParseCommand(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, ErrNotCLI
	}

	for _, arg := range args {
		if arg == "--password" || strings.HasPrefix(arg, "--password=") {
			return Command{}, errors.New("ne passe pas le mot de passe en argument (historique shell) : utilise un TTY ou " + newPasswordEnv)
		}
	}

	switch args[0] {
	case cmdHelp, "-h", "--help":
		if len(args) != 1 {
			return Command{}, fmt.Errorf("arguments inattendus pour %q", args[0])
		}
		return Command{Name: cmdHelp}, nil
	case cmdUsers:
		if len(args) != 1 {
			return Command{}, fmt.Errorf("arguments inattendus pour %q", cmdUsers)
		}
		return Command{Name: cmdUsers}, nil
	case cmdResetPassword:
		if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
			return Command{}, fmt.Errorf("usage : luna-backend %s <username>", cmdResetPassword)
		}
		if len(args) != 2 {
			return Command{}, fmt.Errorf("arguments inattendus pour %q", cmdResetPassword)
		}
		return Command{Name: cmdResetPassword, Username: args[1]}, nil
	case cmdPromoteAdmin:
		if len(args) < 2 || strings.TrimSpace(args[1]) == "" {
			return Command{}, fmt.Errorf("usage : luna-backend %s <username>", cmdPromoteAdmin)
		}
		if len(args) != 2 {
			return Command{}, fmt.Errorf("arguments inattendus pour %q", cmdPromoteAdmin)
		}
		return Command{Name: cmdPromoteAdmin, Username: args[1]}, nil
	default:
		return Command{}, fmt.Errorf("commande inconnue %q\n\n%s", args[0], Usage())
	}
}

func resolveNewPassword(getenv func(string) string, stdinIsTerminal bool, prompt func() (string, error)) (string, error) {
	if getenv == nil {
		return "", errors.New("getenv manquant")
	}

	if raw, ok := getenvLookup(getenv, newPasswordEnv); ok {
		password := strings.TrimRight(raw, "\r\n")
		if password == "" {
			return "", fmt.Errorf("%s est vide", newPasswordEnv)
		}
		return password, nil
	}

	if !stdinIsTerminal {
		return "", fmt.Errorf("pas de TTY : définis %s, ou relance avec un terminal interactif (docker exec -it)", newPasswordEnv)
	}
	if prompt == nil {
		return "", errors.New("saisie interactive indisponible")
	}
	return prompt()
}

func getenvLookup(getenv func(string) string, key string) (string, bool) {
	value := getenv(key)
	if value == "" {
		return "", false
	}
	return value, true
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("le mot de passe doit faire au moins 8 caractères")
	}
	if len(password) > 1000 {
		return fmt.Errorf("le mot de passe doit faire au plus 1000 caractères")
	}
	return nil
}
