package status

type Status struct {
	Message string  `json:"error"`
	Code    int     `json:"code"`
	Details Details `json:"details"`
}

type Details map[DetailKey]interface{}

type DetailKey string

var (
	Page        DetailKey = "Page"
	TableIndex  DetailKey = "TableIndex"
	RowIndex    DetailKey = "RowIndex"
	ColumnIndex DetailKey = "ColumnIndex"
	KeysLength  DetailKey = "KeysLength"
	RowLength   DetailKey = "RowLength"
)

func (e Status) Error() string {
	return e.Message
}

type Option func(*Status)

func WithDetails(d Details) Option {
	return func(e *Status) {
		e.Details = d
	}
}

func NewStatus(msg string, code int, options ...Option) Status {
	e := Status{
		Message: msg,
		Code:    code,
	}

	for _, o := range options {
		o(&e)
	}
	return e
}
