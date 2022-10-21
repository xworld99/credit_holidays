package models

type InternalError struct {
	Err  error
	Type int
}

func (ie *InternalError) Error() string {
	return ie.Err.Error()
}
