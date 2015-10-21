package common

type InternalServerError string

func (e InternalServerError) Error() string {
	return "internal server error: " + string(e)
}

type AuthenticationError string

func (e AuthenticationError) Error() string {
	return "authentication error: " + string(e)
}
