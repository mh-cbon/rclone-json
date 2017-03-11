package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mh-cbon/rclone-json/rclone"
)

// Example demonstrates usage of Rclone-json api.
func Example() {

	// make a new instance of rclone.Cmd.
	cmd := rclone.New("src", "dst")

	// configure it
	// cmd.Stdout = ...
	// cmd.BinPath = ...
	// cmd.Stats = ...

	// Start the process
	cmd.Start()

	// consume the output, convertTo reads output objects, sends them to the encoder.
	cmd.ConvertTo(json.NewEncoder(os.Stdout))

	// wait for process end.
	cmd.Wait()

	fmt.Println("All done !")
}
