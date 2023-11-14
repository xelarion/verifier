package main

import (
	"fmt"
	"time"

	"github.com/NoahCodeGG/verifier"
	"github.com/go-redis/redis/v8"
	storage "github.com/xandercheung/verifier-redis-storage"
)

func main() {
	// use github.com/xandercheung/verifier-redis-storage(github.com/go-redis/redis/v8) as token storage,
	// you can create new TokenStorage
	// init your redisClient
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "",
		DB:       0,
	})
	tokenStorage := storage.RedisStorage{Cli: redisClient}

	// use default verifier
	verifier.InitDefaultVerifier(
		"your key",
		&tokenStorage,
		// options
		verifier.WithSourceName("user"),                   // default "user"
		verifier.WithTokenExpireDuration(10*time.Minute),  // default 15 minutes
		verifier.WithAuthExpireDuration(2*time.Hour),      // default 3 hours
		verifier.WithTempTokenExpireDuration(time.Minute), // default 30 seconds
	)

	// or creates a new Verifier
	// v := verifier.New(
	//	"your key",
	//	&tokenStorage,
	//	verifier.WithTempTokenExpireDuration(time.Minute),
	// )

	sourceId := uint(1)
	token, err := verifier.CreateToken(sourceId, map[string]interface{}{"foo": "bar"})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token)

	// verify token and refresh token(if necessary)
	claims, newToken, err := verifier.VerifyToken(token)
	if err != nil {
		// validation failed
		fmt.Println(err)
	}
	fmt.Println(claims)
	fmt.Println(newToken)

	// verify token only
	claims, isAuthorized := verifier.IsTokenAuthorized(token)
	fmt.Println(claims)
	fmt.Println(isAuthorized)

	// destroy token by sourceId and uuid
	err = verifier.DestroyToken(sourceId, claims.UUID)
	if err != nil {
		fmt.Println(err)
	}

	// destroy all token of sourceId
	err = verifier.DestroyAllToken(sourceId)
	if err != nil {
		fmt.Println(err)
	}

	// refresh token manually
	newToken, err = verifier.RefreshToken(claims.SourceID, claims.UUID, claims.Data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(newToken)
}
