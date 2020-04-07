package main

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// Mocks
type testCommandLine struct {
	mock.Mock
	exitRetCode int
	osArgs      []string
}

func (c *testCommandLine) usagePrint()      { c.Called() }
func (c *testCommandLine) adminUsagePrint() { c.Called() }
func (c *testCommandLine) exit(ret int)     { c.Called(ret) }
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
