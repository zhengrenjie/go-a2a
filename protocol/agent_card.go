package protocol

type (
	// An AgentCard conveys key information:
	// - Overall details (version, name, description, uses)
	// - Skills: A set of capabilities the agent can perform
	// - Default modalities/content types supported by the agent.
	// - Authentication requirements
	AgentCard struct {

		// Human readable name of the agent.
		// (e.g. "Recipe Agent")
		Name string `json:name`

		// A human-readable description of the agent. Used to assist users and
		// other agents in understanding what the agent can do.
		// (e.g. "Agent that helps users with recipes and cooking.")
		Description string `json:description`

		// A URL to the address the agent is hosted at.
		Url string `json:url`

		// The service provider of the agent.
		Provider *Provider `json:"provder,omitempty"`

		// The version of the agent - format is up to the provider. (e.g. "1.0.0")
		Version string `json:version`

		// A URL to documentation for the agent.
		DocumentationUrl *string `json:"documentationUrl,omitempty"`

		// Optional capabilities supported by the agent.
		Capabilities Capabilities `json:capabilities`

		// Authentication requirements for the agent.
		// Intended to match OpenAPI authentication structure.
		Authentication Authentication `json:authentication`

		// The set of interaction modes that the agent supports across all skills. This can be overridden per-skill.
		// Supported mime types for input
		DefaultInputModes []string `json:"defaultInputModes,omitempty"`

		// Supported mime types for output
		DefaultOutputModes []string `json:"defaultOutputModes,omitempty"`

		// Skills are a unit of capability that an agent can perform.
		Skills Skills `json:skills`
	}

	// The service provider of the agent.
	Provider struct {
		Organization string `json:"organization"`
		Url          string `json:url`
	}

	// Optional capabilities supported by the agent.
	Capabilities struct {
		// True if the agent supports SSE (Server-Sent Events)
		Streaming *bool `json:streaming,omitempty`

		// True if the agent can notify updates to client.
		PushNotifications *bool `json:"pushNotifications,omitempty"`

		// True if the agent exposes status change history for tasks
		StateTransitionHistory *bool `json:"stateTransitionHistory,omitempty"`
	}

	// Authentication requirements for the agent.
	// Intended to match OpenAPI authentication structure.
	Authentication struct {
		// e.g. Basic, Bearer
		Schemes []string `json:schemes`

		// Credentials a client should use for private cards
		Credentials *string `json:credentials,omitempty`
	}

	Skills struct {
		// Unique identifier for the agent's skill.
		Id string `json:id`

		// Human readable name of the skill.
		Name string `json:name`

		// Description of the skill - will be used by the client or a human.
		Description string `json:description`

		// Set of tagwords describing classes of capabilities for this specific skill (e.g. "cooking", "customer support", "billing")
		Tags []string `json:tags`

		// The set of example scenarios that the skill can perform.
		// Will be used by the client as a hint to understand how the skill can be used. (e.g. "I need a recipe for bread")
		Examples []string `json:examples,omitempty`

		// The set of interaction modes that the skill supports (if different than the default)
		InputModes  []string `json:"inputModes,omitempty"`
		OutputModes []string `json:"outputModes,omitempty"`
	}
)
