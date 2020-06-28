package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"harvey/app/models"
	"log"
	"time"
)

type GameService struct {
	MongoClient *mongo.Client
	Collection *mongo.Collection
}

type RawGame struct {
	_id string
	Name string
}

func NewGameService(db string, cl string, url string) (*GameService, error) {
	gs := new(GameService)

	client, err := mongo.NewClient(options.Client().ApplyURI(url))
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Cannot create mongo client.")
	}

	gs.MongoClient = client

	ctx, _ := context.WithTimeout(context.Background(), 10 * time.Second)
	err = gs.MongoClient.Connect(ctx)

	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Cannot connect to mongo client.")
	}

	gs.Collection = gs.MongoClient.Database(db).Collection(cl)

	return gs, nil
}

func (c GameService) GetAll() []models.Game {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := c.Collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(ctx)

	var games []models.Game
	for cur.Next(ctx) {
		var result RawGame
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		var game models.Game
		game.Id = result._id
		log.Fatal(result._id)
		game.Name = result.Name

		games = append(games, game)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return games
}

func (c GameService) InsertGame(game models.Game) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	game.Id = uuid.String()

	res, err := c.Collection.InsertOne(ctx, game)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	id := fmt.Sprintf("%v", res.InsertedID)
	return id, nil
}
