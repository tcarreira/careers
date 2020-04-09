package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type testCommandLine struct {
	mock.Mock
	exitRetCode int
	osArgs      []string
}

func (c *testCommandLine) usagePrint(logger *log.Logger)      { c.Called() }
func (c *testCommandLine) adminUsagePrint(logger *log.Logger) { c.Called() }
func (c *testCommandLine) exit(ret int)                       { c.Called(ret) }
func (c *testCommandLine) getArg(idx int) string {
	c.Called(idx)
	return c.osArgs[idx]
}
func (c *testCommandLine) lenArgs() int {
	c.Called()
	return len(c.osArgs)
}

func TestExecutingCommandWithoutArguments(t *testing.T) {
	s := Server{}

	testComm := testCommandLine{
		exitRetCode: 1,
		osArgs: []string{
			"programName",
		},
	}

	// setup expectations
	testComm.On("lenArgs").Return(1)
	testComm.On("usagePrint")
	testComm.On("exit", 1)

	// call the code we are testing
	parseCommandLine(&testComm, &s)

	testComm.AssertExpectations(t)
}

func TestExecutingCommandAdmin(t *testing.T) {
	s := Server{}

	testComm := testCommandLine{
		exitRetCode: 1,
		osArgs: []string{
			"programName",
			"admin",
		},
	}

	// setup expectations
	testComm.On("lenArgs").Return(2)
	testComm.On("getArg", 1).Return("admin")
	testComm.On("adminUsagePrint")
	testComm.On("exit", 1)

	// call the code we are testing
	parseCommandLine(&testComm, &s)

	testComm.AssertExpectations(t)
}

func TestCommandLine_usagePrint(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	testComm := CommandLine{}

	testComm.usagePrint(logger)

	assert.Contains(t, buf.String(), fmt.Sprintf("Usage: %s COMMAND", filepath.Base(os.Args[0])))
}

func TestCommandLine_adminUsagePrint(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	testComm := CommandLine{}

	testComm.adminUsagePrint(logger)

	assert.Contains(t, buf.String(), fmt.Sprintf("Usage: %s admin COMMAND", filepath.Base(os.Args[0])))
}
