package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	. "github.com/theist/buildkitool/buildkite"
)

func checkEnv() bool {
	envNeeded := []string{"BUILDKITE_API_TOKEN", "BUILDKITE_ORG"}
	checkEnv := true
	for _, envVar := range envNeeded {
		_, present := os.LookupEnv(envVar)
		if !present {
			checkEnv = false
			log.Println("Missing needed environment variable ", envVar)
		}
	}
	return checkEnv
}

func printHelp() {
	//TODO: help!
}

func cancelBuild() {
	//TODO: cancel builds
}

func listAgents() {
	//TODO: list agetns status
}

func listBuilds() {
	cli := NewAPIClient(os.Getenv("BUILDKITE_API_TOKEN"), os.Getenv("BUILDKITE_ORG"))
	builds := cli.BuildList()
	if len(builds) == 0 {
		fmt.Print(color.HiGreenString("There aren't any pending builds\n"))
		return
	}
	fmt.Printf("Listing %s builds\n", color.YellowString("%v", len(builds)))
	for _, build := range builds {
		author := color.CyanString(build.Creator.Name)
		number := color.HiYellowString("#%v", *build.Number)
		pipeline := color.HiGreenString(*build.Pipeline.Name)
		branch := color.GreenString(*build.Branch)
		state := ""
		switch *build.State {
		case "running":
			state = color.HiYellowString(*build.State)
		case "scheduled":
			state = color.HiBlueString(*build.State)
		default:
			state = color.RedString(*build.State)
		}

		fmt.Printf("Build %s in %s(%s) by %s -> %s\n", number, pipeline, branch, author, state)
		for _, job := range build.Jobs {
			name := color.GreenString(*job.Name)
			jState := ""
			switch *job.State {
			case "running":
				jState = color.HiYellowString(*job.State)
			case "scheduled":
				jState = color.HiBlueString(*job.State)
			case "passed":
				jState = color.HiGreenString(*job.State)
			default:
				jState = color.RedString(*job.State)
			}
			fmt.Printf("  Job: %s -> %s\n", name, jState)
		}
	}
}

func main() {
	err := godotenv.Load(".env")
	if err == nil {
		log.Println("env variables loaded from .env file")
	}
	if !checkEnv() {
		os.Exit(1)
	}
	command := "help"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	switch command {
	case "help":
		printHelp()
	case "builds":
		listBuilds()
	case "cancel":
		cancelBuild()
	case "agents":
		listAgents()
	default:
		fmt.Print(color.HiRedString("Unknow command: %v\n", command))
		printHelp()
	}
}