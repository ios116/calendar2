package exceptions

type DomainError string

func (e DomainError) Error() string {
	return string(e)
}

var (
	TitleRequired       = DomainError("Title or login is not valid")
	DateRequired        = DomainError("Date required")
	IDRequired        = DomainError("Id required")
	DurationRequired    = DomainError("Duration required")
	AuthorRequired    = DomainError("Author required")

	ObjectDoesNotExist  = DomainError("Object does not exist")
	InternalServerError = DomainError("Internal Server Error")
)
