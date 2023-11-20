package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	api "serveur/api"
	auth "serveur/auth"
	db "serveur/database"
	t "serveur/types"
)

type BadRequest struct {
	Error string `json:"error"`
}

var auth_chan chan auth.AuthQuery

// ----------------------------------------------------------------------------------------------------------------------------------------------------
// -------------------------------------------------------user and account-------------------------------------------------------------------------------
// --------------------------------------------------------------------------------------------------------------------------------------------------------
func handlerUser(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	switch r.Method {
	case "GET":
		getUser(w, r)
	case "PUT":
		register(w, r)
	case "DELETE":
		closeAccount(w, r)
	default:
		w.WriteHeader(400)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: fmt.Sprintf("Méthode non acceptée : %s", r.Method)})
		w.Write(b)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	if r.Method != "POST" {
		fmt.Println("wrong method in logout :", r.Method)
		w.WriteHeader(400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var input struct {
		Auth auth.Auth `json:"auth"`
	}
	err := decoder.Decode(&input)
	if err != nil {
		fmt.Println("error decoding input:", err)
		w.WriteHeader(400)
		return
	}

	fmt.Println("auth :", input.Auth.Session)

	auth := auth.AuthQuery{
		Userid:         input.Auth.Userid,
		Session:        input.Auth.Session,
		Delete_session: true,
		Ret_chan:       make(chan auth.AuthResponse),
	}

	auth_chan <- auth
	auth_res := <-auth.Ret_chan

	fmt.Println("response :", auth_res)

	if auth_res.Ok {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(403)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: "Wrong auth"})
		w.Write(b)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/user/")

	if len(id) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
	}

	fmt.Println("id :", id)

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	result := db.GetUserOut(id, client)

	if result.Success {
		var b []byte
		b, _ = json.Marshal(result.User)
		w.Write(b)
	} else {
		w.WriteHeader(400)
	}
}

func register(w http.ResponseWriter, r *http.Request) { //TODO: verify that the email is different, maybe change the requeste as the users' info
	decoder := json.NewDecoder(r.Body)
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}
	err := decoder.Decode(&input)
	if err != nil {
		fmt.Println("error decoding input:", err)
		w.WriteHeader(400)
		return
	}

	if !strings.HasPrefix(input.Email, "") {
		fmt.Println("error : missing email")
		b, _ := json.Marshal(BadRequest{Error: "missing email"})
		w.Write(b)
		w.WriteHeader(400)
		return
	}

	if !strings.HasPrefix(input.Password, "") {
		fmt.Println("error : missing password")
		b, _ := json.Marshal(BadRequest{Error: "missing password"})
		w.Write(b)
		w.WriteHeader(400)
		return
	}

	if !strings.HasPrefix(input.Username, "") {
		fmt.Println("error : missing username")
		b, _ := json.Marshal(BadRequest{Error: "missing username"})
		w.Write(b)
		w.WriteHeader(400)
		return
	}

	client := db.Connect()
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	userSameEmail := db.GetUserByMail(input.Email, client)
	if userSameEmail.Success {
		fmt.Println("Email already existed:", userSameEmail.User)
		b, _ := json.Marshal(BadRequest{Error: "Email already registered"})
		w.Write(b)
		w.WriteHeader(400)
		return
	}

	var user t.User
	var fav []string
	user.Email = input.Email
	user.Password = input.Password
	user.Username = input.Username
	user.Favoris_id = fav

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	res := db.AddUser(user, client)
	if !res.Success {
		w.WriteHeader(500)
		return
	}

	b, _ := json.Marshal(res.User)
	w.Write(b)
	w.WriteHeader(200)
}

func login(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	if r.Method != "POST" {
		fmt.Println("wrong method in login :", r.Method)
		w.WriteHeader(400)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := decoder.Decode(&input)
	if err != nil {
		fmt.Println("error decoding input:", err)
		w.WriteHeader(400)
		return
	}

	fmt.Println(input)

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	result := db.GetUserByMail(input.Email, client)

	if !result.Success {
		fmt.Println("error finding user:", err)
		w.WriteHeader(401)
		b, _ := json.Marshal(BadRequest{Error: "Email or password invalid"})
		w.Write(b)
		return
	}

	if result.User.Password != input.Password {
		w.WriteHeader(401)
		b, _ := json.Marshal(BadRequest{Error: "Email or password invalid"})
		w.Write(b)
		return
	}

	auth := auth.AuthQuery{
		Userid:      result.User.Id,
		New_session: true,
		Ret_chan:    make(chan auth.AuthResponse),
	}

	auth_chan <- auth
	auth_res := <-auth.Ret_chan

	if auth_res.Ok {
		w.WriteHeader(200)
		b, _ := json.Marshal(auth_res.Auth)
		w.Write(b)
	} else {
		w.WriteHeader(403)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: "Wrong auth"})
		w.Write(b)
	}
}

