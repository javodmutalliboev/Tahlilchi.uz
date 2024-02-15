package model

type BusinessPromotionalPost struct {
	ID                  int                            `json:"id"`
	TitleLatin          string                         `json:"title_latin"`
	DescriptionLatin    string                         `json:"description_latin"`
	TitleCyrillic       string                         `json:"title_cyrillic"`
	DescriptionCyrillic string                         `json:"description_cyrillic"`
	Photos              []BusinessPromotionalPostPhoto `json:"photos"`
	Videos              []string                       `json:"videos"`
	CoverImage          []byte                         `json:"cover_image"`
	Expiration          string                         `json:"expiration"`
	CreatedAt           string                         `json:"created_at"`
	UpdatedAt           string                         `json:"updated_at"`
	Archived            bool                           `json:"archived"`
	Partner             string                         `json:"partner"`
	Completed           bool                           `json:"completed"`
}
