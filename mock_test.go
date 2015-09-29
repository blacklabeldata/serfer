package serfer

import (
	"github.com/hashicorp/serf/serf"

	"github.com/stretchr/testify/mock"
)

// MockEventHandler mocks a basic Event handler.
type MockEventHandler struct {
	mock.Mock
}

// HandleEvent processes member events.
func (m *MockEventHandler) HandleEvent(e serf.Event) {
	m.Called(e)
	return
}

// MockMemberEventHandler mocks MemberEvent handlers.
type MockMemberEventHandler struct {
	mock.Mock
}

// HandleMemberEvent processes member events.
func (m *MockMemberEventHandler) HandleMemberEvent(e serf.MemberEvent) {
	m.Called(e)
	return
}

// MockUserEventHandler mocks UserEvent handlers.
type MockUserEventHandler struct {
	mock.Mock
}

// HandleUserEvent processes UserEvents.
func (m *MockUserEventHandler) HandleUserEvent(e serf.UserEvent) {
	m.Called(e)
	return
}

// MockReconciler mocks a Reconciler.
type MockReconciler struct {
	mock.Mock
}

// Reconcile processes Members.
func (m *MockReconciler) Reconcile(e serf.Member) {
	m.Called(e)
	return
}

// MockQueryEventHandler mocks QueryEvent handlers.
type MockQueryEventHandler struct {
	mock.Mock
}

// HandleQueryEvent processes QueryEvents.
func (m *MockQueryEventHandler) HandleQueryEvent(e serf.Query) {
	m.Called(e)
	return
}

// MockEvent
type MockEvent struct {
	mock.Mock

	Type serf.EventType
	Name string
}

// EventType returns the EventType
func (m *MockEvent) EventType() serf.EventType {
	m.Called()
	return m.Type
}

// String returns the EventType name
func (m *MockEvent) String() string {
	return m.Name
}
