package domain

type EmailAddress interface {
	String() string
	Equals(other EmailAddress) bool
}
