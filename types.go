package paxos

// Value represents any value that can be agreed upon
type Value interface{}

// ProposalID represents a unique identifier for a proposal
type ProposalID struct {
    Number int64
    NodeID string
}

// Message types for Paxos protocol
type PrepareRequest struct {
    ProposalID ProposalID
}

type PrepareResponse struct {
    ProposalID      ProposalID
    AcceptedID      *ProposalID
    AcceptedValue   Value
    Ok              bool
}

type AcceptRequest struct {
    ProposalID ProposalID
    Value     Value
}

type AcceptResponse struct {
    ProposalID ProposalID
    Ok        bool
}

// Proposer interface defines the behavior of a proposer in the Paxos protocol
type Proposer interface {
    Propose(value Value) error
    HandlePrepareResponse(response PrepareResponse)
    HandleAcceptResponse(response AcceptResponse)
}

// Acceptor interface defines the behavior of an acceptor in the Paxos protocol
type Acceptor interface {
    HandlePrepare(request PrepareRequest) PrepareResponse
    HandleAccept(request AcceptRequest) AcceptResponse
}

// Learner interface defines the behavior of a learner in the Paxos protocol
type Learner interface {
    Learn(proposalID ProposalID, value Value)
    GetLearnedValue() Value
}