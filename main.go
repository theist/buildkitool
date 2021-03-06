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
	fmt.Println(`buildkitool

    Command line utility to check several aspects of the buildkite status.
Currently it support the next subcommands:

   * builds 		Check build status
   * config			Print current config and help configuring.
   * help			This help.

To use this command line utility you will need an API token that you can obtain
at https://buildkite.com/user/api-access-tokens. Once you have it write and
modify the output of a "buildkitool config" command.`)
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
	if !checkEnv() {
		os.Exit(1)
	}
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

func printBuildsHelp() {
	fmt.Printf("%v builds\n\n", os.Args[0])
	fmt.Println(`  Prints information of the current builds started or pending in buildkite. Each
  build will be printed with the list of jobs associated and their statuses.

  All the pending build will be presented as "scheduled" and if there's any
  running will be presented as "running" unless there aren't any agents running,
  then it will be shown as "-stalled-"

  --no-jobs, -nj, no-jobs
                    Will not show the jobs associated, only the builds

  --only-pending, -p, pending
                    Will show only pending "scheduled" jobs

  --help, help
                    This help.`)
}

func commandListBuilds() {
	if len(os.Args) < 3 {
		listBuilds(true, true)
		return
	}
	flag := os.Args[2]
	switch flag {
	case "--no-jobs", "no-jobs", "-nj":
		listBuilds(false, false)
	case "--only-pending", "pending", "-p":
		listBuilds(true, false)
	case "--help", "help":
		printBuildsHelp()
	default:
		fmt.Print(color.HiRedString("Unknow Flag: %v\n", flag))
		printBuildsHelp()
	}
}

func main() {
	godotenv.Load(".env")
	userConfigFile, _ := homedir.Expand("~/.buildkitool")
	godotenv.Load(userConfigFile)
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
		fmt.Print(color.HiRedString("Unknown command: %v\n\n", command))
		printHelp()
	}
}