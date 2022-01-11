package repositories

import (
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
)

// DBAddress banco de dados fake de endereços para os testes
var DBAddress = &[]entities.Endereco{}

type addressConnectionFake struct {
	connection *[]entities.Endereco
}

func (db *addressConnectionFake) CreateAddress(address entities.Endereco) (entities.Endereco, error) {
	address.DataCriacao = time.Now()
	address.DataAtualizacao = time.Now()

	*db.connection = append(*db.connection, address)

	return address, nil
}

func (db *addressConnectionFake) UpdateAddress(address entities.Endereco) (entities.Endereco, error) {
	address.DataAtualizacao = time.Now()

	for i, v := range *db.connection {
		if v.ID == address.ID {
			(*db.connection)[i] = address
		}
	}

	return address, nil
}

func (db *addressConnectionFake) FindAddressByID(addressID string) entities.Endereco {
	address := entities.Endereco{}

	for _, v := range *db.connection {
		if v.ID == addressID && !v.DataRemocao.Valid {
			address = v
		}
	}

	return address
}

func (db *addressConnectionFake) FindAddressByFields(street string, neighborhood string, number int) entities.Endereco {
	address := entities.Endereco{}

	for _, v := range *db.connection {
		if v.Logradouro == street && v.Bairro == neighborhood && v.Numero == number {
			address = v
		}
	}

	return address
}

func (db *addressConnectionFake) DeleteAddress(address entities.Endereco) error {
	for i, v := range *db.connection {
		if v.ID == address.ID {
			(*db.connection)[i].DataRemocao.Scan(time.Now())
		}
	}

	return nil
}

func (db *addressConnectionFake) FindAddresses(street string, neighborhood string, number string) []entities.Endereco {
	return *db.connection
}

// NewAddressRepositoryFake cria uma nova instancia de AddressRepository para os testes.
func NewAddressRepositoryFake(database *[]entities.Endereco) repositories.AddressRepository {
	return &addressConnectionFake{
		connection: database,
	}
}
