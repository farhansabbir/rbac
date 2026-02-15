package lib

import "fmt"

type Gatekeeper struct {
	requestsRejected uint64
	requestsAccepted uint64
}

func NewGatekeeper() *Gatekeeper {
	return &Gatekeeper{
		requestsRejected: 0,
		requestsAccepted: 0,
	}
}

func (g *Gatekeeper) incrementRequestsRejected() {
	g.requestsRejected++
}

func (g *Gatekeeper) incrementRequestsAccepted() {
	g.requestsAccepted++
}

func (g *Gatekeeper) IsRequestAllowed(requestcontext *RequestContext) (bool, error) {
	if requestcontext.RequestResourceType == ResourceTypeNone {
		g.incrementRequestsRejected()
		return false, fmt.Errorf("RequestResourceType cannot be ResourceTypeNone")
	}
	rules_for_request := requestcontext.PrincipalProfile.profRuleMap[uint32(requestcontext.RequestResourceType)]
	fmt.Println(rules_for_request)
	g.incrementRequestsAccepted()
	return true, nil
}
