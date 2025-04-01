package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/DIVIgor/chirpy/internal/auth"
	"github.com/google/uuid"
)

const userUpgradeEvent string = "user.upgraded"

type upgradeUserReqBody struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

// Upgrade user plan depending on specific request body and respond with appropriate status code
func (cfg *apiConfig) handlerUpgradeUserPlan(writer http.ResponseWriter, req *http.Request) {
	// check API token
	reqToken, err := auth.GetBearerToken(req.Header, auth.ApiBearer)
	if err != nil {
		respWithErr(writer, http.StatusBadRequest, "Couldn't find token", err)
		return
	}
	if reqToken != cfg.polkaKey {
		respWithErr(writer, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// read request body
	decoder := json.NewDecoder(req.Body)
	reqData := upgradeUserReqBody{}
	err = decoder.Decode(&reqData)
	if err != nil {
		respWithErr(writer, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Ignore events that not match the desired one
	if reqData.Event != userUpgradeEvent {
		writer.WriteHeader(http.StatusNoContent)
		return
	}

	// Upgrade user plan or respond 404
	_, err = cfg.dbQueries.UpgradeUserPlan(req.Context(), reqData.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respWithErr(writer, http.StatusNotFound, "Couldn't find user", err)
		} else {
			respWithErr(writer, http.StatusInternalServerError, "Couldn't upgrade user plan", err)
			return
		}
	}

	writer.WriteHeader(http.StatusNoContent)
}
