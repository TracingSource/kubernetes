package ebtables

import (
	"strings"
	"testing"

	"k8s.io/utils/exec"
	fakeexec "k8s.io/utils/exec/testing"
)

func TestEnsureChain(t *testing.T) {
	fcmd := fakeexec.FakeCmd{
		CombinedOutputScript: []fakeexec.FakeCombinedOutputAction{
			// Does not Exists
			func() ([]byte, error) { return nil, &fakeexec.FakeExitError{Status: 1} },
			// Success
			func() ([]byte, error) { return []byte{}, nil },
			// Exists
			func() ([]byte, error) { return nil, nil },
			// Does not Exists
			func() ([]byte, error) { return nil, &fakeexec.FakeExitError{Status: 1} },
			// Fail to create chain
			func() ([]byte, error) { return nil, &fakeexec.FakeExitError{Status: 2} },
		},
	}
	fexec := fakeexec.FakeExec{
		CommandScript: []fakeexec.FakeCommandAction{
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
		},
	}

	runner := New(&fexec)
	exists, err := runner.EnsureChain(TableFilter, "TEST-CHAIN")
	if exists {
		t.Errorf("expected exists = false")
	}
	if err != nil {
		t.Errorf("expected err = nil")
	}

	exists, err = runner.EnsureChain(TableFilter, "TEST-CHAIN")
	if !exists {
		t.Errorf("expected exists = true")
	}
	if err != nil {
		t.Errorf("expected err = nil")
	}

	exists, err = runner.EnsureChain(TableFilter, "TEST-CHAIN")
	if exists {
		t.Errorf("expected exists = false")
	}
	errStr := "Failed to ensure TEST-CHAIN chain: exit 2, output:"
	if err == nil || !strings.Contains(err.Error(), errStr) {
		t.Errorf("expected error: %q", errStr)
	}
}

func TestEnsureRule(t *testing.T) {
	fcmd := fakeexec.FakeCmd{
		CombinedOutputScript: []fakeexec.FakeCombinedOutputAction{
			// Exists
			func() ([]byte, error) {
				return []byte(`Bridge table: filter

Bridge chain: OUTPUT, entries: 4, policy: ACCEPT
-j TEST
`), nil
			},
			// Does not Exists.
			func() ([]byte, error) {
				return []byte(`Bridge table: filter

Bridge chain: TEST, entries: 0, policy: ACCEPT`), nil
			},
			// Fail to create
			func() ([]byte, error) { return nil, &fakeexec.FakeExitError{Status: 2} },
		},
	}
	fexec := fakeexec.FakeExec{
		CommandScript: []fakeexec.FakeCommandAction{
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
			func(cmd string, args ...string) exec.Cmd { return fakeexec.InitFakeCmd(&fcmd, cmd, args...) },
		},
	}

	runner := New(&fexec)

	exists, err := runner.EnsureRule(Append, TableFilter, ChainOutput, "-j", "TEST")
	if !exists {
		t.Errorf("expected exists = true")
	}
	if err != nil {
		t.Errorf("expected err = nil")
	}

	exists, err = runner.EnsureRule(Append, TableFilter, ChainOutput, "-j", "NEXT-TEST")
	if exists {
		t.Errorf("expected exists = false")
	}
	errStr := "Failed to ensure rule: exit 2, output: "
	if err == nil || err.Error() != errStr {
		t.Errorf("expected error: %q", errStr)
	}
}
