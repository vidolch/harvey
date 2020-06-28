package controllers

import (
	"github.com/revel/revel"
	"harvey/app/models"
	"harvey/app/models/views"
	"harvey/app/services"
	"log"
)

type Game struct {
	*revel.Controller
}

func (c Game) GetAll() revel.Result {
	gs, err := _getGameService()

	if err != nil {
		log.Fatal(err)
		return c.RenderJSON(err)
	}

	res := gs.GetAll()

	response := views.JsonResponse{}
	response.Data = res

	return c.RenderJSON(response)
}

func (c Game) GetById(id string) revel.Result {
	gs, err := _getGameService()

	if err != nil {
		log.Fatal(err)
		return c.RenderJSON(err)
	}

	res := gs.GetById(id)

	response := views.JsonResponse{}
	response.Data = res

	return c.RenderJSON(response)
}

func (c Game) Insert() revel.Result {
	var game models.Game
	c.Params.BindJSON(&game)

	response := views.JsonResponse{}
	gs, err := _getGameService()
	if err != nil {
		response.Error = "Internal error"
		return c.RenderJSON(response)
	}
	id, err := gs.InsertGame(game)

	if err != nil {
		response.Error = "Internal error"
		return c.RenderJSON(response)
	}
	response.Data = id

	return c.RenderJSON(response)
}

func _getGameService() (*services.GameService, error) {
	gs, err := services.NewGameService("harvey", "games", "mongodb://localhost:27017")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return gs, err
}