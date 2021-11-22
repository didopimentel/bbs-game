package domain_context

import (
    "context"
    "errors"
    "github.com/dgrijalva/jwt-go"
)

func ExtractPlayerID(ctx context.Context) (string, error) {
    props := ctx.Value("props").(jwt.MapClaims)
    p, ok := props["PlayerID"]
    if !ok {
        return "", errors.New("player id not present in token")
    }

    return p.(string), nil
}