func closeAccount(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/user/")
	if len(id) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
	}

	fmt.Println("id :", id)

	decoder := json.NewDecoder(r.Body)
	var input struct {
		Auth auth.Auth `json:"auth"`
	}
	err := decoder.Decode(&input)
	if err != nil {
		fmt.Println("error decoding input:", err)
		w.WriteHeader(400)
		return
	}

	auth_q := auth.AuthQuery{
		Userid:   input.Auth.Userid,
		Session:  input.Auth.Session,
		Ret_chan: make(chan auth.AuthResponse),
	}

	auth_chan <- auth_q
	auth_res := <-auth_q.Ret_chan

	if !auth_res.Ok {
		w.WriteHeader(403)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: "Wrong auth"})
		w.Write(b)
		return
	}

	if auth_res.Auth.Userid != id {
		w.WriteHeader(403)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: "userid is not the one in auth"})
		w.Write(b)
		return
	}

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	deleted := db.DeleteUser(auth_res.Auth.Userid, client)
	if !deleted {
		w.WriteHeader(500)
		return
	}

	auth_q = auth.AuthQuery{
		Userid:         input.Auth.Userid,
		Session:        input.Auth.Session,
		Delete_session: true,
		Ret_chan:       make(chan auth.AuthResponse),
	}

	auth_chan <- auth_q
	<-auth_q.Ret_chan

	w.WriteHeader(200)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------
//----------------------------------------------------------films------------------------------------------------------------------------------
//-----------------------------------------------------------------------------------------------------------------------------------------------

// get un film par id
func getFilm(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/film/")
	if r.Method != "GET" {
		fmt.Println("wrong method in getFilm :", r.Method)
		w.WriteHeader(400)
		return
	}

	if len(id) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
	}

	film := api.GetFilm(id)

	if film.Success {
		b, _ := json.Marshal(film.Film)
		w.Write(b)
	} else {
		w.WriteHeader(400)
	}
}

// get un film par genre
func getFilmByGenre(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	genre := strings.TrimPrefix(r.URL.Path, "/search/genre/")
	if r.Method != "GET" {
		fmt.Println("wrong method in getFilm :", r.Method)
		w.WriteHeader(400)
		return
	}

	if len(genre) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
	}

	fmt.Println("genre :", genre)
	films := api.GetFilmsByGenre(genre)

	if len(films) > 0 {
		var filmList []t.Film
		for _, filmId := range films {
			result := api.GetFilm(filmId)
			if result.Success {
				filmList = append(filmList, result.Film)
			}
		}
		jsonResponse, err := json.Marshal(filmList)
		if err != nil {
			fmt.Println("error marshalling json response: ", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		fmt.Println("aucun film trouvé")
		w.WriteHeader(400)
	}
}

// get un film par titre
func searchFilm(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	fmt.Println("recieve request")

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	keyword := strings.TrimPrefix(r.URL.Path, "/search/title/")
	if r.Method != "GET" {
		fmt.Println("wrong method in getFilm :", r.Method)
		w.WriteHeader(400)
		return
	}

	if len(keyword) < 1 {
		fmt.Println("title is missing in parameters")
		w.WriteHeader(400)
		return
	}

	fmt.Println("ready to send to api")

	films := api.SearchFilm(keyword)

	if len(films) > 0 {
		var filmList []t.Film
		for _, filmId := range films {
			result := api.GetFilm(filmId)
			if result.Success {
				filmList = append(filmList, result.Film)
			}
		}

		jsonResponse, err := json.Marshal(filmList)
		if err != nil {
			fmt.Println("error marshalling json response: ", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	} else {
		w.WriteHeader(400)
	}

}

//-------------------------------------------------------------------------------------------------------------------------------------------------
//-----------------------------------------------------------------comments----------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------------

func commentDelete(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var input struct {
		Auth auth.Auth `json:"auth"`
	}
	err := decoder.Decode(&input)
	if err != nil {
		fmt.Println("error decoding input:", err)
		w.WriteHeader(402)
		return
	}

	auth := auth.AuthQuery{
		Userid:   input.Auth.Userid,
		Session:  input.Auth.Session,
		Ret_chan: make(chan auth.AuthResponse),
	}

	fmt.Println("requete : ", auth)

	auth_chan <- auth
	fmt.Println("requete envoyée")
	auth_res := <-auth.Ret_chan

	if !auth_res.Ok {
		w.WriteHeader(403)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: "Wrong auth"})
		w.Write(b)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/comment/")

	if len(id) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
		return
	}

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if db.DeleteComment(id, auth_res.Auth.Userid, client) {
		fmt.Printf("deleted comment id:%s \n", id)
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}

}

func commentsUser(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/comment/user/")
	if r.Method != "GET" {
		fmt.Println("wrong method in getCommentsUser :", r.Method)
		w.WriteHeader(400)
		return
	}

	if len(id) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
	}

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	comments := db.GetCommentsUser(id, client)

	var b []byte
	b, _ = json.Marshal(comments)
	w.Write(b)
}

