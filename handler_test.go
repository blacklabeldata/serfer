package serfer

import (
	"net"
	"testing"

	"github.com/hashicorp/serf/serf"
	log "github.com/mgutz/logxi/v1"
	"github.com/stretchr/testify/suite"
)

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestEventHandlerSuite(t *testing.T) {
	suite.Run(t, new(EventHandlerTestSuite))
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type EventHandlerTestSuite struct {
	suite.Suite
	Handler SerfEventHandler
	Member  serf.Member
}

// Make sure that Handler and Member are set before each test
func (suite *EventHandlerTestSuite) SetupTest() {
	suite.Handler = SerfEventHandler{
		Logger: &log.NullLogger{},
	}

	suite.Member = serf.Member{
		Name:        "",
		Addr:        net.ParseIP("127.0.0.1"),
		Port:        9022,
		Tags:        make(map[string]string),
		Status:      serf.StatusAlive,
		ProtocolMin: serf.ProtocolVersionMin,
		ProtocolMax: serf.ProtocolVersionMax,
		ProtocolCur: serf.ProtocolVersionMax,
		DelegateMin: serf.ProtocolVersionMin,
		DelegateMax: serf.ProtocolVersionMax,
		DelegateCur: serf.ProtocolVersionMax,
	}
}

// Test NodeJoin events are processed properly
func (suite *EventHandlerTestSuite) TestNodeJoined() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberJoin,
		[]serf.Member{suite.Member},
	}

	// Add NodeJoined handler
	m := &MockMemberEventHandler{}
	m.On("HandleMemberEvent", evt).Return()
	suite.Handler.NodeJoined = m

	// Add Reconciler
	r := &MockReconciler{}
	r.On("Reconcile", evt).Return()
	suite.Handler.Reconciler = r

	// Process event
	suite.Handler.HandleEvent(evt)
	m.AssertCalled(suite.T(), "HandleMemberEvent", evt)
	r.AssertCalled(suite.T(), "Reconcile", evt)
}

// Test NodeLeave messages are dispatched properly
func (suite *EventHandlerTestSuite) TestNodeLeave() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberLeave,
		[]serf.Member{suite.Member},
	}

	// Add NodeLeft handler
	m := &MockMemberEventHandler{}
	m.On("HandleMemberEvent", evt).Return()
	suite.Handler.NodeLeft = m

	// Add Reconciler
	r := &MockReconciler{}
	r.On("Reconcile", evt).Return()
	suite.Handler.Reconciler = r

	// Process event
	suite.Handler.HandleEvent(evt)
	m.AssertCalled(suite.T(), "HandleMemberEvent", evt)
	r.AssertCalled(suite.T(), "Reconcile", evt)
}

// Test NodeFailed messages are dispatched properly
func (suite *EventHandlerTestSuite) TestNodeFailed() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberFailed,
		[]serf.Member{suite.Member},
	}

	// Add NodeFailed handler
	m := &MockMemberEventHandler{}
	m.On("HandleMemberEvent", evt).Return()
	suite.Handler.NodeFailed = m

	// Add Reconciler
	r := &MockReconciler{}
	r.On("Reconcile", evt).Return()
	suite.Handler.Reconciler = r

	// Process event
	suite.Handler.HandleEvent(evt)
	m.AssertCalled(suite.T(), "HandleMemberEvent", evt)
	r.AssertCalled(suite.T(), "Reconcile", evt)
}

// Test NodeReaped messages are dispatched properly
func (suite *EventHandlerTestSuite) TestNodeReaped() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberReap,
		[]serf.Member{suite.Member},
	}

	// Add NodeReaped handler
	m := &MockMemberEventHandler{}
	m.On("HandleMemberEvent", evt).Return()
	suite.Handler.NodeReaped = m

	// Add Reconciler
	r := &MockReconciler{}
	r.On("Reconcile", evt).Return()
	suite.Handler.Reconciler = r

	// Process event
	suite.Handler.HandleEvent(evt)
	m.AssertCalled(suite.T(), "HandleMemberEvent", evt)
	r.AssertCalled(suite.T(), "Reconcile", evt)
}

// Test UserEvent messages are dispatched properly
func (suite *EventHandlerTestSuite) TestUserEvent() {

	// Create Member Event
	evt := serf.UserEvent{
		LTime:    serf.LamportTime(0),
		Name:     "Event",
		Payload:  make([]byte, 0),
		Coalesce: false,
	}

	// Add UserEvent handler
	m := &MockUserEventHandler{}
	m.On("HandleUserEvent", evt).Return()
	suite.Handler.UserEvent = m

	// Add Reconciler
	r := &MockReconciler{}
	r.On("Reconcile", evt).Return()
	suite.Handler.Reconciler = r

	// Process event
	suite.Handler.HandleEvent(evt)
	m.AssertCalled(suite.T(), "HandleUserEvent", evt)
	r.AssertNotCalled(suite.T(), "Reconcile", evt)
}

