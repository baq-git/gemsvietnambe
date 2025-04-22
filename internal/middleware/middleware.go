package middleware

type contextK string

const (
	UserIDContextK     contextK = "userID"
	ContextValuesK     contextK = "contextValues"
	SecretkeyContextK  contextK = "secretkey"
	RefreshkeyContextK contextK = "refreshkey"
)
