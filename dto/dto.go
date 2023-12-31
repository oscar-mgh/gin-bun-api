package dto

type GenreDto struct {
	Name string `json:"name"`
}
type MovieDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Year        int    `json:"year"`
	GenreID     int    `json:"genre_id"`
}
type UserDto struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	ProfileID int64  `json:"profile_id"`
}
type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginRespuestaDto struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}
