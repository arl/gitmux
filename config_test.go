package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
)

var _updateGolden = flag.Bool("update", false, "update golden files")

// This test ensures that new features do not change gitmux output when used
// with a default configuration.
func TestOutputNonRegression(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping in -short mode")
	}

	tmpdir, err := ioutil.TempDir("", t.Name())
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tmpdir)

	t.Logf("test working directory: %q", tmpdir)
	cloneAndHack(t, tmpdir)

	cmd := exec.Command("go", "run", ".", "-printcfg")

	b, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command error: %s: %s\noutput: %s", cmdString(cmd), err, string(b))
	}

	defcfg := path.Join(tmpdir, "default.cfg")
	if err := ioutil.WriteFile(defcfg, b, os.ModePerm); err != nil {
		t.Fatalf("Can't write %q: %s", defcfg, err)
	}

	repodir := path.Join(tmpdir, "gitmux")
	cmd = exec.Command("go", "run", ".", "-cfg", defcfg, repodir)

	got, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command %q failed:\n%s\nerr: %s", cmdString(cmd), b, err)
	}

	goldenFile := path.Join("testdata", "default.output.golden")

	if *_updateGolden {
		if err := ioutil.WriteFile(goldenFile, got, os.ModePerm); err != nil {
			t.Fatalf("Can't update golden file %q: %s", goldenFile, err)
		}
	}

	want, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(want, got) {
		t.Fatalf("got:\n%s\nwant:\n%s", want, got)
	}
}

func cmdString(cmd *exec.Cmd) string {
	return strings.Join(append([]string{cmd.Path}, cmd.Args...), " ")
}

func run(t *testing.T, name string, args ...string) {
	t.Helper()

	cmd := exec.Command(name, args...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("Command %q failed:\n%s\nerr: %s", cmdString(cmd), out, err)
	}
}

func cloneAndHack(t *testing.T, dir string) {
	t.Helper()

	popd1, err := pushdir(dir)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := popd1(); err != nil {
			t.Fatalf("popd1: %v", err)
		}
	}()

	run(t, "git", "clone", "git://github.com/arl/gitmux.git")

	popd2, err := pushdir("gitmux")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := popd2(); err != nil {
			t.Fatalf("popd2: %v", err)
		}
	}()

	if err := ioutil.WriteFile("dummy", []byte("dummy"), os.ModePerm); err != nil {
		t.Fatalf("write dummy: %s", err)
	}

	run(t, "git", "add", "dummy")
	run(t, "git", "commit", "-m", "add dummy file")

	if err := ioutil.WriteFile("dummy2", []byte("dummy2"), os.ModePerm); err != nil {
		t.Fatalf("write dummy2: %s", err)
	}

	run(t, "git", "add", "dummy2")
	run(t, "git", "stash")

	if err := ioutil.WriteFile("file1", nil, os.ModePerm); err != nil {
		t.Fatalf("write file1: %s", err)
	}

	if err := ioutil.WriteFile("file2", nil, os.ModePerm); err != nil {
		t.Fatalf("write file2: %s", err)
	}

	if err := ioutil.WriteFile("file3", nil, os.ModePerm); err != nil {
		t.Fatalf("write file3: %s", err)
	}

	run(t, "git", "add", "file1")
	run(t, "git", "add", "file2")

	if err := ioutil.WriteFile("file2", []byte("foo"), os.ModePerm); err != nil {
		t.Fatalf("write file2: %s", err)
	}
}
