package cli

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"luna-backend/types"
)

func TestParseCommandEmptyIsNotCLI(t *testing.T) {
	_, err := ParseCommand(nil)
	if !errors.Is(err, ErrNotCLI) {
		t.Fatalf("ParseCommand(nil) error = %v, want ErrNotCLI", err)
	}

	_, err = ParseCommand([]string{})
	if !errors.Is(err, ErrNotCLI) {
		t.Fatalf("ParseCommand([]) error = %v, want ErrNotCLI", err)
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    Command
		wantErr string
	}{
		{name: "help", args: []string{"help"}, want: Command{Name: cmdHelp}},
		{name: "short help", args: []string{"-h"}, want: Command{Name: cmdHelp}},
		{name: "long help", args: []string{"--help"}, want: Command{Name: cmdHelp}},
		{name: "users", args: []string{"users"}, want: Command{Name: cmdUsers}},
		{name: "reset-password", args: []string{"reset-password", "emilien"}, want: Command{Name: cmdResetPassword, Username: "emilien"}},
		{name: "promote-admin", args: []string{"promote-admin", "emilien"}, want: Command{Name: cmdPromoteAdmin, Username: "emilien"}},
		{name: "users extra", args: []string{"users", "extra"}, wantErr: "arguments inattendus"},
		{name: "reset without user", args: []string{"reset-password"}, wantErr: "usage"},
		{name: "reset empty user", args: []string{"reset-password", "  "}, wantErr: "usage"},
		{name: "reset extra", args: []string{"reset-password", "emilien", "more"}, wantErr: "arguments inattendus"},
		{name: "password flag", args: []string{"reset-password", "--password=secret", "emilien"}, wantErr: "ne passe pas le mot de passe"},
		{name: "password flag only", args: []string{"--password"}, wantErr: "ne passe pas le mot de passe"},
		{name: "unknown", args: []string{"drop-db"}, wantErr: "commande inconnue"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseCommand(tc.args)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("error %q does not contain %q", err.Error(), tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestResolveNewPasswordPrefersEnv(t *testing.T) {
	got, err := resolveNewPassword(func(key string) string {
		if key == newPasswordEnv {
			return "correcthorse\n"
		}
		return ""
	}, true, func() (string, error) {
		t.Fatal("interactive prompt should not run when env is set")
		return "", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != "correcthorse" {
		t.Fatalf("got %q, want trimmed env password", got)
	}
}

func TestResolveNewPasswordRequiresTTYOrEnv(t *testing.T) {
	_, err := resolveNewPassword(func(string) string { return "" }, false, func() (string, error) {
		t.Fatal("prompt should not run without TTY")
		return "", nil
	})
	if err == nil || !strings.Contains(err.Error(), newPasswordEnv) {
		t.Fatalf("expected %s hint, got %v", newPasswordEnv, err)
	}
}

func TestResolveNewPasswordInteractive(t *testing.T) {
	got, err := resolveNewPassword(func(string) string { return "" }, true, func() (string, error) {
		return "interactive-secret", nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != "interactive-secret" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatUserTable(t *testing.T) {
	empty := formatUserTable(nil)
	if !strings.Contains(empty, "Aucun compte") {
		t.Fatalf("empty table = %q", empty)
	}

	withAdmin := formatUserTable([]*types.UserAccount{
		{Username: "emilien", Email: "emilien@example.com", Admin: true, Enabled: true},
		{Username: "guest", Email: "guest@example.com", Admin: false, Enabled: false},
	})
	if !strings.Contains(withAdmin, "emilien") || !strings.Contains(withAdmin, "oui") {
		t.Fatalf("table missing admin row: %q", withAdmin)
	}
	if strings.Contains(withAdmin, "Aucun administrateur") {
		t.Fatalf("unexpected missing-admin warning: %q", withAdmin)
	}

	withoutAdmin := formatUserTable([]*types.UserAccount{
		{Username: "guest", Email: "guest@example.com", Admin: false, Enabled: true},
	})
	if !strings.Contains(withoutAdmin, "promote-admin") {
		t.Fatalf("expected promote-admin hint, got %q", withoutAdmin)
	}
}

func TestUsageMentionsReset(t *testing.T) {
	usage := Usage()
	for _, needle := range []string{cmdUsers, cmdResetPassword, cmdPromoteAdmin, newPasswordEnv} {
		if !strings.Contains(usage, needle) {
			t.Fatalf("usage missing %q", needle)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	if err := validatePassword("short"); err == nil {
		t.Fatal("expected error for short password")
	}
	if err := validatePassword("long-enough-password"); err != nil {
		t.Fatal(err)
	}
}

func TestMaybeRunHelpDoesNotOpenRuntime(t *testing.T) {
	err := maybeRun([]string{"help"}, func() (*Runtime, error) {
		t.Fatal("help must not open the database")
		return nil, nil
	}, io.Discard, io.Discard, os.Stdin, os.Getenv, func(int) bool { return false })
	if err != nil {
		t.Fatal(err)
	}
}
