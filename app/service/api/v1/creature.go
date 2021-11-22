package v1

import (
    "bbs-game/domain/entities"
    "encoding/json"
    "errors"
    "github.com/gorilla/mux"
    "io/ioutil"
    "net/http"
)

type CreatureService interface {
    Create(item entities.Creature) (entities.Creature, error)
    Get(id string) (entities.Creature, error)
}

type CreatureAPI struct {
    service CreatureService
}

func NewCreatureAPI(service CreatureService) *CreatureAPI {
    return &CreatureAPI{service: service}
}

func (api *CreatureAPI) Create(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }

    var creature entities.Creature
    err = json.Unmarshal(body, &creature)
    if err != nil {
        panic(err)
    }

    createdCreature, err := api.service.Create(creature)
    if err != nil {
        panic(err)
    }
    bytes, err := json.Marshal(createdCreature)
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    w.Write(bytes)
}


func (api *CreatureAPI) Get(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, ok := vars["id"]
    if !ok {
        panic(errors.New("no id"))
    }

    creature, err := api.service.Get(id)
    if err != nil {
        panic(err)
    }

    bytes, err := json.Marshal(creature)
    if err != nil {
        panic(err)
    }
    w.Header().Set("Content-Type", "application/json")

    w.Write(bytes)
}
