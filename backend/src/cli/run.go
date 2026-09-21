package cli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"luna-backend/auth"
	"luna-backend/config"
	"luna-backend/crypto"
	"luna-backend/db"
	lunaerrors "luna-backend/errors"
	"luna-backend/types"

	"github.com/sirupsen/logrus"
	"golang.org/x/term"
)

type Runtime struct {
	Config *config.CommonConfig
	DB     *db.Database
	Logger *logrus.Entry
}

type runtimeOpener func() (*Runtime, error)

func MaybeRun(args []string, open runtimeOpener) error {
	return maybeRun(args, open, os.Stdout, os.Stderr, os.Stdin, os.Getenv, term.IsTerminal)
}

func maybeRun(
	args []string,
	open runtimeOpener,
	stdout io.Writer,
	stderr io.Writer,
	stdin *os.File,
	getenv func(string) string,
	isTerminal func(int) bool,
) error {
	cmd, err := ParseCommand(args)
	if err != nil {
		return err
	}

	if cmd.Name == cmdHelp {
		_, err = io.WriteString(stdout, Usage())
		return err
	}

	if open == nil {
		return fmt.Errorf("runtime manquant")
	}
	runtime, err := open()
	if err != nil {
		return err
	}

	switch cmd.Name {
	case cmdUsers:
		return runListUsers(runtime, stdout)
	case cmdResetPassword:
		return runResetPassword(runtime, cmd.Username, stdout, stderr, stdin, getenv, isTerminal)
	case cmdPromoteAdmin:
		return runPromoteAdmin(runtime, cmd.Username, stdout)
	default:
		return fmt.Errorf("commande inconnue %q", cmd.Name)
	}
}

func runListUsers(runtime *Runtime, stdout io.Writer) error {
	accounts, err := withReadTx(runtime, func(tx *db.Transaction) ([]*types.UserAccount, error) {
		accounts, tr := tx.Queries().ListUserAccounts()
		if tr != nil {
			return nil, fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
		}
		return accounts, nil
	})
	if err != nil {
		return err
	}
	_, err = io.WriteString(stdout, formatUserTable(accounts))
	return err
}

func runPromoteAdmin(runtime *Runtime, username string, stdout io.Writer) error {
	return withWriteTx(runtime, func(tx *db.Transaction) error {
		account, tr := tx.Queries().GetUserAccountByUsername(username)
		if tr != nil {
			return userLookupError(username, tr)
		}
		if account.Admin {
			_, err := fmt.Fprintf(stdout, "Le compte %s est déjà administrateur.\n", account.Username)
			return err
		}
		if tr = tx.Queries().PromoteUserToAdmin(account.Id); tr != nil {
			return fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
		}
		_, err := fmt.Fprintf(stdout, "Le compte %s est maintenant administrateur.\n", account.Username)
		return err
	})
}

func runResetPassword(
	runtime *Runtime,
	username string,
	stdout io.Writer,
	stderr io.Writer,
	stdin *os.File,
	getenv func(string) string,
	isTerminal func(int) bool,
) error {
	exists, err := crypto.SymmetricKeyFileExists(runtime.Config, "passwordPepper")
	if err != nil {
		return fmt.Errorf("impossible de vérifier le pepper : %w", err)
	}
	if !exists {
		return fmt.Errorf("fichier de pepper introuvable (%s/passwordPepper.key). Vérifie DATA_PATH : il doit être identique à celui du backend", runtime.Config.Env.GetKeysPath())
	}

	account, err := withReadTx(runtime, func(tx *db.Transaction) (*types.UserAccount, error) {
		account, tr := tx.Queries().GetUserAccountByUsername(username)
		if tr != nil {
			return nil, userLookupError(username, tr)
		}
		return account, nil
	})
	if err != nil {
		return err
	}

	stdinFd := int(stdin.Fd())
	password, err := resolveNewPassword(getenv, isTerminal(stdinFd), func() (string, error) {
		return promptNewPassword(stderr, stdinFd)
	})
	if err != nil {
		return err
	}
	if err = validatePassword(password); err != nil {
		return err
	}

	secured, tr := auth.SecurePassword(password, runtime.Config)
	if tr != nil {
		return fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
	}

	err = withWriteTx(runtime, func(tx *db.Transaction) error {
		current, lookupTr := tx.Queries().GetUserAccountByUsername(username)
		if lookupTr != nil {
			return userLookupError(username, lookupTr)
		}
		if tr = tx.Queries().UpdatePassword(current.Id, secured); tr != nil {
			return fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
		}
		if tr = tx.Queries().ForceEnableUser(current.Id); tr != nil {
			return fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
		}
		if tr = tx.Queries().DeleteSessions(current.Id); tr != nil {
			return fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
		}
		return nil
	})
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(stdout, "Mot de passe réinitialisé pour %s. Toutes les sessions de ce compte ont été invalidées.\n", account.Username)
	return err
}

