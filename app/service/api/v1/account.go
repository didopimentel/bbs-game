package v1

import (
    "bbs-game/app/service/api/v1/dto"
    "bbs-game/domain/account"
    "encoding/json"
    "errors"
    "io/ioutil"
    "net/http"
)

type AccountService interface {
    Create(input account.CreateInput) (account.CreateOutput, error)
    Login(input account.LoginInput) (account.LoginOutput, error)
}

type AccountAPI struct {
    service AccountService
}

func NewAccountAPI(accountService AccountService) *AccountAPI {
    return &AccountAPI{service: accountService}
}

func (api *AccountAPI) Create(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    var accountDTO dto.CreateAccountRequest
    err = json.Unmarshal(body, &accountDTO)
    if err != nil {
        panic(err)
    }

    input := account.CreateInput{
        Email:      accountDTO.Email,
        Password:   accountDTO.Password,
        PlayerName: accountDTO.PlayerName,
    }
    output, err := api.service.Create(input)
    if err != nil {
        panic(err)
    }
    bytes, err := json.Marshal(dto.CreateAccountResponse{
        AccountID:  output.AccountID,
        PlayerName: output.PlayerName,
        PlayerID:   output.PlayerID,
        Email:      output.Email,
    })
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(bytes)
}


func (api *AccountAPI) Login(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    var loginDTO dto.LoginRequest
    err = json.Unmarshal(body, &loginDTO)
    if err != nil {
        panic(err)
    }

    input := account.LoginInput{
        Email:      loginDTO.Email,
        Password:   loginDTO.Password,
    }
    output, err := api.service.Login(input)
    if err != nil {
        if errors.As(account.ErrInvalidPassword, &err) {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }

        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(output.Token))
    w.WriteHeader(http.StatusOK)
}