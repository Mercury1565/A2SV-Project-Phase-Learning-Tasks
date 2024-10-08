package bootstrap

import "go.mongodb.org/mongo-driver/mongo"

type Application struct {
	Env   *Env
	Mongo *mongo.Client
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = NewMongoDBClient(app.Env)
	return *app
}

func (app *Application) CloseMongoDBConnection() {
	CloseMongoDBClient(app.Mongo)
}
