package v1

import (
    "bbs-game/domain/entities"
    "bbs-game/extensions/domain_context"
    "encoding/json"
    "errors"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
)

type PlayerService interface {
    Create(item entities.Player) (entities.Player, error)
    Get(id string) (entities.Player, error)
}

type PlayerAPI struct {
    service PlayerService
}

func NewPlayerAPI(service PlayerService) *PlayerAPI {
    return &PlayerAPI{service: service}
}

func (api *PlayerAPI) Create(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    var player entities.Player
    err = json.Unmarshal(body, &player)
    if err != nil {
        panic(err)
    }

    createdItem, err := api.service.Create(player)
    if err != nil {
        panic(err)
    }
    bytes, err := json.Marshal(createdItem)
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(bytes)
}

func (api *PlayerAPI) GetPlayerByToken(w http.ResponseWriter, r *http.Request) {
    playerID, err := domain_context.ExtractPlayerID(r.Context())
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    item, err := api.service.Get(playerID)
    if err != nil {
        panic(err)
    }

    bytes, err := json.Marshal(item)
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Type", "application/json")

    w.Write(bytes)
}

func (api *PlayerAPI) Get(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, ok := vars["id"]
    if !ok {
        panic(errors.New("no id"))
    }

    createdItem, err := api.service.Get(id)
    if err != nil {
        panic(err)
    }

    bytes, err := json.Marshal(createdItem)
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Type", "application/json")

    w.Write(bytes)
}
