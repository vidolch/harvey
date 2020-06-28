package controllers

import (
	"github.com/revel/revel"
	"harvey/app/models/views"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) ApiTest() revel.Result {
	//gs, err := services.NewGameService("harvey", "games", "mongodb://localhost:27017")
	//
	//if err != nil {
	//	log.Fatal(err)
	//	return c.RenderJSON(err)
	//}
	//
	//res := gs.GetAll()

	response := views.JsonResponse{}
	//response.Name = res

	return c.RenderJSON(response)
}
