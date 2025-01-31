// Similar to the simple example but filtering the hosts
package main

import (
	"os"

	"github.com/nornir-automation/gornir/pkg/gornir"
	"github.com/nornir-automation/gornir/pkg/plugins/inventory"
	"github.com/nornir-automation/gornir/pkg/plugins/logger"
	"github.com/nornir-automation/gornir/pkg/plugins/output"
	"github.com/nornir-automation/gornir/pkg/plugins/runner"
	"github.com/nornir-automation/gornir/pkg/plugins/task"
)

func main() {
	logger := logger.NewLogrus(false)

	inventory, err := inventory.FromYAMLFile("/go/src/github.com/nornir-automation/gornir/examples/hosts.yaml")
	if err != nil {
		logger.Fatal(err)
	}

	gr := &gornir.Gornir{
		Inventory: inventory,
		Logger:    logger,
	}

	// define a function we will use to filter the hosts
	filter := func(gr *gornir.Gornir, h *gornir.Host) bool {
		return h.Hostname == "dev1.group_1" || h.Hostname == "dev4.group_2"
	}

	// Before calling Gornir.RunS we call Gornir.Filter and pass the function defined
	// above. This will narrow down the inventor to the hosts matching the filter
	results, err := gr.Filter(filter).RunS(
		"What's my ip?",
		runner.Parallel(),
		&task.RemoteCommand{Command: "ip addr | grep \\/24 | awk '{ print $2 }'"},
	)
	if err != nil {
		logger.Fatal(err)
	}

	output.RenderResults(os.Stdout, results, true)
}
