package yamlvfs_test

import (
	"testing"
	"testing/fstest"

	"github.com/seb-gan/yamlvfs"
)

func TestLoad(t *testing.T) {
	fsys, err := yamlvfs.Load(`
config.yml: |
  name: myapp
  version: 1.0.0
src/:
  main.go: |
    package main
    func main() {}
  utils/:
    helper.go: |
      package utils
empty/:
`)
	if err != nil {
		t.Fatal(err)
	}

	// fstest.TestFS validates the entire fs.FS contract
	err = fstest.TestFS(fsys, "config.yml", "src/main.go", "src/utils/helper.go", "empty")
	if err != nil {
		t.Fatal(err)
	}
}
