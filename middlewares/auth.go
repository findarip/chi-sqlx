package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret []byte

var JWT_SIGNING_METHOD = jwt.SigningMethodHS256
var JWT_SIGNATURE_KEY = []byte("secret")

func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Signing method invalid")
			} else if method != JWT_SIGNING_METHOD {
				return nil, fmt.Errorf("Signing method invalid")
			}

			return JWT_SIGNATURE_KEY, nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(context.Background(), "userInfo", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
		// if tokenString == "" {
		// 	// w.Write([]byte(`empty authorization`))
		// 	web.WriteToResponseBody(w, 401, http.StatusText(http.StatusUnauthorized), nil, "Token Not Found", nil)
		// 	return
		// }
		// fmt.Println("string /n", tokenString)

		// splitToken := strings.Split(tokenString, " ")
		// if splitToken[1] != "" {
		// 	w.Write([]byte(`empty token`))
		// 	web.WriteToResponseBody(w, 401, http.StatusText(http.StatusUnauthorized), nil, "empty token", nil)

		// }

		// token, err := jwt.Parse(splitToken[1], func(token *jwt.Token) (interface{}, error) {
		// 	// Don't forget to validate the alg is what you expect:
		// 	fmt.Println("token /n", token)

		// 	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		// 	}

		// 	// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		// 	return hmacSampleSecret, nil
		// })

		// if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// 	// fmt.Println(claims["exp"], claims["username"])
		// 	fmt.Printf("claim data %v ", claims)
		// } else {
		// 	fmt.Printf("masuk eror")
		// 	fmt.Println(err)
		// }

		// return

		// splitToken := strings.Split(tokenString, " ")
		// if splitToken[1] != "" {
		// w.Write([]byte(`empty token`))
		// web.WriteToResponseBody(w, 401, http.StatusText(http.StatusUnauthorized), nil, "empty token", nil)

		// }

		// next.ServeHTTP(w, r)
		// fmt.Println("string /n", tokenString)

	})
}
