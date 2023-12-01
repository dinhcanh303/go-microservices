package constant

const (
	//HEADER
	ApiKey        = "x-api-key"
	ClientID      = "x-client-id"
	Authorization = "authorization"
	RefreshToken  = "refresh-token"

	//UUID
	NullUUID = "00000000-0000-0000-0000-000000000000"
	//Context
	User = "X-Forwarded-User"
)

var HeaderMap = []string{
	ApiKey,
	ClientID,
	Authorization,
	RefreshToken,
}
