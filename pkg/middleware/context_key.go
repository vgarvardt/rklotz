package middleware

const prefix = "rKlotz context key "

type contextKey string

func (c contextKey) String() string {
	return prefix + string(c)
}
