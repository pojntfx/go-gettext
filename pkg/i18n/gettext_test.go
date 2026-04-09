package i18n

import (
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalTranslatesString(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T, domain, dir string)
		translate func(domain, msgid string) string
	}{
		{
			name: "InitI18n() and Local()",
			setup: func(t *testing.T, domain, dir string) {
				if err := InitI18n(domain, dir, slog.Default()); err != nil {
					t.Fatalf("could not init i18n: %v", err)
				}
			},
			translate: func(_, msgid string) string { return Local(msgid) },
		},
		{
			name: "InitI18n() and L()",
			setup: func(t *testing.T, domain, dir string) {
				if err := InitI18n(domain, dir, slog.Default()); err != nil {
					t.Fatalf("could not init i18n: %v", err)
				}
			},
			translate: func(_, msgid string) string { return L(msgid) },
		},
		{
			name: "BindI18n() and LocalDomain()",
			setup: func(t *testing.T, domain, dir string) {
				if err := BindI18n(domain, dir, slog.Default()); err != nil {
					t.Fatalf("could not bind i18n: %v", err)
				}
			},
			translate: func(domain, msgid string) string { return LocalDomain(domain, msgid) },
		},
		{
			name: "BindI18n() and LD()",
			setup: func(t *testing.T, domain, dir string) {
				if err := BindI18n(domain, dir, slog.Default()); err != nil {
					t.Fatalf("could not bind i18n: %v", err)
				}
			},
			translate: func(domain, msgid string) string { return LD(domain, msgid) },
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			lcMessagesDir := filepath.Join(tmpDir, "de", "LC_MESSAGES")
			if err := os.MkdirAll(lcMessagesDir, os.ModePerm); err != nil {
				t.Fatalf("could not create temporary LC_MESSAGES dir: %v", err)
			}

			poFile := filepath.Join(tmpDir, "de.po")
			if err := os.WriteFile(poFile, []byte(`msgid ""
msgstr ""
"Content-Type: text/plain; charset=UTF-8\n"

msgid "General"
msgstr "Allgemein"`), os.ModePerm); err != nil {
				t.Fatalf("could not create .po file: %v", err)
			}

			if out, err := exec.CommandContext(t.Context(), "msgfmt", "-o", filepath.Join(lcMessagesDir, "test.mo"), poFile).CombinedOutput(); err != nil {
				t.Fatalf("could not run msgfmt (output %s): %v", out, err)
			}

			t.Setenv("LC_ALL", "de_DE.UTF-8")

			tc.setup(t, "test", tmpDir)

			assert.Equal(t, "Allgemein", tc.translate("test", "General"), "should be translated outside of a container")

			testMap := map[int]string{1: tc.translate("test", "General")}
			assert.Equal(t, "Allgemein", testMap[1], "should be translated inside a map")

			testArray := [1]string{tc.translate("test", "General")}
			assert.Equal(t, "Allgemein", testArray[0], "should be translated inside an array")

			testSlice := []string{tc.translate("test", "General")}
			assert.Equal(t, "Allgemein", testSlice[0], "should be translated inside a slice")
		})
	}
}
