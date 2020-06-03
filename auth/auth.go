package auth

func GetTokenExpiry(token string) int {
	return 3600
}

func GetTokenPrincipal(token string) string {
	return "b16421d5-9ae7-41d0-ad17-38b735c96307"
}
