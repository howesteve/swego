package worker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/howesteve/swego/swerker/stdio/internal/lichdata"

	"github.com/tinylib/msgp/msgp"
)

func mockOnly(t *testing.T) {
	if *useWorker {
		t.SkipNow()
	}
}

func swizzle(mock string, args ...string) func() {
	if !*useWorker {
		execCommand = testExecCommand("Test" + mock + "_SubProcess")
	}

	execCmdArgs = append(args, "-dangerous_enable_test_functions")

	return func() {
		execCommand = exec.Command
		execCmdArgs = nil
	}
}

// see https://npf.io/2015/06/testing-exec-command/
func testExecCommand(testName string) func(string, ...string) *exec.Cmd {
	return func(command string, cmdArgs ...string) *exec.Cmd {
		args := []string{"-test.run=" + testName, "--", command}
		args = append(args, cmdArgs...)
		cmd := exec.Command(os.Args[0], args...)
		cmd.Env = []string{"GO_TEST_SUBPROCESS=1"}
		return cmd
	}
}

func readInput() msgp.Raw {
	data, err := lichdata.ReadFrom(os.Stdin)
	if err != nil {
		if err == io.EOF || err == lichdata.ErrNoLength {
			os.Exit(0)
		}

		writePanic("failed to read request", err.Error())
	}

	return data
}

func writeResponse(e msgp.Encodable) {
	w := msgp.NewWriter(lichdata.NewWriter(os.Stdout))
	err := msgp.Encode(w, e)
	if err != nil {
		writePanic("failed to write response", err.Error())
	}
}

func writeErrorMap(msg string, dbg ...string) {
	em := ErrorMap{"err": msg}
	if len(dbg) != 0 {
		em["dbg"] = dbg[0]
	}

	writeResponse(em)
}

func writePanic(msg string, dbg ...string) {
	if len(dbg) != 0 {
		fmt.Fprintln(os.Stderr, "DEBUG: "+dbg[0])
	}

	fmt.Fprintln(os.Stderr, "ERROR: "+msg)
	os.Exit(1)
}

func writeInitalFuncs() {
	writeResponse(Funcs{
		"rpc_funcs", // required by worker RPC system
		"test_crash",
		"test_error",
		"worker_is_mocked",
	})
}

func TestNoFuncs_SubProcess(t *testing.T) {
	if os.Getenv("GO_TEST_SUBPROCESS") != "1" {
		t.SkipNow()
	}

	os.Exit(1)
}

func TestFuncsPanic_SubProcess(t *testing.T) {
	if os.Getenv("GO_TEST_SUBPROCESS") != "1" {
		t.SkipNow()
	}

	writePanic("funcs panic")
}

func TestInvalidFuncsData_SubProcess(t *testing.T) {
	if os.Getenv("GO_TEST_SUBPROCESS") != "1" {
		t.SkipNow()
	}

	fmt.Fprintln(os.Stdout, "invalid funcs data")
	os.Exit(1)
}

func TestInvalidFuncsType_SubProcess(t *testing.T) {
	if os.Getenv("GO_TEST_SUBPROCESS") != "1" {
		t.SkipNow()
	}

	w := msgp.NewWriter(lichdata.NewWriter(os.Stdout))
	w.WriteNil()
	w.Flush()
	os.Exit(1)
}

func TestParseFuncs_SubProcess(t *testing.T) {
	if os.Getenv("GO_TEST_SUBPROCESS") != "1" {
		t.SkipNow()
	}

	writeInitalFuncs()
	os.Exit(0)
}

func TestCall_SubProcess(t *testing.T) {
	if os.Getenv("GO_TEST_SUBPROCESS") != "1" {
		t.SkipNow()
	}

	writeInitalFuncs()
	readInput()
	writeInitalFuncs()
	os.Exit(0)
}

func TestCallError_SubProcess(t *testing.T) {
	if os.Getenv("GO_TEST_SUBPROCESS") != "1" {
		t.SkipNow()
	}

	writeInitalFuncs()
	readInput()
	writeErrorMap("test_error called", "func=test_error")
	os.Exit(0)
}

func TestCrashedCall_SubProcess(t *testing.T) {
	if os.Getenv("GO_TEST_SUBPROCESS") != "1" {
		t.SkipNow()
	}

	writeInitalFuncs()
	readInput()
	writePanic("test_crash called", "func=test_crash")
}

func TestUnexpectedExit_SubProcess(t *testing.T) {
	if os.Getenv("GO_TEST_SUBPROCESS") != "1" {
		t.SkipNow()
	}

	writeInitalFuncs()
	readInput()
	os.Exit(1)
}
