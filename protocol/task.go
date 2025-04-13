package protocol

const (
	TaskStateSubmitted     TaskState = "submitted"
	TaskStateWorking       TaskState = "working"
	TaskStateInputRequired TaskState = "input-required"
	TaskStateCompleted     TaskState = "completed"
	TaskStateCanceled      TaskState = "canceled"
	TaskStateFailed        TaskState = "failed"
	TaskStateUnknown       TaskState = "unknown"

	RoleUser  Role = "user"
	RoleAgent Role = "agent"
)

// Enum type.
type (
	TaskState string
	Role      string
)

// Task.
type (
	Task struct {
		// Unique identifier for the task.
		ID string `json:"id"`

		// Client-generated id for the session holding the task.
		SessionID string `json:"session_id"`

		// Current status of the task.
		Status string `json:"status"`

		// History of messages exchanged between the agent and the client.
		History []Message `json:"history,omitempty"`

		// Collection of artifacts created by the agent.
		Artifacts []Artifact `json:"artifacts,omitempty"`

		// Extension metadata.
		Metadata map[string]any `json:"metadata,omitempty"`
	}

	// TaskState and accompanying message.
	TaskStatus struct {

		// See [TaskState]
		State TaskState `json:"state"`

		// Additional status updates for client.
		Message *Message `json:"message,omitempty"`

		// ISO datetime value.
		Timestamp *string `json:"timestamp,omitempty"`
	}

	// sent by server during sendSubscribe or subscribe requests.
	TaskStatusUpdateEvent struct {
		// Unique identifier for the task.
		ID string `json:"id"`

		Status TaskStatus `json:"status"`

		// Indicates the end of the event stream.
		Final bool `json:"final"`

		// Extension metadata.
		Metadata map[string]any `json:"metadata,omitempty"`
	}

	// sent by server during sendSubscribe or subscribe requests
	TaskArtifactUpdateEvent struct {
		// Unique identifier for the task.
		ID string `json:"id"`

		Artifact Artifact `json:"artifact"`

		// Extension metadata.
		Metadata map[string]any `json:"metadata,omitempty"`
	}

	// Sent by the client to the agent to create, continue, or restart a task.
	TaskSendParams struct {
		// Unique identifier for the task.
		ID string `json:"id"`

		// Server creates a new sessionId for new tasks if not set.
		SessionID *string `json:"session_id,omitempty"`

		// Message to send to the agent.
		Message Message `json:"message"`

		// Number of recent messages to be retrieved.
		HistoryLength *int `json:"history_length,omitempty"`

		// Where the server should send notifications when disconnected.
		PushNotification *PushNotificationConfig `json:"push_notification,omitempty"`

		// Extension metadata.
		Metadata map[string]any `json:"metadata,omitempty"`
	}
)

// Artifacts are generated as an end result of a Task.
// Artifacts are immutable, can be named, and can have multiple parts.
// A streaming response can append parts to existing Artifacts.
//
// A single Task can generate many Artifacts.
// For example, "create a webpage" could create separate HTML and image Artifacts.
type Artifact struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Parts       []Part  `json:"parts"`

	// Extension metadata.
	Metadata map[string]any `json:"metadata,omitempty"`

	Index     int   `json:"index"`
	Append    *bool `json:"append,omitempty"`
	LastChunk *bool `json:"last_chunk,omitempty"`
}

// Message contains any content that is not an Artifact.
// This can include things like agent thoughts, user context, instructions, errors, status, or metadata.
//
// All content from a client comes in the form of a Message.
// Agents send Messages to communicate status or to provide instructions (whereas generated results are sent as Artifacts).
//
// A Message can have multiple parts to denote different pieces of content.
// For example, a user request could include a textual description from a user and then multiple files used as context from the client.
type Message struct {
	Role Role `json:"role"`

	Parts []Part `json:"parts"`

	// Extension metadata.
	Metadata map[string]any `json:"metadata,omitempty"`
}

// Parts.
type (

	// A fully formed piece of content exchanged between a client and a remote agent as part of a Message or an Artifact.
	// Each Part has its own content type and metadata.
	Part struct {
		Type string `json:"type"`

		// Extension metadata.
		Metadata map[string]any `json:"metadata,omitempty"`
	}

	TextPart struct {
		Part `json:",inline"`
		Text string `json:"text"`
	}

	FilePart struct {
		Part `json:",inline"`
		File struct {
			Name     *string `json:"name,omitempty"`
			MimeType *string `json:"mime_type,omitempty"`

			// Base64 encoded content.
			Bytes *string `json:"bytes,omitempty"`
			Uri   *string `json:"uri,omitempty"`
		} `json:"file"`
	}

	DataPart struct {
		Part `json:",inline"`
		Data map[string]any `json:"data,omitempty"`
	}
)

type (
	PushNotificationConfig struct {
		Url            string          `json:"url"`
		Token          *string         `json:"token,omitempty"`
		Authentication *Authentication `json:"authentication,omitempty"`
	}

	TaskPushNotificationConfig struct {
		ID                     string                 `json:"id"`
		PushNotificationConfig PushNotificationConfig `json:"push_notification"`
	}
)
