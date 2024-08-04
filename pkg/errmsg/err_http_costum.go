package errmsg

type HttpError struct {
	Code   int
	Errors map[string][]string
	Msg    string
}

func (e *HttpError) Error() string {
	return e.Msg
}

func NewHttpErrors(errCode int, opts ...Option) *HttpError {
	err := &HttpError{
		Code:   errCode,
		Errors: make(map[string][]string),
		Msg:    "Permintaan Anda gagal diproses",
	}

	for _, opt := range opts {
		opt(err)
	}

	return err
}

func (e *HttpError) Add(field, msg string) {
	e.Errors[field] = append(e.Errors[field], msg)
}

func (e *HttpError) HasErrors() bool {
	return len(e.Errors) > 0
}

type Option func(*HttpError)

func WithMessage(msg string) Option {
	return func(err *HttpError) {
		err.Msg = msg
	}
}

func WithErrors(field string, msg string) Option {
	return func(err *HttpError) {
		err.Errors[field] = append(err.Errors[field], msg)
	}
}

func errorCustomHandler(err *HttpError) (int, *HttpError) {
	return err.Code, err
}
