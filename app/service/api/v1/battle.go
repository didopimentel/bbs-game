package v1

import (
	"bbs-game/app/service/api/v1/dto"
	"bbs-game/domain/battle"
	"bbs-game/domain/entities"
	"bbs-game/extensions/domain_context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type BattleService interface {
	Create(input battle.CreateInput) (battle.CreateOutput, error)
	CreateAction(battle.CreateActionInput) (entities.BattleAction, error)
	Get(string) (entities.Battle, error)
	ListActionsByBattle(string) ([]entities.BattleAction, error)
	GenerateNextRound(battle.GenerateNextRoundInput) (battle.GenerateNextRoundOutput, error)
}

type BattleAPI struct {
	service BattleService
}

func NewBattleAPI(service BattleService) *BattleAPI {
	return &BattleAPI{service: service}
}

func toBattleParticipantsEntities(participants []dto.BattleParticipant) []entities.BattleParticipant {
	p := make([]entities.BattleParticipant, 0)
	for _, bp := range participants {
		p = append(p, entities.BattleParticipant{
			ParticipantType: entities.BattleParticipantType(bp.ParticipantType),
			ParticipantID:   bp.ParticipantID,
		})
	}

	return p
}

func (api *BattleAPI) Create(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var battleDTO dto.BattleCreateInput
	err = json.Unmarshal(body, &battleDTO)
	if err != nil {
		panic(err)
	}

	input := battle.CreateInput{
		BattleParticipants: toBattleParticipantsEntities(battleDTO.BattleParticipants),
	}
	createdBattle, err := api.service.Create(input)
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(createdBattle)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(bytes)
}

func (api *BattleAPI) GenerateNextRound(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	battleID, ok := vars["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	playerID, err := domain_context.ExtractPlayerID(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	input := battle.GenerateNextRoundInput{
		BattleID: battleID,
		PlayerID: playerID,
	}
	response, err := api.service.GenerateNextRound(input)
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(dto.GenerateNextRoundResponse{
		BattleActions:    response.BattleActions,
		PlayerDied:       response.PlayerDied,
		CreatureDied:     response.CreatureDied,
		ExperienceGained: response.ExperienceGained,
		GainedLevel:      response.GainedLevel,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (api *BattleAPI) CreateAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	battleID, ok := vars["id"]
	if !ok {
		panic(errors.New("no id"))
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var battleActionDTO dto.CreateActionInput
	err = json.Unmarshal(body, &battleActionDTO)
	if err != nil {
		panic(err)
	}

	createdBattle, err := api.service.CreateAction(battle.CreateActionInput{
		BattleID:   battleID,
		CauserID:   battleActionDTO.CauserID,
		TargetID:   battleActionDTO.TargetID,
		ActionType: entities.BattleActionType(battleActionDTO.ActionType),
	})
	if err != nil {
		panic(err)
	}
	bytes, err := json.Marshal(createdBattle)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(bytes)
}

func (api *BattleAPI) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		panic(errors.New("no id"))
	}

	battle, err := api.service.Get(id)
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(battle)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(bytes)
}

func (api *BattleAPI) ListActionsByBattle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		panic(errors.New("no id"))
	}

	actions, err := api.service.ListActionsByBattle(id)
	if err != nil {
		panic(err)
	}

	bytes, err := json.Marshal(actions)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(bytes)
}

