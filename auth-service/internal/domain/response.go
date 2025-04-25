package domain

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}

type LoginResponse struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccessTokenExpiresIn  int    `json:"accessTokenExpiresIn"`
	RefreshTokenExpiresIn int    `json:"refreshTokenExpiresIn"`
}

type RefreshTokenResponse struct {
	AccessToken           string `json:"accessToken"`
	RefreshToken          string `json:"refreshToken"`
	AccessTokenExpiresIn  int    `json:"accessTokenExpiresIn"`
	RefreshTokenExpiresIn int    `json:"refreshTokenExpiresIn"`
}
