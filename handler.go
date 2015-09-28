package serfer

import (
	"strings"

	"github.com/hashicorp/serf/serf"
	log "github.com/mgutz/logxi/v1"
)

const (
	// StatusReap is used to update the status of a node if we
	// are handling a EventMemberReap
	StatusReap = serf.MemberStatus(-1)
)

// EventHandler processes generic Serf events. Depending on the
// event type, more processing may be needed.
type EventHandler interface {
	HandleEvent(serf.Event)
}

// MemberEventHandler handles membership change events.
type MemberEventHandler interface {
	HandleMemberEvent(serf.MemberEvent)
}

// UserEventHandler handles user events.
type UserEventHandler interface {
	HandleUserEvent(serf.UserEvent)
}

// QueryEventHandler handles Serf query events.
type QueryEventHandler interface {
	HandleQueryEvent(serf.Query)
}

// Reconciler is used to reconcile Serf events wilth an external process, like Raft.
type Reconciler interface {
	Reconcile(serf.Member)
}

// IsLeaderFunc should return true if the local node is the cluster leader.
type IsLeaderFunc func() bool

// SerfEventHandler is used to dispatch various Serf events to separate event handlers.
type SerfEventHandler struct {

	// ServicePrefix is used to filter out unknown events.
	ServicePrefix string

	// ReconcileOnJoin determines if the Reconiler is called when a node joins the cluster.
	ReconcileOnJoin bool

	// ReconcileOnLeave determines if the Reconiler is called when a node leaves the cluster.
	ReconcileOnLeave bool

	// ReconcileOnFail determines if the Reconiler is called when a node fails.
	ReconcileOnFail bool

	// ReconcileOnUpdate determines if the Reconiler is called when a node updates.
	ReconcileOnUpdate bool

	// ReconcileOnReap determines if the Reconiler is called when a node is reaped from the cluster.
	ReconcileOnReap bool

	// IsLeader determines if the local node is the cluster leader.
	IsLeader IsLeaderFunc

	// IsLeaderEventFunc determines if an event is a leader election event based on the event name.
	IsLeaderEvent func(string) bool

	// LeaderElectionHandler processes leader election events.
	LeaderElectionHandler UserEventHandler

	// UserEvent processes known, non-leader election events.
	UserEvent UserEventHandler

	// UnknownEventHandler processes unkown events.
	UnknownEventHandler UserEventHandler

	// Called when a Member joins the cluster.
	NodeJoined MemberEventHandler

	// Called when a Member leaves the cluster by sending a leave message.
	NodeLeft MemberEventHandler

	// Called when a Member has been detected as failed.
	NodeFailed MemberEventHandler

	// Called when a Member has been Readed from the cluster.
	NodeReaped MemberEventHandler

	// Called when a Member has been updated.
	NodeUpdated MemberEventHandler

	// Called when a membership event occurs.
	Reconciler Reconciler

	// Called when a serf.Query is received.
	QueryHandler QueryEventHandler

	// Logs output
	Logger log.Logger
}

// HandleEvent processes a generic Serf event and dispatches it to the appropriate
// destination.
func (s SerfEventHandler) HandleEvent(e serf.Event) {
	if e == nil {
		return
	}

	var reconcile bool
	switch e.EventType() {

	// If the event is a Join event, call NodeJoined and then reconcile event with
	// persistent storage.
	case serf.EventMemberJoin:
		reconcile = s.ReconcileOnJoin
		if s.NodeJoined != nil {
			s.NodeJoined.HandleMemberEvent(e.(serf.MemberEvent))
		}

	// If the event is a Leave event, call NodeLeft and then reconcile event with
	// persistent storage.
	case serf.EventMemberLeave:
		reconcile = s.ReconcileOnLeave
		if s.NodeLeft != nil {
			s.NodeLeft.HandleMemberEvent(e.(serf.MemberEvent))
		}

	// If the event is a Failed event, call NodeFailed and then reconcile event with
	// persistent storage.
	case serf.EventMemberFailed:
		reconcile = s.ReconcileOnFail
		if s.NodeFailed != nil {
			s.NodeFailed.HandleMemberEvent(e.(serf.MemberEvent))
		}

	// If the event is a Reap event, reconcile event with persistent storage.
	case serf.EventMemberReap:
		reconcile = s.ReconcileOnReap
		if s.NodeReaped != nil {
			s.NodeReaped.HandleMemberEvent(e.(serf.MemberEvent))
		}

	// If the event is a user event, handle leader elections, user events and unknown events.
	case serf.EventUser:
		s.handleUserEvent(e.(serf.UserEvent))

	// If the event is an Update event, call NodeUpdated
	case serf.EventMemberUpdate:
		reconcile = s.ReconcileOnUpdate
		if s.NodeUpdated != nil {
			s.NodeUpdated.HandleMemberEvent(e.(serf.MemberEvent))
		}

	// If the event is a query, call Query Handler
	case serf.EventQuery:
		if s.QueryHandler != nil {
			s.QueryHandler.HandleQueryEvent(*e.(*serf.Query))
		}
	default:
		s.Logger.Warn("unhandled Serf Event: %#v", e)
		return
	}

	// Reconcile event with external storage
	if reconcile && s.Reconciler != nil {
		s.reconcile(e.(serf.MemberEvent))
	}
}

// reconcile is used to reconcile Serf events with the strongly
// consistent store if we are the current leader
func (s *SerfEventHandler) reconcile(me serf.MemberEvent) {

	// Do nothing if we are not the leader.
	if !s.IsLeader() {
		return
	}

	// Check if this is a reap event
	isReap := me.EventType() == serf.EventMemberReap

	// Queue the members for reconciliation
	for _, m := range me.Members {
		// Change the status if this is a reap event
		if isReap {
			m.Status = StatusReap
		}

		// Call reconcile
		if s.Reconciler != nil {
			s.Reconciler.Reconcile(m)
		}
	}
}

// handleUserEvent is called when a user event is received from both local and remote nodes.
func (s *SerfEventHandler) handleUserEvent(event serf.UserEvent) {
	switch name := event.Name; {

	// Handles leader election events
	case s.IsLeaderEvent(name):
		s.Logger.Info("serfer: New leader elected: %s", event.Payload)

		// Process leader election event
		if s.LeaderElectionHandler != nil {
			s.LeaderElectionHandler.HandleUserEvent(event)
		}

	// Handle service events
	case s.isServiceEvent(name):
		event.Name = s.getRawEventName(name)
		s.Logger.Debug("serfer: user event: %s", event.Name)

		// Process user event
		if s.UserEvent != nil {
			s.UserEvent.HandleUserEvent(event)
		}

	// Handle unknown user events
	default:
		s.Logger.Warn("serfer: unknown event: %v", event)

		// Process unknown event
		if s.UnknownEventHandler != nil {
			s.UnknownEventHandler.HandleUserEvent(event)
		}
	}
}

// getRawEventName is used to get the raw event name
func (s *SerfEventHandler) getRawEventName(name string) string {
	return strings.TrimPrefix(name, s.ServicePrefix+":")
}

// isServiceEvent checks if a serf event is a known event
func (s *SerfEventHandler) isServiceEvent(name string) bool {
	return strings.HasPrefix(name, s.ServicePrefix+":")
}
