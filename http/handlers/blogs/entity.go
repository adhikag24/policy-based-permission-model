package blogs

import "github.com/adhikag24/policy-based-permission-model/http/handlers/shared"

type (
	CommonRequest[T any] shared.CommonRequest[T]
	Response[T any]      shared.Response[T]
)
