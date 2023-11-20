package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	t "serveur/types"
)

const rapid_api_kei = "187ee318demsh976606a5e61f41ap1178e8jsnd22cb78fec44"
const rapid_api_host = "moviesminidatabase.p.rapidapi.com"

func shuffle(src []string) []string {
	dest := make([]string, len(src))
	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}

	return dest
}

func containsString(array []string, word string) bool {
	return indexString(array, word) != -1
}

func indexString(array []string, word string) int {
	for i, w := range array {
		if w == word {
			return i
		}
	}
	return -1
}

// recherche un film par son id
func GetFilm(id string) t.Result_film {
	if id == "" {
		return t.Result_film{Success: false}
	}

	url := fmt.Sprintf("https://moviesminidatabase.p.rapidapi.com/movie/id/%s", id)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", rapid_api_kei)
	req.Header.Add("X-RapidAPI-Host", rapid_api_host)

	res, e := http.DefaultClient.Do(req)

	if e != nil {
		return t.Result_film{Success: false}
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var objmap map[string]json.RawMessage
	err := json.Unmarshal(body, &objmap)
	if err != nil {
		panic(err)
	}

	var title string
	var description string
	var image string

	var results map[string]json.RawMessage
	json.Unmarshal(objmap["results"], &results)

	if len(results) < 1 {
		return t.Result_film{Success: false}
	}

	json.Unmarshal(results["title"], &title)
	json.Unmarshal(results["description"], &description)
	json.Unmarshal(results["image_url"], &image)

	film := t.Film{
		Id:          id,
		Title:       title,
		Description: description,
		Image:       image,
	}

	return t.Result_film{Success: true, Film: film}
}

func getGenresFilm(id string) []string {
	if id == "" {
		return make([]string, 0)
	}

	url := fmt.Sprintf("https://moviesminidatabase.p.rapidapi.com/movie/id/%s", id)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", rapid_api_kei)
	req.Header.Add("X-RapidAPI-Host", rapid_api_host)

	res, e := http.DefaultClient.Do(req)

	if e != nil {
		return make([]string, 0)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var objmap map[string]json.RawMessage
	err := json.Unmarshal(body, &objmap)
	if err != nil {
		panic(err)
	}

	var results map[string]json.RawMessage
	json.Unmarshal(objmap["results"], &results)

	if len(results) < 1 {
		return make([]string, 0)
	}

	var genres []struct {
		Id    int    `json:"id"`
		Genre string `json:"genre"`
	}

	json.Unmarshal(results["gen"], &genres)

	gens := make([]string, 0)
	for _, g := range genres {
		gens = append(gens, g.Genre)
	}

	return gens
}

func GetFilmsByGenre(genre string) []string {
	url := fmt.Sprintf("https://moviesminidatabase.p.rapidapi.com/movie/byGen/%s", genre)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", rapid_api_kei)
	req.Header.Add("X-RapidAPI-Host", rapid_api_host)

	res, e := http.DefaultClient.Do(req)

	if e != nil {
		return make([]string, 0)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var objmap map[string]json.RawMessage
	err := json.Unmarshal(body, &objmap)
	if err != nil {
		panic(err)
	}

	var results []struct {
		Imdb_id string `json:"imdb_id"`
		Title   string `json:"title"`
	}
	json.Unmarshal(objmap["results"], &results)

	if len(results) < 1 {
		return make([]string, 0)
	}

	films := make([]string, 0)

	for _, f := range results {
		films = append(films, f.Imdb_id)
	}

	return films
}

// recherche film par titre
func SearchFilm(title string) []string {
	url := fmt.Sprintf("https://moviesminidatabase.p.rapidapi.com/movie/imdb_id/byTitle/%s/", title)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", rapid_api_kei)
	req.Header.Add("X-RapidAPI-Host", rapid_api_host)

	res, e := http.DefaultClient.Do(req)

	if e != nil {
		return make([]string, 0)
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var objmap map[string]json.RawMessage
	err := json.Unmarshal(body, &objmap)
	if err != nil {
		panic(err)
	}

	var results []struct {
		Imdb_id string `json:"imdb_id"`
		Title   string `json:"title"`
	}
	json.Unmarshal(objmap["results"], &results)

	if len(results) < 1 {
		return make([]string, 0)
	}

	films := make([]string, 0)

	for _, f := range results {
		films = append(films, f.Imdb_id)
	}

	return films
}

func GetPopularFilms(max int) t.Result_films {
	url := "https://moviesminidatabase.p.rapidapi.com/movie/order/byPopularity/"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", rapid_api_kei)
	req.Header.Add("X-RapidAPI-Host", rapid_api_host)

	res, e := http.DefaultClient.Do(req)

	if e != nil {
		fmt.Println("error api : error request", e)
		return t.Result_films{Success: false}
	}

	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var objmap map[string]json.RawMessage
	err := json.Unmarshal(body, &objmap)
	if err != nil {
		panic(err)
	}

	var results []struct {
		Imdb_id    string `json:"imdb_id"`
		Title      string `json:"title"`
		Popularity int    `json:"popularity"`
	}

	json.Unmarshal(objmap["results"], &results)

	if len(results) < 1 {
		fmt.Println("error api not enough films")
		return t.Result_films{Success: false}
	}

	fmt.Println("nb popular films to populate :", max)

	nb := 0
	films := make([]t.Film, 0)
	for _, id := range results {
		film := GetFilm(id.Imdb_id)
		if film.Success {
			films = append(films, film.Film)
			nb++
			if nb >= max {
				return t.Result_films{Success: true, Films: films}
			}
		}
	}

	fmt.Println("nb popular films :", len(films))

	return t.Result_films{Success: true, Films: films}
}

func GetSimilarFilms(ids []string, max int) t.Result_films {
	if len(ids) == 0 {
		fmt.Println("no film ids")
		return t.Result_films{
			Success: true,
			Films:   make([]t.Film, 0),
		}
	}

	genres := make([]string, 0)
	for _, id := range ids {
		for _, g := range getGenresFilm(id) {
			if !containsString(genres, g) {
				genres = append(genres, g)
			}
		}
	}

	films := make([]string, 0)
	for _, g := range genres {
		gs := GetFilmsByGenre(g)
		for _, f := range gs {
			if !containsString(ids, f) {
				films = append(films, f)
			}
		}
	}

	films = shuffle(films)
	nb := 0
	result_films := make([]t.Film, 0)
	for _, f := range films {
		film := GetFilm(f)
		if film.Success {
			result_films = append(result_films, film.Film)
			nb++
			if nb >= max {
				break
			}
		}
	}

	return t.Result_films{
		Success: true,
		Films:   result_films,
	}
}
