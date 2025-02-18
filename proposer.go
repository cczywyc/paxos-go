package paxos

import (
    "errors"
    "sync"
)

// Proposer represents a node that can make proposals in the Paxos protocol
type proposer struct {
    nodeID          string
    proposalNumber  int64
    currentID       *ProposalID
    value           Value
    prepareCount    int
    acceptCount     int
    acceptors       []Acceptor
    mu              sync.Mutex
    quorumSize      int
}

// NewProposer creates a new Proposer instance
func NewProposer(nodeID string, acceptors []Acceptor) Proposer {
    return &proposer{
        nodeID:     nodeID,
        acceptors:  acceptors,
        quorumSize: (len(acceptors) / 2) + 1,
    }
}

// Propose initiates the Paxos protocol with a value
func (p *proposer) Propose(value Value) error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.currentID != nil {
        return errors.New("proposal already in progress")
    }

    p.proposalNumber++
    p.currentID = &ProposalID{Number: p.proposalNumber, NodeID: p.nodeID}
    p.value = value
    p.prepareCount = 0
    p.acceptCount = 0

    // Phase 1: Prepare
    prepareReq := PrepareRequest{ProposalID: *p.currentID}
    for _, acceptor := range p.acceptors {
        go func(a Acceptor) {
            response := a.HandlePrepare(prepareReq)
            p.HandlePrepareResponse(response)
        }(acceptor)
    }

    return nil
}

// HandlePrepareResponse processes responses from acceptors during the prepare phase
func (p *proposer) HandlePrepareResponse(response PrepareResponse) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.currentID == nil || response.ProposalID != *p.currentID {
        return // Ignore old responses
    }

    if response.Ok {
        p.prepareCount++
        if response.AcceptedValue != nil && response.AcceptedID != nil {
            // If there's a previously accepted value, we must propose that instead
            p.value = response.AcceptedValue
        }

        // If we have a quorum, begin Phase 2
        if p.prepareCount >= p.quorumSize {
            acceptReq := AcceptRequest{
                ProposalID: *p.currentID,
                Value:     p.value,
            }
            for _, acceptor := range p.acceptors {
                go func(a Acceptor) {
                    response := a.HandleAccept(acceptReq)
                    p.HandleAcceptResponse(response)
                }(acceptor)
            }
        }
    }
}

// HandleAcceptResponse processes responses from acceptors during the accept phase
func (p *proposer) HandleAcceptResponse(response AcceptResponse) {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.currentID == nil || response.ProposalID != *p.currentID {
        return // Ignore old responses
    }

    if response.Ok {
        p.acceptCount++
        if p.acceptCount >= p.quorumSize {
            // Consensus achieved!
            p.currentID = nil // Reset for next proposal
        }
    }
}