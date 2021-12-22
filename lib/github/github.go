package github

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const (
	// HostedIDSuffix the GitHub hosted attestation type
	HostedIDSuffix = "/Attestations/GitHubHostedActions@v1"
	// SelfHostedIDSuffix the GitHub self hosted attestation type
	SelfHostedIDSuffix = "/Attestations/SelfHostedActions@v1"
	// BuildType URI indicating what type of build was performed. It determines the meaning of invocation, buildConfig and materials.
	BuildType = "https://github.com/Attestations/GitHubActionsWorkflow@v1"
	// PayloadContentType used to define the Envelope content type
	// See: https://github.com/in-toto/attestation#provenance-example
	PayloadContentType = "application/vnd.in-toto+json"
)

// Environment is the environment from which provenance is generated.
type Environment struct {
	Context *Context
	Runner  *RunnerContext
}

// NewEnvironment retrieves the github and runner contexts from environment
// variables
func NewEnvironment() (*Environment, error) {
	event, err := eventFromPath(os.Getenv("GITHUB_EVENT_PATH"))
	if err != nil {
		return nil, err
	}

	return &Environment{
		Context: &Context{
			Action:     os.Getenv("GITHUB_ACTION"),
			Actions:    (os.Getenv("GITHUB_ACTIONS") == "true"),
			ActionPath: os.Getenv("GITHUB_ACTION_PATH"),
			Actor:      os.Getenv("GITHUB_ACTOR"),
			BaseRef:    os.Getenv("GITHUB_BASE_REF"),
			Event:      event,
			EventName:  os.Getenv("GITHUB_EVENT_NAME"),
			EventPath:  os.Getenv("GITHUB_EVENT_PATH"),
			HeadRef:    os.Getenv("GITHUB_HEAD_REF"),
			Job:        os.Getenv("GITHUB_JOB"),
			Ref:        os.Getenv("GITHUB_REF"),
			Repository: os.Getenv("GITHUB_REPOSITORY"),
			RunID:      os.Getenv("GITHUB_RUN_ID"),
			RunNumber:  os.Getenv("GITHUB_RUN_NUMBER"),
			SHA:        os.Getenv("GITHUB_SHA"),
			Workflow:   os.Getenv("GITHUB_WORKFLOW"),
		},
		Runner: &RunnerContext{
			Arch:      os.Getenv("RUNNER_ARCH"),
			Name:      os.Getenv("RUNNER_NAME"),
			OS:        os.Getenv("RUNNER_OS"),
			Temp:      os.Getenv("RUNNER_TEMP"),
			ToolCache: os.Getenv("RUNNER_TOOL_CACHE"),
		},
	}, nil
}

func (e *Environment) BuilderID() string {
	repoURI := "https://github.com/" + e.Context.Repository

	if e.Context.Actions {
		return repoURI + HostedIDSuffix
	}

	return repoURI + SelfHostedIDSuffix
}

// Context holds all the information set on Github runners in relation to the job
//
// This information is retrieved from environment variables
type Context struct {
	Action     string
	Actions    bool
	ActionPath string
	Actor      string
	BaseRef    string
	Event      AnyEvent
	EventName  string
	EventPath  string
	HeadRef    string
	Job        string
	Ref        string
	Repository string
	RunID      string
	RunNumber  string
	SHA        string
	Workflow   string
}

// RunnerContext holds information about the given Github Runner in which a workflow executes
//
// This information is retrieved from environment variables
type RunnerContext struct {
	Arch      string
	Name      string
	OS        string
	Temp      string
	ToolCache string
}

// AnyEvent holds the inputs from a Github workflow
//
// See https://docs.github.com/en/actions/reference/events-that-trigger-workflows
// The only Event with dynamically-provided input is workflow_dispatch which
// exposes the user params at the key "input."
type AnyEvent struct {
	Inputs json.RawMessage `json:"inputs"`
}

func eventFromPath(path string) (AnyEvent, error) {
	event := AnyEvent{}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return event, err
	}
	if err := json.Unmarshal([]byte(f), &event); err != nil {
		return event, err
	}

	return event, nil
}
