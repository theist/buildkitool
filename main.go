package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/mitchellh/go-homedir"
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

func printConfig() {
	userConfigFile, _ := homedir.Expand("~/.buildkitool")

	fmt.Printf("# place and fill if needed these lines in a local file called .env\n")
	fmt.Printf("# or in your home dir as %v\n", userConfigFile)
	fmt.Printf("# or find a way to set it as environment variables\n")
	fmt.Printf("# local .env takes precedence over home file, any of these will override already setted environment variables\n\n")
	fmt.Printf("BUILDKITE_ORG=%v\n", os.Getenv("BUILDKITE_ORG"))
	fmt.Printf("BUILDKITE_API_TOKEN=%v\n", os.Getenv("BUILDKITE_API_TOKEN"))
}

func listBuilds(printJobs, printFinishedJobs bool) {
	cli := NewAPIClient(os.Getenv("BUILDKITE_API_TOKEN"), os.Getenv("BUILDKITE_ORG"))
	builds := cli.BuildList()
	if len(builds) == 0 {
		fmt.Print(color.HiGreenString("There aren't any pending builds\n"))
		return
	}
	agentsAvailable := cli.AvailableAgents()
	fmt.Printf("Listing %s builds for %s agents\n", color.YellowString("%v", len(builds)), color.YellowString("%v", agentsAvailable))
	for _, build := range builds {
		author := color.CyanString(build.Creator.Name)
		number := color.HiYellowString("#%v", *build.Number)
		pipeline := color.HiGreenString(*build.Pipeline.Name)
		branch := color.GreenString(*build.Branch)
		state := ""
		switch *build.State {
		case "running":
			state = color.HiYellowString(*build.State)
			if agentsAvailable == 0 {
				state = color.HiRedString("-stalled- no agents available")
			}
		case "scheduled":
			state = color.HiBlueString(*build.State)
		default:
			state = color.RedString(*build.State)
		}

		fmt.Printf("Build %s in %s(%s) by %s -> %s\n", number, pipeline, branch, author, state)
		if printJobs {
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
				if printFinishedJobs || *job.State == "scheduled" {
					fmt.Printf("  Job: %s -> %s\n", name, jState)
				}
			}
		}
	}
}

func commandListBuilds() {
	if len(os.Args) < 3 {
		listBuilds(true, true)
		return
	}
	flag := os.Args[2]
	switch flag {
	case "--no-jobs":
		listBuilds(false, false)
	case "--only-pending":
		listBuilds(true, false)
	}
}

func main() {
	godotenv.Load(".env")
	userConfigFile, _ := homedir.Expand("~/.buildkitool")
	godotenv.Load(userConfigFile)
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
		commandListBuilds()
	case "cancel":
		cancelBuild()
	case "agents":
		listAgents()
	case "config":
		printConfig()
	default:
		fmt.Print(color.HiRedString("Unknow command: %v\n", command))
		printHelp()
	}
}
