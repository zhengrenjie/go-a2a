package protocol

import "context"

const (
	MethodSendTask                 A2AMethod = "tasks/send"
	MethodGetTask                  A2AMethod = "tasks/get"
	MethodCancelTask               A2AMethod = "tasks/cancel"
	MethodSetTaskPushNotifications A2AMethod = "tasks/pushNotification/set"
	MethodGetTaskPushNotifications A2AMethod = "tasks/pushNotification/get"
	MethodSubscribeTask            A2AMethod = "tasks/sendSubscribe"
	MethodResubscribeTask          A2AMethod = "tasks/resubscribe"
)

type A2AMethod string

// IA2AProtocol defines the interface for the A2A protocol.
// Either the client or the server must implement this interface.
type IA2AProtocol interface {

	// AgentCard returns the agent card of the protocol.
	// Remote Agents that support A2A are required to publish an Agent Card in JSON format describing the agent’s capabilities/skills and authentication mechanism.
	// Clients use the Agent Card information to identify the best agent that can perform a task and leverage A2A to communicate with that remote agent.
	AgentCard() AgentCard

	SendTask(ctx context.Context, params *TaskSendParams) (*Task, error)

	GetTask(ctx context.Context, params *TaskSendParams) (*Task, error)

	CancelTask(ctx context.Context, params *TaskSendParams) (*Task, error)

	SetTaskPushNotifications(ctx context.Context, params *TaskPushNotificationConfig) (*TaskPushNotificationConfig, error)

	GetTaskPushNotifications(ctx context.Context, params *TaskPushNotificationConfig) (*TaskPushNotificationConfig, error)
}
