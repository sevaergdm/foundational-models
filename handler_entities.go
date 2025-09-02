package main

import (
	"net/http"
)

func (cfg *apiConfig) handlerGetEntities(w http.ResponseWriter, r *http.Request) {
	entities := []CoreEntity{}
	for _, value := range cfg.entitiesCache {
		entities = append(entities, value)
	}
	respondWithPrettyJSON(w, http.StatusOK, entities)
}

func (cfg *apiConfig) handlerGetEntity(w http.ResponseWriter, r *http.Request) {
	entityName := r.PathValue("entityName")
	
	requestedEntity, ok := cfg.entitiesCache[entityName]
	if !ok {
		respondWithError(w, http.StatusNotFound, "Unable to find requested entity", nil)
		return
	}
	respondWithPrettyJSON(w, http.StatusOK, requestedEntity)
}
