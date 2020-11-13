package model

import (
	"context"
	"galasejahtera/pkg/dao"

	"go.mongodb.org/mongo-driver/mongo"
)

// Model ...
type Model struct {
	userDAO     dao.IUserDAO
	authDAO     dao.IAuthDAO
	zoneDAO     dao.IZoneDAO
	activityDAO dao.IActivityDAO
	faqDAO      dao.IFaqDAO
}

// InitModel ...
func InitModel(client *mongo.Client) IModel {
	return &Model{
		userDAO:     dao.InitUserDAO(client),
		authDAO:     dao.InitAuthDAO(client),
		zoneDAO:     dao.InitZoneDAO(client),
		activityDAO: dao.InitActivityDAO(client),
		faqDAO:      dao.InitFaqDAO(client),
	}
}

func (m *Model) InitMongoIndex(ctx context.Context) {
	m.authDAO.InitIndex(ctx)
}
