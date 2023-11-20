package types

type User struct {
	Id         string   `json:"id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	Favoris_id []string `json:"favoris_id"`
}

type Commentaire struct {
	Id      string `json:"id"`
	User_id string `json:"user_id"`
	Contenu string `json:"contenu"`
	Film_id string `json:"film_id"`
}

type Film struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type Commentaire_user struct {
	Id      string `json:"id"`
	User_id string `json:"user_id"`
	Contenu string `json:"contenu"`
	Film    Film   `json:"film"`
}

type Commentaire_film struct {
	Id      string   `json:"id"`
	User    User_out `json:"user"`
	Contenu string   `json:"contenu"`
	Film_id string   `json:"film_id"`
}

type User_out struct {
	Id         string   `json:"id"`
	Username   string   `json:"username"`
	Favoris_id []string `json:"favoris_id"`
}

type Film_Fav struct {
	User_id  string `json:"user_id"`
	Film_id  string `json:"film_id"`
	Favorite bool   `json:"favorite"`
}

type Result_fav struct {
	Success  bool     `json:"success"`
	Favorite Film_Fav `json:"favorite"`
}

type Result_film struct {
	Success bool `json:"success"`
	Film    Film `json:"film"`
}

type Result_films struct {
	Success bool   `json:"success"`
	Films   []Film `json:"films"`
}

type Result_com struct {
	Success bool             `json:"success"`
	Comment Commentaire_film `json:"comment"`
}

type Result_user struct {
	Success bool `json:"success"`
	User    User `json:"user"`
}

type Result_user_out struct {
	Success bool     `json:"success"`
	User    User_out `json:"user"`
}