func commentsFilm(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/comment/film/")
	if r.Method != "GET" {
		fmt.Println("wrong method in getCommentsFilm :", r.Method)
		w.WriteHeader(400)
		return
	}

	if len(id) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
	}

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	comments := db.GetCommentsFilm(id, client)

	var b []byte
	b, _ = json.Marshal(comments)
	w.Write(b)

}

//---------------------------------------------------------------------------------------------------------------------------------------------------
//---------------------------------------------------------------recommends----------------------------------------------------------------------------
//----------------------------------------------------------------------------------------------------------------------------------------------------

func recommendFilm(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/recommend/film/")
	if r.Method != "GET" {
		fmt.Println("wrong method in getRecommendFilm :", r.Method)
		w.WriteHeader(400)
		return
	}

	if len(id) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
	}

	recommend_films := 5

	films := api.GetSimilarFilms(append(make([]string, 0), id), recommend_films)

	if !films.Success {
		w.WriteHeader(500)
		b, _ := json.Marshal(BadRequest{Error: fmt.Sprintf("Error while recommending %d films similar to id %s", recommend_films, id)})
		w.Write(b)
		return
	}

	w.WriteHeader(200)
	b, _ := json.Marshal(films.Films)
	w.Write(b)
}

func recommendUser(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/recommend/user/")
	if r.Method != "GET" {
		fmt.Println("wrong method in recommendUser :", r.Method)
		w.WriteHeader(400)
		return
	}

	if len(id) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
	}

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	user := db.GetUserOut(id, client)
	if !user.Success {
		w.WriteHeader(500)
		b, _ := json.Marshal(BadRequest{Error: fmt.Sprintf("Error while retrieving user %s", id)})
		w.Write(b)
		return
	}

	recommend_films := 10

	films := api.GetSimilarFilms(user.User.Favoris_id, recommend_films)

	if !films.Success {
		w.WriteHeader(500)
		b, _ := json.Marshal(BadRequest{Error: fmt.Sprintf("Error while recommending %d films of user %s", recommend_films, id)})
		w.Write(b)
		return
	}

	if len(films.Films) < recommend_films {
		fmt.Println("populating with popular films", recommend_films-len(films.Films))
		fs := api.GetPopularFilms(recommend_films - len(films.Films))
		if fs.Success {
			films.Films = append(films.Films, fs.Films...)
		} else {
			fmt.Println("error")
		}
	}

	w.WriteHeader(200)
	b, _ := json.Marshal(films.Films)
	w.Write(b)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------
//------------------------------------------------------------------------fav---------------------------------------------------------------------------
//------------------------------------------------------------------------------------------------------------------------------------------------------

func favorites(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/favorites/user/")
	if r.Method != "GET" {
		fmt.Println("wrong method in favorites :", r.Method)
		w.WriteHeader(400)
		return
	}

	if len(id) < 1 {
		fmt.Println("id is missing in parameters")
		w.WriteHeader(400)
	}

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	favs := db.GetUserFavoriteFilms(id, client)

	b, _ := json.Marshal(favs)
	w.Write(b)
}

func handleFavorite(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	switch r.Method {
	case "GET":
		isFavorite(w, r)
	case "POST":
		setFavorite(w, r)
	default:
		w.WriteHeader(400)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: fmt.Sprintf("Méthode non acceptée : %s", r.Method)})
		w.Write(b)
	}
}

func setFavorite(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var input struct {
		Auth     auth.Auth `json:"auth"`
		User_id  string    `json:"user_id"`
		Film_id  string    `json:"film_id"`
		Favorite bool      `json:"favorite"`
	}
	err := decoder.Decode(&input)
	if err != nil {
		fmt.Println("error decoding input:", err)
		w.WriteHeader(400)
		return
	}

	auth := auth.AuthQuery{
		Userid:   input.Auth.Userid,
		Session:  input.Auth.Session,
		Ret_chan: make(chan auth.AuthResponse),
	}

	auth_chan <- auth
	auth_res := <-auth.Ret_chan

	if !auth_res.Ok {
		w.WriteHeader(403)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: "Wrong auth"})
		w.Write(b)
		return
	}

	if input.Film_id == "" || input.User_id == "" || input.User_id != auth_res.Auth.Userid {
		err := BadRequest{Error: "missing film_id or user_id"}
		w.WriteHeader(400)
		var b []byte
		b, _ = json.Marshal(err)
		w.Write(b)
		return
	}

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	res := db.SetFilmUserFavorite(auth_res.Auth.Userid, input.Film_id, input.Favorite, client)

	if res.Success {
		b, _ := json.Marshal(res.Favorite)
		w.Write(b)
	} else {
		w.WriteHeader(500)
	}

}

