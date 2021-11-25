package account

import (
    "bbs-game/domain/entities"
    "github.com/dgrijalva/jwt-go"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "time"
)

var JwtKey = []byte("my_secret_key")

type Service struct {
    db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
    return &Service{
        db: db,
    }
}

type CreateInput struct {
    Email string
    Password string
    PlayerName string
}
type CreateOutput struct {
    AccountID string
    Email string
    PlayerID string
    PlayerName string
}
func (s *Service) Create(input CreateInput) (CreateOutput, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
    if err != nil {
        return CreateOutput{}, err
    }

    account := entities.Account{
        Email: input.Email,
        Password: string(bytes),
    }

    player := entities.Player{
        Name:       input.PlayerName,
        Level:      1,
        Damage: "1d6+5",
        Experience: 0,
        HP:         100,
        TotalHP:    100,
    }

    err = s.db.Transaction(func(tx *gorm.DB) error {
        tx = s.db.Table("accounts").Create(&account)
        if tx.Error != nil {
            return tx.Error
        }

        player.AccountID = account.ID

        tx = s.db.Table("players").Create(&player)
        if tx.Error != nil {
            return tx.Error
        }

        return nil
    })
    if err != nil {
        return CreateOutput{}, err
    }


    return CreateOutput{
        AccountID: account.ID,
        Email:      account.Email,
        PlayerID:   player.ID,
        PlayerName: player.Name,
    }, nil

}

type LoginInput struct {
    Email string
    Password string
}
type LoginOutput struct {
    Token string
    ExpiresAt int64
}
type Claims struct {
    Email string
    PlayerID string
    jwt.StandardClaims
}
func (s *Service) Login(input LoginInput) (LoginOutput, error) {
    account := entities.Account{}

    tx := s.db.Table("accounts").Where("email = ?", input.Email).First(&account)
    if tx.Error != nil {
        return LoginOutput{}, tx.Error
    }

    if !CheckPasswordHash(input.Password, account.Password) {
        return LoginOutput{}, ErrInvalidPassword
    }


    player := entities.Player{}
    tx = s.db.Table("players").Where("account_id = ?", account.ID).First(&player)
    if tx.Error != nil {
        return LoginOutput{}, tx.Error
    }

    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &Claims{
        Email: input.Email,
        PlayerID: player.ID,
        StandardClaims: jwt.StandardClaims{
            // In JWT, the expiry time is expressed as unix milliseconds
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(JwtKey)
    if err != nil {
        return LoginOutput{}, err
    }

    return LoginOutput{
        Token:     tokenString,
        ExpiresAt: expirationTime.Unix(),
    }, nil
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

