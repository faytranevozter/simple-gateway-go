package response

// Error Simple error response
func Error(code int, message string) Base {
	return Base{
		Status:  code,
		Message: message,
		Data:    make(map[string]interface{}),
	}
}

// Success Success response
func Success(i interface{}) Base {
	return Base{
		Status:  200,
		Message: "success",
		Data:    i,
	}
}