func isFavorite(w http.ResponseWriter, r *http.Request) {
	tmpF := strings.TrimPrefix(r.URL.Path, "/favorite/")

	ids := strings.Split(tmpF, "/")
	if len(ids) < 2 {
		fmt.Println("missing arguments")
		w.WriteHeader(400)
		return
	}

	uid := ids[1]
	if len(uid) < 1 {
		fmt.Println("userid is missing in parameters")
		w.WriteHeader(400)
		return
	}

	fid := ids[0]
	if len(fid) < 1 {
		fmt.Println("filmid is missing in parameters")
		w.WriteHeader(400)
		return
	}

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	fav := db.IsFilmUserFavorite(fid, uid, client)

	b, _ := json.Marshal(fav)
	w.Write(b)
}

//-----------------------------------------------------------------------------------------------------------------------------------------------------
//-------------------------------------------------------------------comments----------------------------------------------------------------------------
//-------------------------------------------------------------------------------------------------------------------------------------------------------

func comment(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}
	switch r.Method {
	case "DELETE":
		commentDelete(w, r)
	case "PUT":
		commentPut(w, r)
	default:
		w.WriteHeader(400)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: fmt.Sprintf("Méthode non acceptée : %s", r.Method)})
		w.Write(b)
	}
}

func commentPut(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var input struct {
		Auth    auth.Auth `json:"auth"`
		User_id string    `json:"user_id"`
		Film_id string    `json:"film_id"`
		Contenu string    `json:"contenu"`
	}
	err := decoder.Decode(&input)
	if err != nil {
		fmt.Println("error decoding input:", err)
		w.WriteHeader(400)
		return
	}

	auth := auth.AuthQuery{
		Userid:   input.Auth.Userid,
		Session:  input.Auth.Session,
		Ret_chan: make(chan auth.AuthResponse),
	}

	fmt.Println("requete : ", auth)

	auth_chan <- auth
	fmt.Println("requete envoyée")
	auth_res := <-auth.Ret_chan

	if !auth_res.Ok {
		w.WriteHeader(403)
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: "Wrong auth"})
		w.Write(b)
		return
	}

	if input.User_id == "" || input.Film_id == "" || input.User_id != auth_res.Auth.Userid {
		err := BadRequest{Error: "missing user_id or film_id"}
		w.WriteHeader(400)
		var b []byte
		b, _ = json.Marshal(err)
		w.Write(b)
		return
	}

	client := db.Connect()

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	res := db.PutComment(auth_res.Auth.Userid, input.Film_id, input.Contenu, client)
	if res.Success {
		w.WriteHeader(200)
		b, _ := json.Marshal(res.Comment)
		w.Write(b)
	} else {
		var b []byte
		b, _ = json.Marshal(BadRequest{Error: "impossible de créer le commentaire"})
		w.Write(b)
		w.WriteHeader(400)
	}

}

func all(w http.ResponseWriter, r *http.Request) {
	enableCORS(&w)
	w.WriteHeader(200)
	fmt.Println("ICI???")
	fmt.Println(r.URL.Path)
}

func launchAuth() {
	auth.Authentificator(auth_chan)
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Add("Access-Control-Allow-Headers", "Content-Type, X-Auth-Token, Origin, Authorization")
	(*w).Header().Set("Context-Type", "application/json")
}

func main() {
	auth_chan = make(chan auth.AuthQuery)
	go launchAuth()
	http.HandleFunc("/", all)
	http.HandleFunc("/login/", login)
	http.HandleFunc("/logout/", logout)
	http.HandleFunc("/user/", handlerUser)
	http.HandleFunc("/comment/", comment)
	http.HandleFunc("/comment/user/", commentsUser)
	http.HandleFunc("/comment/film/", commentsFilm)
	http.HandleFunc("/film/", getFilm)
	http.HandleFunc("/favorites/user/", favorites)
	http.HandleFunc("/favorite/", handleFavorite)
	http.HandleFunc("/recommend/user/", recommendUser)
	http.HandleFunc("/recommend/film/", recommendFilm)
	http.HandleFunc("/search/title/", searchFilm)
	http.HandleFunc("/search/genre/", getFilmByGenre)
	fmt.Println("Prêt à écouter")
	log.Fatal(http.ListenAndServe(":54321", nil))
}
