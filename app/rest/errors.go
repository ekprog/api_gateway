package rest

import "microservice/app/core"

func ServerError() core.StatusResponse {
	return core.StatusResponse{
		Status: core.Status{
			Code:    core.ServerError,
			Message: "We will fix it soon....",
		},
	}
}

func ValidationError(msg ...string) core.StatusResponse {

	msg_ := core.ValidationError
	if len(msg) > 0 {
		msg_ = msg[0]
	}

	return core.StatusResponse{
		Status: core.Status{
			Code:    core.ValidationError,
			Message: msg_,
		},
	}
}

func NotFoundError() core.StatusResponse {
	return core.StatusResponse{
		Status: core.Status{
			Code: core.NotFound,
		},
	}
}

func UnauthorizedError() core.StatusResponse {
	return core.StatusResponse{
		Status: core.Status{
			Code: core.Unauthorised,
		},
	}
}
