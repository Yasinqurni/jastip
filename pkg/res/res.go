// Package res will wrap response with the standard JSON structure
package res

// BaseResponse is struct contains base response field including
// Success bool is response status is it success or not
// Message is the response message
type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// SuccessResponse is the struct that contains BaseResponse and Data
type SuccessResponse struct {
	*BaseResponse
	Data  interface{} `json:"data"`
	Meta  interface{} `json:"meta,omitempty"`
	Links interface{} `json:"links,omitempty"`
}

// FailedResponse is the struct that contains BaseResponse and Error
type FailedResponse struct {
	*BaseResponse
	Error interface{} `json:"error"`
}

// JSON function will wrap the response struct based on the success bool value
// Return SuccessResponse if its true and FailedResponse if its false
// Contains success bool, message string and dataOrError interface{} as required parameter
func JSON(success bool, message string, dataOrError interface{}, args ...interface{}) interface{} {
	baseResponse := BaseResponse{
		Success: success,
		Message: message,
	}
	if success {
		// Set links and meta if exists
		var links interface{}
		var meta interface{}
		if len(args) > 0 {
			links = args[0]
		}
		if len(args) > 1 {
			meta = args[1]
		}

		// Return success response
		return SuccessResponse{
			BaseResponse: &baseResponse,
			Data:         dataOrError,
			Links:        links,
			Meta:         meta,
		}
	}

	return FailedResponse{
		BaseResponse: &baseResponse,
		Error:        dataOrError,
	}
}
