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

type IA2AProtocol interface {
	AgentCard() AgentCard

	SendTask(ctx context.Context, params *TaskSendParams) (*Task, error)

	GetTask(ctx context.Context, params *TaskSendParams) (*Task, error)

	CancelTask(ctx context.Context, params *TaskSendParams) (*Task, error)

	SetTaskPushNotifications(ctx context.Context, params *TaskPushNotificationConfig) (*TaskPushNotificationConfig, error)

	GetTaskPushNotifications(ctx context.Context, params *TaskPushNotificationConfig) (*TaskPushNotificationConfig, error)
}