func promptNewPassword(stderr io.Writer, fd int) (string, error) {
	first, err := readHiddenInput(stderr, fd, "Nouveau mot de passe : ")
	if err != nil {
		return "", err
	}
	second, err := readHiddenInput(stderr, fd, "Confirmation : ")
	if err != nil {
		return "", err
	}
	if first != second {
		return "", fmt.Errorf("les mots de passe ne correspondent pas")
	}
	return first, nil
}

func readHiddenInput(stderr io.Writer, fd int, prompt string) (string, error) {
	if _, err := io.WriteString(stderr, prompt); err != nil {
		return "", err
	}
	secret, err := term.ReadPassword(fd)
	_, _ = io.WriteString(stderr, "\n")
	if err != nil {
		return "", fmt.Errorf("lecture du mot de passe : %w", err)
	}
	return string(secret), nil
}

func userLookupError(username string, tr *lunaerrors.ErrorTrace) error {
	if tr.GetStatus() == http.StatusNotFound {
		return fmt.Errorf("utilisateur introuvable : %s\nAstuce : luna-backend users", username)
	}
	return fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
}

func formatUserTable(users []*types.UserAccount) string {
	if len(users) == 0 {
		return "Aucun compte enregistré.\nLe premier utilisateur inscrit devient administrateur.\n"
	}

	var builder strings.Builder
	writer := tabwriter.NewWriter(&builder, 0, 4, 2, ' ', 0)
	fmt.Fprintln(writer, "USERNAME\tEMAIL\tADMIN\tENABLED")
	hasAdmin := false
	for _, user := range users {
		fmt.Fprintf(writer, "%s\t%s\t%s\t%s\n", user.Username, user.Email, yesNo(user.Admin), yesNo(user.Enabled))
		if user.Admin {
			hasAdmin = true
		}
	}
	_ = writer.Flush()
	if !hasAdmin {
		builder.WriteString("\nAucun administrateur. Pour en désigner un :\n  luna-backend promote-admin <username>\n")
	}
	return builder.String()
}

func yesNo(value bool) string {
	if value {
		return "oui"
	}
	return "non"
}

func withReadTx[T any](runtime *Runtime, fn func(*db.Transaction) (T, error)) (T, error) {
	var zero T
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, tr := runtime.DB.BeginReadOnlyTransaction(ctx)
	if tr != nil {
		return zero, fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
	}
	defer func() {
		_ = tx.Rollback(runtime.Logger)
	}()

	result, err := fn(tx)
	if err != nil {
		return zero, err
	}
	if tr = tx.Commit(runtime.Logger); tr != nil {
		return zero, fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
	}
	return result, nil
}

func withWriteTx(runtime *Runtime, fn func(*db.Transaction) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, tr := runtime.DB.BeginTransaction(ctx)
	if tr != nil {
		return fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
	}
	defer func() {
		_ = tx.Rollback(runtime.Logger)
	}()

	if err := fn(tx); err != nil {
		return err
	}
	if tr = tx.Commit(runtime.Logger); tr != nil {
		return fmt.Errorf("%s", tr.Serialize(lunaerrors.LvlPlain))
	}
	return nil
}
