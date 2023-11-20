package auth

import (
	"fmt"
	"math/rand"
)

// new_session && delete_session == false
type AuthQuery struct {
	Userid         string
	Session        string
	New_session    bool
	Delete_session bool
	Ret_chan       chan AuthResponse
}

type Auth struct {
	Userid  string `json:"userid"`
	Session string `json:"session"`
}

type AuthResponse struct {
	Auth Auth
	Ok   bool
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func Authentificator(input chan AuthQuery) {
	var userid_to_session = make(map[string]string)

	for {
		request := <-input
		fmt.Println("new request : ", request)

		// pas de userid
		if request.Userid == "" {
			request.Ret_chan <- AuthResponse{
				Ok: false,
			}
			continue
		}

		if request.New_session {
			session := RandomString(32)
			userid_to_session[request.Userid] = session
			request.Ret_chan <- AuthResponse{
				Auth: Auth{
					Userid:  request.Userid,
					Session: session,
				},
				Ok: true,
			}
		} else if request.Delete_session {
			if val, ok := userid_to_session[request.Userid]; ok && val == request.Session {
				delete(userid_to_session, request.Userid)
				request.Ret_chan <- AuthResponse{
					Auth: Auth{
						Userid:  request.Userid,
						Session: request.Session,
					},
					Ok: true,
				}
			} else {
				request.Ret_chan <- AuthResponse{
					Auth: Auth{
						Userid:  request.Userid,
						Session: request.Session,
					},
					Ok: false,
				}
			}
		} else {
			val, ok := userid_to_session[request.Userid]
			request.Ret_chan <- AuthResponse{
				Auth: Auth{
					Userid:  request.Userid,
					Session: request.Session,
				},
				Ok: ok && val == request.Session,
			}
		}
		fmt.Println("request handled")
	}
}
