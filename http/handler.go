package http

import (
	handlersblogs "github.com/adhikag24/policy-based-permission-model/http/handlers/blogs"
	handlersfunnels "github.com/adhikag24/policy-based-permission-model/http/handlers/funnels"
	handlerspolicies "github.com/adhikag24/policy-based-permission-model/http/handlers/policies"
)

type Handlers struct {
	Policies *handlerspolicies.Handler
	Funnels  *handlersfunnels.Handler
	Blogs    *handlersblogs.Handler
}
