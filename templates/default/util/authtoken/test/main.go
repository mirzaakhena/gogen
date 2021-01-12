package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mirzaakhena/govue/templates/shared/util/authtoken"
)

func main() {

	// sk := authtoken.GenerateSecretKey()
	// fmt.Printf("%v\n\n", sk)

	secretKey := []byte("0bde05a8c622a12e6932f0748055261085f6a557c6dccd7be61073795532bd7ab7f55887aa01efba17790397304fc6dd6efed2d1a8389fdf9336a40d08fd220f")
	jwtToken := authtoken.New("Diriku", secretKey)

	tokenString2 := jwtToken.Generate("LOGIN", "USER", struct {
		Name string `json:"name"`
	}{Name: "mirza"}, 20*time.Hour)
	fmt.Printf("%v\n\n", tokenString2)

	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJVU0VSIiwiZXhwIjoxNjA4NDcwMTU0LCJpYXQiOjE2MDgzOTgxNTQsImlzcyI6IkRpcmlrdSIsIm5iZiI6MTYwODM5ODE1NCwic3ViIjoiTE9HSU4iLCJleHRlbmRJbmZvIjp7Im5hbWUiOiJtaXJ6YSJ9fQ.dmWDq1e0yD8LYPK4UCYiq9YXzVv2PiScKQi9Wx-Evac"
	claims, valid := jwtToken.Validate("LOGIN", tokenString)
	if !valid {
		fmt.Printf("not valid %v\n", claims)
		os.Exit(0)
	}
	fmt.Printf("valid %v\n", claims)

	// claims.

}
