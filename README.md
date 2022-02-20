# verifier

Golang JWT token verifier with storage(default Redis)

## Usage

```shell
go get -u github.com/xandercheung/verifier
```

```go
package main

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/xandercheung/verifier"
	"time"
)

func main() {
	var redisClient *redis.Client
	// ...
	// use redis(github.com/go-redis/redis/v8) as token storage, you can create new TokenStorage
	tokenStorage := verifier.RedisStorage{Client: redisClient}

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

	// creates a new Verifier
	//v := verifier.New(
	//	"your key",
	//	&tokenStorage,
	//	verifier.WithTempTokenExpireDuration(time.Minute),
	//)

	var sourceId uint = 1
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
	err = verifier.DestroyToken(sourceId, claims.Uuid)
	if err != nil {
		fmt.Println(err)
	}

	// destroy all token of sourceId
	err = verifier.DestroyAllToken(sourceId)
	if err != nil {
		fmt.Println(err)
	}

	// refresh token manually
	newToken, err = verifier.RefreshToken(claims.SourceId, claims.Uuid, claims.Data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(newToken)
}

```
