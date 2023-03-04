package domain

type VerifyRequest struct {
	AccessToken string `json:"access_token"`
}

type VerifyUser struct {
	Id string `json:"id"`
}

type VerifyResponse struct {
	Status Status
	User   *VerifyUser
}
