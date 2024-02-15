package model

type BusinessPromotionalPostPhoto struct {
	ID        int    `json:"id"`
	BPP       int    `json:"bpp"`
	FileName  string `json:"file_name"`
	File      []byte `json:"file"`
	CreatedAt string `json:"created_at"`
}
