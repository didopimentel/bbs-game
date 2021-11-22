package middlewares

import (
    "bbs-game/domain/account"
    "context"
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "net/http"
    "strings"
)

func Authentication(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/api/v1/accounts" || r.URL.Path == "/api/v1/accounts/login" {
            next.ServeHTTP(w, r)
        } else {
            authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
            if len(authHeader) != 2 {
                fmt.Println("Malformed token")
                w.WriteHeader(http.StatusUnauthorized)
                w.Write([]byte("Malformed Token"))
            } else {
                jwtToken := authHeader[1]
                token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
                    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                        return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
                    }
                    return account.JwtKey, nil
                })

                if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
                    ctx := context.WithValue(r.Context(), "props", claims)
                    // Access domain_context values in handlers like this
                    // props, _ := r.Context().Value("props").(jwt.MapClaims)
                    next.ServeHTTP(w, r.WithContext(ctx))
                } else {
                    fmt.Println(err)
                    w.WriteHeader(http.StatusUnauthorized)
                    w.Write([]byte("Unauthorized"))
                }
            }

        }
    })
}
