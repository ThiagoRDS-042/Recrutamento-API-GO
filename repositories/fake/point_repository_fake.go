package repositories

import (
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
	"github.com/gofrs/uuid"
)

// DBPoint banco de dados fake de pontos para os testes
var DBPoint = &[]entities.Ponto{}

type pointConnectionFake struct {
	connection        *[]entities.Ponto
	connectionClient  *[]entities.Cliente
	connectionAddress *[]entities.Endereco
}

func (db *pointConnectionFake) CreatePoint(point entities.Ponto) (entities.Ponto, error) {
	pointID, _ := uuid.NewV4()

	point.ID = pointID.String()
	point.DataCriacao = time.Now()
	point.DataAtualizacao = time.Now()

	*db.connection = append(*db.connection, point)

	return point, nil
}

func (db *pointConnectionFake) UpdatePoint(point entities.Ponto) (entities.Ponto, error) {
	point.DataAtualizacao = time.Now()
	point.DataRemocao.Valid = false

	for i, pointValue := range *db.connection {
		if pointValue.ID == point.ID {
			(*db.connection)[i] = point
		}
	}

	return point, nil
}

func (db *pointConnectionFake) FindPointByID(pointID string) entities.Ponto {
	point := entities.Ponto{}

	for _, pointValue := range *db.connection {
		if pointValue.ID == pointID && !pointValue.DataRemocao.Valid {
			point = pointValue
		}
	}

	return point
}

func (db *pointConnectionFake) FindPointByClientIDAndAddressID(clientID string, addressID string) entities.Ponto {
	point := entities.Ponto{}

	for _, pointValue := range *db.connection {
		if pointValue.ClienteID == clientID && pointValue.EnderecoID == addressID {
			point = pointValue
		}
	}

	return point
}

func (db *pointConnectionFake) FindPointsByClientID(clientID string) []entities.Ponto {
	points := []entities.Ponto{}

	for _, pointValue := range *db.connection {
		if pointValue.ClienteID == clientID && !pointValue.DataRemocao.Valid {
			points = append(points, pointValue)
		}
	}

	return points
}

func (db *pointConnectionFake) FindPointsByAddressID(addressID string) []entities.Ponto {
	points := []entities.Ponto{}

	for _, pointValue := range *db.connection {
		if pointValue.EnderecoID == addressID && !pointValue.DataRemocao.Valid {
			points = append(points, pointValue)
		}
	}

	return points
}

func (db *pointConnectionFake) DeletePoint(point entities.Ponto) error {
	for i, pointValue := range *db.connection {
		if pointValue.ID == point.ID {
			(*db.connection)[i].DataRemocao.Scan(time.Now())
		}
	}

	return nil
}

func (db *pointConnectionFake) FindPoints(clientID string, addressID string) []entities.Ponto {
	points := []entities.Ponto{}

	if clientID != "" && addressID != "" {
		for _, pointValue := range *db.connection {
			if pointValue.ClienteID == clientID && pointValue.EnderecoID == addressID &&
				!pointValue.DataRemocao.Valid {
				points = append(points, pointValue)
			}
		}
	} else if clientID != "" {
		for _, pointValue := range *db.connection {
			if pointValue.ClienteID == clientID && !pointValue.DataRemocao.Valid {
				points = append(points, pointValue)
			}
		}
	} else if addressID != "" {
		for _, pointValue := range *db.connection {
			if pointValue.EnderecoID == addressID && !pointValue.DataRemocao.Valid {
				points = append(points, pointValue)
			}
		}
	} else {
		for _, pointValue := range *db.connection {
			if !pointValue.DataRemocao.Valid {
				points = append(points, pointValue)
			}
		}
	}

	if len(points) != 0 {
		for i, point := range points {
			for _, client := range *db.connectionClient {
				if point.ClienteID == client.ID {
					points[i].Cliente = client
				}
			}

			for _, address := range *db.connectionAddress {
				if point.EnderecoID == address.ID {
					points[i].Endereco = address
				}
			}
		}
	}

	return points
}

// NewPointRepositoryFake cria uma nova instancia de PointRepository para os testes.
func NewPointRepositoryFake(database *[]entities.Ponto, connectionClient *[]entities.Cliente, connectionAddress *[]entities.Endereco) repositories.PointRepository {
	return &pointConnectionFake{
		connection:        database,
		connectionClient:  connectionClient,
		connectionAddress: connectionAddress,
	}
}
