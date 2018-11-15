package repos

import (
	"crypto/sha256"
	"fmt"

	"github.com/couchbase/gocb"
)

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	Country string `json:"country"`
	Postal  string `json:"postal"`
}

type User struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Sha      string `json:"session"`
	Address  string `json:"address"`
}

// var users map[string]User
var cluster *gocb.Cluster
var bucket *gocb.Bucket

func init() {
	// users = make(map[string]User)
	// var err error
	cluster, _ = gocb.Connect("couchbase://localhost")
	_ = cluster.Authenticate(gocb.PasswordAuthenticator{
		Username: "Wisa",
		Password: "couchAdmin",
	})

	bucket, _ = cluster.OpenBucket("HelloBucket", "")
}

func computeSha(username, password string) string {
	unencrypt := "%s:%s"
	unencrypt = fmt.Sprintf(unencrypt, username, password)

	sha := sha256.Sum256([]byte(unencrypt))
	return string(sha[:])
}

// ValidateUser : this function is intended to used for login
func ValidateUser(username, password string) (User, bool) {
	user, found := FindUser(username)

	// no user
	if !found {
		return User{}, false
	}

	// check user validation
	sha := computeSha(username, password)

	// user.Sha = string([]byte(user.Sha))

	// for _, c := range sha {
	// 	fmt.Printf("%#U", c)
	// }
	// fmt.Println("\n")
	// for _, c := range user.Sha {
	// 	fmt.Printf("%#U", c)
	// }
	// fmt.Println("")

	sha_r1 := []rune(user.Sha)
	sha_r2 := []rune(sha)
	isShaEqual := len(sha_r1) == len(sha_r2)

	if isShaEqual {
		for i, r1 := range sha_r1 {
			// fmt.Println(i)
			r2 := sha_r2[i]
			if r1 != r2 {
				isShaEqual = false
			}
			i++
		}
	}

	if isShaEqual {
		return user, true
	} else {
		return User{}, false
	}
}

func FindUser(username string) (User, bool) {
	// user := users[username]
	// if username == "" { // TODO: better check this somewhere else
	// 	return User{}, false
	// }

	user := User{}
	_, err := bucket.Get(username, &user)

	if err == nil {
		return user, true
	} else {
		return User{}, false
	}
}

func Register(username, password, email string) User {
	sha := computeSha(username, password)

	// create new user
	user := User{
		Username: username,
		Email:    email,
		Sha:      sha,
	}
	// users[username] = user

	bucket.Upsert(username, user, 0)

	return user
}
