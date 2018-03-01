package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/natethinks/instruu-api/internal/respond"
	"github.com/natethinks/instruu-api/internal/store"
)

// NewJWT accepts a user and creates a JWT representing that user
func NewJWT(user store.User) (string, error) {
	fmt.Println(user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
	})

	tokenString, err := token.SignedString([]byte("my_not_secret_key"))

	return tokenString, err
}

// CheckJWT retrieves user info from the JWT and stores it in the header for an endpoint to use, but allows access without a JWT
// this is mostly for get requests since post requests will almost always require a user to be logged in as something
func CheckJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check to see if the web token is valid, should be stored in a cookie since that's what i'm going to use
		// also check to make sure the request is coming from a valid resource, XSRF CSRF checks

		fmt.Println(r.Cookies())

		authCookie, _ := r.Cookie("auth")
		tokenString := authCookie.String()
		fmt.Println(tokenString)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte("my_not_secret_key"), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println(claims["foo"], claims["nbf"])
		} else {
			fmt.Println(err)
		}

		h.ServeHTTP(w, r)
	})
}

// SecureCheckJWT blocks access to an endpoint unless a user is logged in with a valid JWT
func SecureCheckJWT(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check to see if the web token is valid, should be stored in a cookie since that's what i'm going to use
		// also check to make sure the request is coming from a valid resource, XSRF CSRF checks

		fmt.Println(r.Cookies())

		jwt, err := r.Cookie("auth")
		if err != nil {
			// no cookie found by that name
			w.WriteHeader(http.StatusUnauthorized)
			respond.JSON(w, errors.New("Missing or invalid token"))
		}
		fmt.Println(jwt)

		// if the validation fails
		if 1 == 1 {
			log.Println("invalid JWT")

			err := errors.New("Missing or invalid token")

			w.WriteHeader(http.StatusUnauthorized)
			respond.JSON(w, err)
			return
			// if the validation succeeds. continue with the serving of the endpoint that this wraps
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
