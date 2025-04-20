package protocol

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockProtocolUnimplemented(t *testing.T) {
	var (
		ctx  = context.Background()
		mock = NewMockA2AProtocol()
	)

	assert.Panics(t, func() {
		mock.AgentCard()
	})

	assert.Panics(t, func() {
		mock.SendTask(ctx, &TaskSendParams{})
	})

	assert.Panics(t, func() {
		mock.GetTask(ctx, &TaskSendParams{})
	})

	assert.Panics(t, func() {
		mock.CancelTask(ctx, &TaskSendParams{})
	})

	assert.Panics(t, func() {
		mock.SetTaskPushNotifications(ctx, &TaskPushNotificationConfig{})
	})

	assert.Panics(t, func() {
		mock.GetTaskPushNotifications(ctx, &TaskPushNotificationConfig{})
	})

	assert.Panics(t, func() {
		mock.SubscribeTask(ctx, &TaskSendParams{})
	})

	assert.Panics(t, func() {
		mock.ResubscribeTask(ctx, &TaskSendParams{})
	})
}
