package database

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	api "serveur/api"
	t "serveur/types"
)

const uri = "mongodb+srv://pc3r:pass@cluster0.tm95mok.mongodb.net/?retryWrites=true&w=majority"

func Connect() *mongo.Client {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)

	if err != nil {
		panic(err)
	}

	return client
}

func getUser(id string, client *mongo.Client) t.Result_user {
	coll := client.Database("films").Collection("users")
	var result t.User
	err := coll.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(&result)
	if err != nil {
		fmt.Println("error in getUser :", err, " with id :", id)
		return t.Result_user{Success: false}
	}
	return t.Result_user{Success: true, User: result}
}

func GetUserOut(id string, client *mongo.Client) t.Result_user_out {
	result := getUser(id, client)

	if !result.Success {
		return t.Result_user_out{Success: false}
	}

	user := t.User_out{
		Id:         result.User.Id,
		Username:   result.User.Username,
		Favoris_id: result.User.Favoris_id,
	}

	return t.Result_user_out{Success: true, User: user}
}

func GetUserByMail(id string, client *mongo.Client) t.Result_user {
	coll := client.Database("films").Collection("users")
	var result t.User
	err := coll.FindOne(context.TODO(), bson.D{{Key: "email", Value: id}}).Decode(&result)
	if err != nil {
		fmt.Println("error in getUser :", err, " with id :", id)
		return t.Result_user{Success: false}
	}
	return t.Result_user{Success: true, User: result}
}

func AddUser(user t.User, client *mongo.Client) t.Result_user {
	coll := client.Database("films").Collection("users")
	opts := options.Find().SetSort(bson.D{{Key: "id", Value: -1}}).SetLimit(1)
	cursor, e := coll.Find(context.TODO(), bson.D{}, opts)

	var results []t.User
	var count int

	if e != nil {
		fmt.Println(e)
		count = 0
	} else {
		if err := cursor.All(context.TODO(), &results); err != nil {
			fmt.Println("erreur in addUser :", err)
			return t.Result_user{Success: false}
		}

		lastUser := results[0]
		scount := strings.TrimPrefix(lastUser.Id, "u")
		count, _ = strconv.Atoi(scount)
		count += 1
	}
	user.Id = fmt.Sprintf("u%d", count)
	user.Favoris_id = make([]string, 0)
	_, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		fmt.Println("error in addUser :", err, " with user :", user)
		return t.Result_user{Success: false}
	}
	return t.Result_user{Success: true, User: user}
}

func DeleteUser(id string, client *mongo.Client) bool {
	coll := client.Database("films").Collection("users")
	_, err := coll.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}})
	if err != nil {
		fmt.Println("error in deleteUser :", err, " with id :", id)
		return false
	}
	coll = client.Database("films").Collection("comment")
	_, err = coll.DeleteMany(context.TODO(), bson.D{{Key: "id", Value: id}})
	if err != nil {
		fmt.Println("error in deleteUser (comments) :", err, " with id :", id)
		return false
	}

	return true
}

func GetCommentsUser(id_user string, client *mongo.Client) []t.Commentaire_user {
	var coll = client.Database("films").Collection("comment")
	var results []t.Commentaire
	opts := options.Find().SetSort(bson.D{{Key: "id", Value: -1}})
	var cursor, err = coll.Find(context.TODO(), bson.D{{Key: "user_id", Value: id_user}}, opts)

	if err != nil {
		fmt.Println("error in getCommentsUser :", err, " with id :", id_user)
		return make([]t.Commentaire_user, 0)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return make([]t.Commentaire_user, 0)
	}

	var comments = make([]t.Commentaire_user, 0)

	for _, com := range results {
		f := api.GetFilm(com.Film_id)
		if f.Success {
			comment_no_film := t.Commentaire_user{
				Id:      com.Id,
				User_id: com.User_id,
				Contenu: com.Contenu,
				Film:    f.Film,
			}
			comments = append(comments, comment_no_film)
		}
	}

	return comments
}

func GetCommentsFilm(id_film string, client *mongo.Client) []t.Commentaire_film {
	var coll = client.Database("films").Collection("comment")
	var results []t.Commentaire
	opts := options.Find().SetSort(bson.D{{Key: "id", Value: -1}})
	var cursor, err = coll.Find(context.TODO(), bson.D{{Key: "film_id", Value: id_film}}, opts)

	if err != nil {
		fmt.Println("error in getCommentsFilm :", err, " with id :", id_film)
		return make([]t.Commentaire_film, 0)
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return make([]t.Commentaire_film, 0)
	}

	var comments = make([]t.Commentaire_film, 0)

	for _, com := range results {
		user := GetUserOut(com.User_id, client)
		if user.Success {
			comment_no_film := t.Commentaire_film{
				Id:      com.Id,
				User:    user.User,
				Contenu: com.Contenu,
				Film_id: com.Film_id,
			}
			comments = append(comments, comment_no_film)
		}
	}

	return comments
}

