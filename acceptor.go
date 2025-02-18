package paxos

// Acceptor represents a node that can accept proposals in the Paxos protocol
type acceptor struct {
    promisedID    *ProposalID
    acceptedID    *ProposalID
    acceptedValue Value
}

// NewAcceptor creates a new Acceptor instance
func NewAcceptor() Acceptor {
    return &acceptor{}
}

// HandlePrepare processes a prepare request from a proposer
func (a *acceptor) HandlePrepare(request PrepareRequest) PrepareResponse {
    // If we haven't promised anyone yet, or this proposal is higher than our promise
    if a.promisedID == nil || (request.ProposalID.Number > a.promisedID.Number) {
        a.promisedID = &request.ProposalID
        return PrepareResponse{
            ProposalID:    request.ProposalID,
            AcceptedID:    a.acceptedID,
            AcceptedValue: a.acceptedValue,
            Ok:           true,
        }
    }

    // Reject if the proposal number is less than what we've promised
    return PrepareResponse{
        ProposalID: request.ProposalID,
        Ok:         false,
    }
}

// HandleAccept processes an accept request from a proposer
func (a *acceptor) HandleAccept(request AcceptRequest) AcceptResponse {
    // Accept only if the proposal ID is greater than or equal to what we've promised
    if a.promisedID == nil || request.ProposalID.Number >= a.promisedID.Number {
        a.promisedID = &request.ProposalID
        a.acceptedID = &request.ProposalID
        a.acceptedValue = request.Value
        return AcceptResponse{
            ProposalID: request.ProposalID,
            Ok:         true,
        }
    }

    // Reject if the proposal number is less than what we've promised
    return AcceptResponse{
        ProposalID: request.ProposalID,
        Ok:         false,
    }
}