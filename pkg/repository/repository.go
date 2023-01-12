package repository

type TokenStorage interface {
	Get(chatId int, bucket Bucket) (string, error)
	Save(chatId int, token string, bucket Bucket) error
}

type Bucket string

const (
	AccessTokens  Bucket = "access_tokens"
	RequestTokens Bucket = "request_tokens"
)
