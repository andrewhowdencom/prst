package shell

import (
	"os"
	"path/filepath"
	"testing"
)

// TestInitScriptGolden compares InitScript output against golden files in
// testdata/. If output changes, update the golden files and verify in code
// review.
func TestInitScriptGolden(t *testing.T) {
	tests := []struct {
		name    string
		sh      Shell
		numbers []int
		golden  string
	}{
		{"bash 1", Bash, []int{1}, "bash_init_1.txt"},
		{"zsh 1", Zsh, []int{1}, "zsh_init_1.txt"},
		{"bash 0 1 2", Bash, []int{0, 1, 2}, "bash_init_012.txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.sh.InitScript(tt.numbers)

			path := filepath.Join("testdata", tt.golden)
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading golden file %q: %v", path, err)
			}

			if got != string(want) {
				t.Errorf("InitScript() mismatch (-got +want):\n--- got ---\n%s\n--- want ---\n%s", got, want)
			}
		})
	}
}
