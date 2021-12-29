package repositories

import (
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/database"
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"gorm.io/gorm"
)

// PointRepository representa o contracto de PointRepository.
type PointRepository interface {
	CreatePoint(point entities.Ponto) (entities.Ponto, error)
	UpdatePoint(point entities.Ponto) (entities.Ponto, error)
	FindPointByID(pointID string) entities.Ponto
	FindPointByClientIDAndAddressID(clientID string, addressID string) entities.Ponto
	FindPointsByClientID(clientID string) []entities.Ponto
	FindPointsByAddressID(addressID string) []entities.Ponto
	DeletePoint(point entities.Ponto) error
	FindPoints(clientID string, addressID string) []entities.Ponto
}

type pointConnection struct {
	connection *gorm.DB
}

func (db *pointConnection) CreatePoint(point entities.Ponto) (entities.Ponto, error) {
	err := db.connection.Create(&point).Error
	if err != nil {
		return point, err
	}

	return point, nil
}

func (db *pointConnection) UpdatePoint(point entities.Ponto) (entities.Ponto, error) {
	err := db.connection.Save(&point).Error
	if err != nil {
		return point, err
	}

	return point, nil
}

func (db *pointConnection) FindPointByID(pointID string) entities.Ponto {
	point := entities.Ponto{}

	err := db.connection.First(&point, "id = ?", pointID).Error
	if err != nil {
		log.Println(err.Error())
	}

	return point
}

func (db *pointConnection) FindPointByClientIDAndAddressID(clientID string, addressID string) entities.Ponto {
	point := entities.Ponto{}

	err := db.connection.Unscoped().First(&point, "cliente_id = ? AND endereco_id = ?", clientID, addressID).Error
	if err != nil {
		log.Println(err.Error())
	}

	return point
}

func (db *pointConnection) FindPointsByClientID(clientID string) []entities.Ponto {
	points := []entities.Ponto{}

	err := db.connection.Find(&points, "cliente_id = ?", clientID).Error
	if err != nil {
		log.Println(err.Error())
	}

	return points
}

func (db *pointConnection) FindPointsByAddressID(addressID string) []entities.Ponto {
	points := []entities.Ponto{}

	err := db.connection.Find(&points, "endereco_id = ?", addressID).Error
	if err != nil {
		log.Println(err.Error())
	}

	return points
}

func (db *pointConnection) DeletePoint(point entities.Ponto) error {
	err := db.connection.Delete(&point).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *pointConnection) FindPoints(clientID string, addressID string) []entities.Ponto {
	points := []entities.Ponto{}

	var err error

	switch {
	case clientID != "" && addressID != "":
		err = db.connection.Preload("Cliente").Preload("Endereco").
			Find(&points, "cliente_id = ? AND endereco_id = ?", clientID, addressID).Error

	case clientID != "":
		err = db.connection.Preload("Cliente").Preload("Endereco").
			Find(&points, "cliente_id = ?", clientID).Error

	case addressID != "":
		err = db.connection.Preload("Cliente").Preload("Endereco").
			Find(&points, "endereco_id = ?", addressID).Error

	default:
		err = db.connection.Preload("Cliente").Preload("Endereco").
			Find(&points).Error
	}

	if err != nil {
		log.Println(err.Error())
	}

	return points
}

// NewPointRepository cria uma nova instancia de PointRepository.
func NewPointRepository() PointRepository {
	return &pointConnection{
		connection: database.GetDB(),
	}
}
