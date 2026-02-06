package http

import handlerspolicies "github.com/adhikag24/policy-based-permission-model/http/handlers/policies"

type Handlers struct {
	Policies *handlerspolicies.Handler
}