func SetFilmUserFavorite(uid string, fid string, fav bool, client *mongo.Client) t.Result_fav {
	if fav {
		return addUserFavorite(uid, fid, client)
	}
	return removeUserFavorite(uid, fid, client)
}

func addUserFavorite(id_user string, id_film string, client *mongo.Client) t.Result_fav {
	var favs = GetUserFavoriteIds(id_user, client)
	index := indexString(favs, id_film)
	fav_res := t.Film_Fav{User_id: id_user, Film_id: id_film, Favorite: true}

	if index != -1 {
		return t.Result_fav{Success: true, Favorite: fav_res}
	}

	favs = append(favs, id_film)

	var coll = client.Database("films").Collection("users")

	filter := bson.D{{Key: "id", Value: id_user}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "favoris_id", Value: favs}}}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return t.Result_fav{Success: false}
	}

	return t.Result_fav{Success: true, Favorite: fav_res}
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func removeUserFavorite(id_user string, id_film string, client *mongo.Client) t.Result_fav {
	var favs = GetUserFavoriteIds(id_user, client)
	index := indexString(favs, id_film)

	fav_res := t.Film_Fav{User_id: id_user, Film_id: id_film, Favorite: false}

	if index == -1 {
		return t.Result_fav{Success: true, Favorite: fav_res}
	}

	favs = remove(favs, index)

	var coll = client.Database("films").Collection("users")

	filter := bson.D{{Key: "id", Value: id_user}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "favoris_id", Value: favs}}}}

	_, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return t.Result_fav{Success: false}
	}

	return t.Result_fav{Success: true, Favorite: fav_res}
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

func GetUserFavoriteIds(id_user string, client *mongo.Client) []string {
	user := getUser(id_user, client)
	if user.Success {
		return user.User.Favoris_id
	}
	return make([]string, 0)
}

func IsFilmUserFavorite(id_film string, id_user string, client *mongo.Client) t.Film_Fav {
	return t.Film_Fav{
		User_id:  id_user,
		Film_id:  id_film,
		Favorite: containsString(GetUserFavoriteIds(id_user, client), id_film),
	}
}

func GetUserFavoriteFilms(id_user string, client *mongo.Client) []t.Film {
	favs_id := GetUserFavoriteIds(id_user, client)
	var films []t.Film
	films = make([]t.Film, 0)
	for _, f := range favs_id {
		film := api.GetFilm(f)
		if film.Success {
			films = append(films, film.Film)
		}
	}
	return films
}

func DeleteComment(id string, userid string, client *mongo.Client) bool {
	coll := client.Database("films").Collection("comment")
	_, err := coll.DeleteOne(context.TODO(), bson.D{{Key: "id", Value: id}, {Key: "user_id", Value: userid}})
	if err != nil {
		fmt.Println("error in deleteComment :", err, " with id : ", id)
		return false
	}

	return true
}

func PutComment(id_user string, id_film string, contenu string, client *mongo.Client) t.Result_com {
	coll := client.Database("films").Collection("comment")
	opts := options.Find().SetSort(bson.D{{Key: "id", Value: -1}}).SetLimit(1)
	cursor, e := coll.Find(context.TODO(), bson.D{}, opts)

	var results []t.User
	var count int

	if e != nil {
		fmt.Println(e)
		count = 0
	} else {
		if err := cursor.All(context.TODO(), &results); err != nil {
			fmt.Println("erreur in PutComment :", err)
			return t.Result_com{Success: false}
		}

		lastUser := results[0]
		fmt.Println("number", lastUser.Id)
		scount := strings.TrimPrefix(lastUser.Id, "c")
		fmt.Println("number", scount)
		count, _ = strconv.Atoi(scount)
		count += 1
	}
	fmt.Println("ajout commentaire", count)
	com := t.Commentaire{
		Id:      fmt.Sprintf("c%d", count),
		User_id: id_user,
		Contenu: contenu,
		Film_id: id_film,
	}

	_, err := coll.InsertOne(context.TODO(), com)

	retour := t.Commentaire_film{
		Id:      com.Id,
		Contenu: com.Contenu,
		Film_id: com.Film_id,
		User:    GetUserOut(id_user, client).User,
	}

	if err != nil {
		fmt.Println("erreur dans putComment :")
		fmt.Println(err)

		return t.Result_com{Success: false}
	}

	user := getUser(com.User_id, client)
	if user.Success {
		return t.Result_com{Success: true, Comment: retour}
	}
	return t.Result_com{Success: false}
}
