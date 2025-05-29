---

# buildkitool

A command-line tool for interacting with Buildkite, designed to streamline and automate your CI/CD workflows.

## Features

- Trigger and manage Buildkite builds from the terminal
- Query build and pipeline status
- Manage artifacts and logs
- Written in Go for speed and portability

## Installation

```sh
go install github.com/theist/buildkitool@latest
```

Or clone and build manually:

```sh
git clone https://github.com/theist/buildkitool.git
cd buildkitool
go build -o buildkitool
```

## Usage

```sh
buildkitool [command] [flags]
```

Example commands:

```sh
buildkitool trigger --pipeline my-pipeline
buildkitool status --build 12345
```

For a full list of commands and options:

```sh
buildkitool --help
```

## Configuration

Set your Buildkite API token as an environment variable:

```sh
export BUILDKITE_API_TOKEN=your_token_here
```

## Contributing

Contributions, issues, and feature requests are welcome! Please open an issue or submit a pull request.

## License

MIT License
