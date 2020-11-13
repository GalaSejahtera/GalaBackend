package model

import (
	"galasejahtera/pkg/dao"

	"go.mongodb.org/mongo-driver/mongo"
)

// Model ...
type Model struct {
	userDAO dao.IUserDAO
	authDAO dao.IAuthDAO
	faqDAO  dao.IFaqDAO
}

// InitModel ...
func InitModel(client *mongo.Client) IModel {
	return &Model{
		userDAO: dao.InitUserDAO(client),
		authDAO: dao.InitAuthDAO(client),
		faqDAO:  dao.InitFaqDAO(client),
	}
}
