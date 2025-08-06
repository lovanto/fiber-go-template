package repository

const (
	UnauthorizedErrorMessage              string = "unauthorized, check expiration time of your token"
	ForbiddenErrorMessage                 string = "permission denied, check credentials of your token"
	ForbiddenDataModificationErrorMessage string = "permission denied, only the creator can action their data"
	NotFoundErrorMessage                  string = "data not found"
)
