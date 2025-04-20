package protocol

import "context"

type MockA2AProtocol struct {
}

// AgentCard implements IA2AProtocol.
func (m *MockA2AProtocol) AgentCard() AgentCard {
	panic("unimplemented")
}

// CancelTask implements IA2AProtocol.
func (m *MockA2AProtocol) CancelTask(ctx context.Context, params *TaskSendParams) (*Task, error) {
	panic("unimplemented")
}

// GetTask implements IA2AProtocol.
func (m *MockA2AProtocol) GetTask(ctx context.Context, params *TaskSendParams) (*Task, error) {
	panic("unimplemented")
}

// GetTaskPushNotifications implements IA2AProtocol.
func (m *MockA2AProtocol) GetTaskPushNotifications(ctx context.Context, params *TaskPushNotificationConfig) (*TaskPushNotificationConfig, error) {
	panic("unimplemented")
}

// ResubscribeTask implements IA2AProtocol.
func (m *MockA2AProtocol) ResubscribeTask(ctx context.Context, params *TaskSendParams) (chan any, error) {
	panic("unimplemented")
}

// SendTask implements IA2AProtocol.
func (m *MockA2AProtocol) SendTask(ctx context.Context, params *TaskSendParams) (*Task, error) {
	panic("unimplemented")
}

// SetTaskPushNotifications implements IA2AProtocol.
func (m *MockA2AProtocol) SetTaskPushNotifications(ctx context.Context, params *TaskPushNotificationConfig) (*TaskPushNotificationConfig, error) {
	panic("unimplemented")
}

// SubscribeTask implements IA2AProtocol.
func (m *MockA2AProtocol) SubscribeTask(ctx context.Context, params *TaskSendParams) (chan any, error) {
	panic("unimplemented")
}

func NewMockA2AProtocol() *MockA2AProtocol {
	return &MockA2AProtocol{}
}

var _ IA2AProtocol = (*MockA2AProtocol)(nil)
