package protocol

import "context"

type A2AProtocol interface {
	AgentCard() AgentCard

	SendTask(ctx context.Context, params *TaskSendParams) (*Task, error)

	GetTask(ctx context.Context, params *TaskSendParams) (*Task, error)

	CancelTask(ctx context.Context, params *TaskSendParams) (*Task, error)

	SetTaskPushNotifications(ctx context.Context, params *TaskPushNotificationConfig) (*TaskPushNotificationConfig, error)

	GetTaskPushNotifications(ctx context.Context, params *TaskPushNotificationConfig) (*TaskPushNotificationConfig, error)
}