// Test NodeUpdated messages are dispatched properly
func (suite *EventHandlerTestSuite) TestNodeUpdated() {

	// Create Member Event
	evt := serf.MemberEvent{
		serf.EventMemberUpdate,
		[]serf.Member{suite.Member},
	}

	// Add NodeReaped handler
	m := &MockMemberEventHandler{}
	m.On("HandleMemberEvent", evt).Return()
	suite.Handler.NodeUpdated = m

	// Add Reconciler
	r := &MockReconciler{}
	r.On("Reconcile", evt).Return()
	suite.Handler.Reconciler = r

	// Process event
	suite.Handler.HandleEvent(evt)
	m.AssertCalled(suite.T(), "HandleMemberEvent", evt)
	r.AssertCalled(suite.T(), "Reconcile", evt)
}

// Test QueryEvent messages are dispatched properly
func (suite *EventHandlerTestSuite) TestQueryEvent() {

	// Create Query
	query := serf.Query{
		LTime:   serf.LamportTime(0),
		Name:    "Event",
		Payload: make([]byte, 0),
	}

	// Add UserEvent handler
	m := &MockQueryEventHandler{}
	m.On("HandleQueryEvent", query).Return()
	suite.Handler.QueryHandler = m

	// Add Reconciler
	r := &MockReconciler{}
	r.On("Reconcile", query).Return()
	suite.Handler.Reconciler = r

	// Process event
	suite.Handler.HandleEvent(&query)
	m.AssertCalled(suite.T(), "HandleQueryEvent", query)
	r.AssertNotCalled(suite.T(), "Reconcile", query)
}

// Test nil messages are not dispatched properly
func (suite *EventHandlerTestSuite) TestNilEvent() {

	// Add NodeJoined handler
	m1 := &MockMemberEventHandler{}
	suite.Handler.NodeJoined = m1

	// Add NodeLeft handler
	m2 := &MockMemberEventHandler{}
	suite.Handler.NodeLeft = m2

	// Add NodeFailed handler
	m3 := &MockMemberEventHandler{}
	suite.Handler.NodeFailed = m3

	// Add NodeReaped handler
	m4 := &MockMemberEventHandler{}
	suite.Handler.NodeReaped = m4

	// Add NodeUpdated handler
	m5 := &MockMemberEventHandler{}
	suite.Handler.NodeUpdated = m5

	// Add UserEvent handler
	u1 := &MockUserEventHandler{}
	suite.Handler.UserEvent = u1

	// Add UserEvent handler
	q1 := &MockQueryEventHandler{}
	suite.Handler.QueryHandler = q1

	// Add Reconciler
	r1 := &MockReconciler{}
	suite.Handler.Reconciler = r1

	// Process event
	suite.Handler.HandleEvent(nil)
	m1.AssertNotCalled(suite.T(), "HandleMemberEvent")
	m2.AssertNotCalled(suite.T(), "HandleMemberEvent")
	m3.AssertNotCalled(suite.T(), "HandleMemberEvent")
	m4.AssertNotCalled(suite.T(), "HandleMemberEvent")
	m5.AssertNotCalled(suite.T(), "HandleMemberEvent")
	u1.AssertNotCalled(suite.T(), "HandleUserEvent")
	q1.AssertNotCalled(suite.T(), "HandleQueryEvent")
	r1.AssertNotCalled(suite.T(), "Reconcile")
}

// Test unknown messages are not dispatched properly
func (suite *EventHandlerTestSuite) TestUnknownEvent() {

	// Add NodeJoined handler
	m1 := &MockMemberEventHandler{}
	suite.Handler.NodeJoined = m1

	// Add NodeLeft handler
	m2 := &MockMemberEventHandler{}
	suite.Handler.NodeLeft = m2

	// Add NodeFailed handler
	m3 := &MockMemberEventHandler{}
	suite.Handler.NodeFailed = m3

	// Add NodeReaped handler
	m4 := &MockMemberEventHandler{}
	suite.Handler.NodeReaped = m4

	// Add NodeUpdated handler
	m5 := &MockMemberEventHandler{}
	suite.Handler.NodeUpdated = m5

	// Add UserEvent handler
	u1 := &MockUserEventHandler{}
	suite.Handler.UserEvent = u1

	// Add UserEvent handler
	q1 := &MockQueryEventHandler{}
	suite.Handler.QueryHandler = q1

	// Add Reconciler
	r1 := &MockReconciler{}
	suite.Handler.Reconciler = r1

	// Process event
	t1 := &MockEvent{Name: "UnknownType", Type: serf.EventType(-1)}
	t1.On("EventType").Return()
	suite.Handler.HandleEvent(t1)

	// Test Assertions
	t1.AssertCalled(suite.T(), "EventType")
	m1.AssertNotCalled(suite.T(), "HandleMemberEvent")
	m2.AssertNotCalled(suite.T(), "HandleMemberEvent")
	m3.AssertNotCalled(suite.T(), "HandleMemberEvent")
	m4.AssertNotCalled(suite.T(), "HandleMemberEvent")
	m5.AssertNotCalled(suite.T(), "HandleMemberEvent")
	u1.AssertNotCalled(suite.T(), "HandleUserEvent")
	q1.AssertNotCalled(suite.T(), "HandleQueryEvent")
	r1.AssertNotCalled(suite.T(), "Reconcile")
}
