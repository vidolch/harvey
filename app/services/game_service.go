package services

import (
	"context"
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
	Id string
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
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	cur, err := c.Collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}

	defer cur.Close(ctx)

	var games []models.Game
	for cur.Next(ctx) {
		var result models.Game
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}

		games = append(games, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	ctx.Done()
	return games
}

func (c GameService) GetById(id string) models.Game {
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	filter := bson.M{"id" : id}
	var result models.Game
	err := c.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	ctx.Done()
	return result
}

func (c GameService) InsertGame(game models.Game) (string, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	game.Id = uuid.String()

	_, err = c.Collection.InsertOne(ctx, game)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	id := game.Id
	ctx.Done()
	return id, nil
}
