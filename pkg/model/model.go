package model

import (
	"galasejahtera/pkg/dao"

	"go.mongodb.org/mongo-driver/mongo"
)

// Model ...
type Model struct {
	userDAO   dao.IUserDAO
	authDAO   dao.IAuthDAO
	reportDAO dao.IReportDAO
	covidDAO  dao.ICovidDAO
	dailyDAO  dao.IDailyDAO
}

// InitModel ...
func InitModel(client *mongo.Client) IModel {
	return &Model{
		userDAO:   dao.InitUserDAO(client),
		authDAO:   dao.InitAuthDAO(client),
		reportDAO: dao.InitReportDAO(client),
		covidDAO:  dao.InitCovidDAO(client),
		dailyDAO:  dao.InitDailyDAO(client),
	}
}
