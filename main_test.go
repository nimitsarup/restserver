package main

import (
	"flag"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/nimitsarup/restserver/features/steps"
)

type componenttestSuite struct {
}

var componentFlag = flag.Bool("component", false, "perform component tests")

func (c *componenttestSuite) InitializeScenario(ctx *godog.ScenarioContext) {
	restApiFtr, err := steps.NewUsersApiComponent()
	if err != nil {
		panic(err)
	}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		if err := restApiFtr.Reset(); err != nil {
			panic(err)
		}
	})

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		restApiFtr.Close()
	})

	restApiFtr.RegisterSteps(ctx)
}

func (t *componenttestSuite) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	// any suite initialization
}

func TestComponent(t *testing.T) {
	if *componentFlag {
		var opts = godog.Options{
			Output: colors.Colored(os.Stdout),
			Paths:  flag.Args(),
			Format: "pretty",
		}

		ts := &componenttestSuite{}

		status := godog.TestSuite{
			Name:                 "component_tests",
			ScenarioInitializer:  ts.InitializeScenario,
			TestSuiteInitializer: ts.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip()
	}
}
