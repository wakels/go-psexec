package main

import (
	"bufio"
	"github.com/fatih/color"
	"github.com/nproc/parseargs-go"
	"log"
	"os"
	"strings"
	"time"
)

func handleInteractiveMode() {
	reader := bufio.NewReader(os.Stdin)
	for {
		red := color.New(color.FgCyan)
		red.Printf("%s %s %s => ", time.Now().Format("15:04:05"), *serverFlag, *executorFlag)

		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(color.RedString("ERROR reading string: %s", err.Error()))
		}
		line = strings.Trim(line, "\r\n")

		if strings.EqualFold(line, "exit") {
			color.Green("Exit command received. Good bye.")
			os.Exit(0)
		}

		exeAndArgs, err := parseargs.Parse(line)
		if err != nil {
			color.Red("Cannot parse line '%s', error: %s", line, err.Error())
			continue
		}

		var exe string
		var args []string = []string{}

		exe = exeAndArgs[0]
		if len(exeAndArgs) > 1 {
			args = exeAndArgs[1:]
		}

		color.Green("Exe '%s' and args '%#v'", exe, args)
		color.Yellow("-------------------------------------")
		println()
		execute(exe, args...)
		println()
		color.Yellow("-------------------------------------")
		println()
		println()
	}
}
