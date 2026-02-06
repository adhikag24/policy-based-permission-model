package http

import (
	handlersfunnels "github.com/adhikag24/policy-based-permission-model/http/handlers/funnels"
	handlerspolicies "github.com/adhikag24/policy-based-permission-model/http/handlers/policies"
)

type Handlers struct {
	Policies *handlerspolicies.Handler
	Funnels  *handlersfunnels.Handler
}
