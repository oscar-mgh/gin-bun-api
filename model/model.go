package model

import "github.com/uptrace/bun"

type GenreModel struct {
	bun.BaseModel `bun:"table:genres"`

	ID   int64  `bun:",pk,autoincrement" json:"id"`
	Name string `bun:"name,notnull" json:"name"`
	Slug string `bun:"slug,notnull" json:"slug"`
}

type MovieModel struct {
	bun.BaseModel `bun:"table:movies"`

	ID          int64      `bun:",pk,autoincrement" json:"id"`
	Name        string     `bun:"name,notnull" json:"name"`
	Slug        string     `bun:"slug,notnull" json:"slug"`
	Description string     `bun:"description" json:"description"`
	Year        int        `bun:"year" json:"year"`
	GenreID     int64      `bun:"genre_id" json:"genre_id"`
	Genre       GenreModel `bun:"rel:belongs-to,join:genre_id=id" json:"genre"`
}

type MovieImage struct {
	bun.BaseModel `bun:"table:movie_images"`

	ID      int64      `bun:",pk,autoincrement" json:"id"`
	Name    string     `bun:"name,notnull" json:"name"`
	MovieID int64      `bun:"movie_id" json:"movie_id"`
	Movie   MovieModel `bun:"rel:belongs-to,join:movie_id=id" json:"movie"`
}

type ProfileModel struct {
	bun.BaseModel `bun:"table:profiles"`

	ID   int64  `bun:",pk,autoincrement" json:"id"`
	Name string `bun:"name,notnull" json:"name"`
}

type UserModel struct {
	bun.BaseModel `bun:"table:users"`

	ID        int64        `bun:",pk,autoincrement" json:"id"`
	Name      string       `bun:"name,notnull" json:"name"`
	Email     string       `bun:"email,notnull" json:"email"`
	Phone     string       `bun:"phone,notnull" json:"phone"`
	Password  string       `bun:"password,notnull" json:"password"`
	ProfileID int64        `bun:"profile_id" json:"profile_id"`
	Profile   ProfileModel `bun:"rel:belongs-to,join:profile_id=id" json:"profile"`
}
