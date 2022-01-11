package repositories

import (
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
)

// DBPoint banco de dados fake de pontos para os testes
var DBPoint = &[]entities.Ponto{}

type pointConnectionFake struct {
	connection *[]entities.Ponto
}

func (db *pointConnectionFake) CreatePoint(point entities.Ponto) (entities.Ponto, error) {
	point.DataCriacao = time.Now()
	point.DataAtualizacao = time.Now()

	*db.connection = append(*db.connection, point)

	return point, nil
}

func (db *pointConnectionFake) UpdatePoint(point entities.Ponto) (entities.Ponto, error) {
	point.DataAtualizacao = time.Now()

	for i, v := range *db.connection {
		if v.ID == point.ID {
			(*db.connection)[i] = point
		}
	}

	return point, nil
}

func (db *pointConnectionFake) FindPointByID(pointID string) entities.Ponto {
	point := entities.Ponto{}

	for _, v := range *db.connection {
		if v.ID == pointID && !v.DataRemocao.Valid {
			point = v
		}
	}

	return point
}

func (db *pointConnectionFake) FindPointByClientIDAndAddressID(clientID string, addressID string) entities.Ponto {
	point := entities.Ponto{}

	for _, v := range *db.connection {
		if v.ClienteID == clientID && v.EnderecoID == addressID {
			point = v
		}
	}

	return point
}

func (db *pointConnectionFake) FindPointsByClientID(clientID string) []entities.Ponto {
	points := []entities.Ponto{}

	for _, v := range *db.connection {
		if v.ClienteID == clientID && !v.DataRemocao.Valid {
			points = append(points, v)
		}
	}

	return points
}

func (db *pointConnectionFake) FindPointsByAddressID(addressID string) []entities.Ponto {
	points := []entities.Ponto{}

	for _, v := range *db.connection {
		if v.EnderecoID == addressID && !v.DataRemocao.Valid {
			points = append(points, v)
		}
	}

	return points
}

func (db *pointConnectionFake) DeletePoint(point entities.Ponto) error {
	for i, v := range *db.connection {
		if v.ID == point.ID {
			(*db.connection)[i].DataRemocao.Scan(time.Now())
		}
	}

	return nil
}

func (db *pointConnectionFake) FindPoints(clientID string, addressID string) []entities.Ponto {
	points := []entities.Ponto{}

	if clientID != "" && addressID != "" {
		for _, v := range *db.connection {
			if v.ClienteID == clientID && v.EnderecoID == addressID && !v.DataRemocao.Valid {
				points = append(points, v)
			}
		}
	} else if clientID != "" {
		for _, v := range *db.connection {
			if v.ClienteID == clientID && !v.DataRemocao.Valid {
				points = append(points, v)
			}
		}
	} else if addressID != "" {
		for _, v := range *db.connection {
			if v.EnderecoID == addressID && !v.DataRemocao.Valid {
				points = append(points, v)
			}
		}
	} else {
		points = append(points, *db.connection...)
	}

	return points
	// err := db.connection.Preload("Cliente").Preload("Endereco").Find(&points, sqlQuery).Error
}

// NewPointRepositoryFake cria uma nova instancia de PointRepository para os testes.
func NewPointRepositoryFake(database *[]entities.Ponto) repositories.PointRepository {
	return &pointConnectionFake{
		connection: database,
	}
}
