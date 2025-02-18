package paxos

import (
    "sync"
)

// Learner represents a node that learns the agreed-upon value in the Paxos protocol
type learner struct {
    mu            sync.RWMutex
    learnedValue  Value
    learnedID     *ProposalID
}

// NewLearner creates a new Learner instance
func NewLearner() Learner {
    return &learner{}
}

// Learn records a value that has been accepted by a quorum of acceptors
func (l *learner) Learn(proposalID ProposalID, value Value) {
    l.mu.Lock()
    defer l.mu.Unlock()

    // Only learn if this is a newer proposal than what we've already learned
    if l.learnedID == nil || proposalID.Number > l.learnedID.Number {
        l.learnedValue = value
        l.learnedID = &proposalID
    }
}

// GetLearnedValue returns the most recently learned value
func (l *learner) GetLearnedValue() Value {
    l.mu.RLock()
    defer l.mu.RUnlock()
    return l.learnedValue
}