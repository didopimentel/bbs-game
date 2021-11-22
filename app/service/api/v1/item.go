package v1

import (
    "bbs-game/domain/entities"
    "encoding/json"
    "errors"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
)

type ItemService interface {
    Create(item entities.Item) (entities.Item, error)
    Get(id string) (entities.Item, error)
}

type ItemAPI struct {
    service ItemService
}

func NewItemAPI(service ItemService) *ItemAPI {
    return &ItemAPI{service: service}
}

func (api *ItemAPI) Create(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    var item entities.Item
    err = json.Unmarshal(body, &item)
    if err != nil {
        panic(err)
    }

    createdItem, err := api.service.Create(item)
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


func (api *ItemAPI) Get(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, ok := vars["id"]
    if !ok {
        panic(errors.New("no id"))
    }

    item, err := api.service.Get(id)
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
