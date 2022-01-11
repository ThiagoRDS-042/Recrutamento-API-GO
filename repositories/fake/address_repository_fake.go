package repositories

import (
	"strconv"
	"strings"
	"time"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	repositories "github.com/ThiagoRDS-042/Recrutamento-API-GO/repositories/postgres"
	"github.com/gofrs/uuid"
)

// DBAddress banco de dados fake de endereços para os testes
var DBAddress = &[]entities.Endereco{}

type addressConnectionFake struct {
	connection *[]entities.Endereco
}

func (db *addressConnectionFake) CreateAddress(address entities.Endereco) (entities.Endereco, error) {
	addressID, _ := uuid.NewV4()

	address.ID = addressID.String()
	address.DataCriacao = time.Now()
	address.DataAtualizacao = time.Now()

	*db.connection = append(*db.connection, address)

	return address, nil
}

func (db *addressConnectionFake) UpdateAddress(address entities.Endereco) (entities.Endereco, error) {
	address.DataAtualizacao = time.Now()
	address.DataRemocao.Valid = false

	for i, addressValue := range *db.connection {
		if addressValue.ID == address.ID {
			(*db.connection)[i] = address
		}
	}

	return address, nil
}

func (db *addressConnectionFake) FindAddressByID(addressID string) entities.Endereco {
	address := entities.Endereco{}

	for _, addressValue := range *db.connection {
		if addressValue.ID == addressID && !addressValue.DataRemocao.Valid {
			address = addressValue
		}
	}

	return address
}

func (db *addressConnectionFake) FindAddressByFields(street string, neighborhood string, number int) entities.Endereco {
	address := entities.Endereco{}

	for _, addressValue := range *db.connection {
		if addressValue.Logradouro == street && addressValue.Bairro == neighborhood && addressValue.Numero == number {
			address = addressValue
		}
	}

	return address
}

func (db *addressConnectionFake) DeleteAddress(address entities.Endereco) error {
	for i, addressValue := range *db.connection {
		if addressValue.ID == address.ID {
			(*db.connection)[i].DataRemocao.Scan(time.Now())
		}
	}

	return nil
}

func (db *addressConnectionFake) FindAddresses(street string, neighborhood string, number string) []entities.Endereco {
	address := []entities.Endereco{}

	if street != "" && neighborhood != "" && number != "" {
		numberConverted, _ := strconv.Atoi(number)

		for _, addressValue := range *db.connection {
			if strings.Contains(addressValue.Logradouro, street) && strings.Contains(addressValue.Bairro, neighborhood) &&
				addressValue.Numero == numberConverted && !addressValue.DataRemocao.Valid {
				address = append(address, addressValue)
			}
		}

	} else if street != "" && neighborhood != "" {

		for _, addressValue := range *db.connection {
			if strings.Contains(addressValue.Logradouro, street) && strings.Contains(addressValue.Bairro, neighborhood) &&
				!addressValue.DataRemocao.Valid {
				address = append(address, addressValue)
			}
		}

	} else if street != "" && number != "" {
		numberConverted, _ := strconv.Atoi(number)

		for _, addressValue := range *db.connection {
			if strings.Contains(addressValue.Logradouro, street) && addressValue.Numero == numberConverted &&
				!addressValue.DataRemocao.Valid {
				address = append(address, addressValue)
			}
		}

	} else if neighborhood != "" && number != "" {
		numberConverted, _ := strconv.Atoi(number)

		for _, addressValue := range *db.connection {
			if strings.Contains(addressValue.Bairro, neighborhood) && addressValue.Numero == numberConverted &&
				!addressValue.DataRemocao.Valid {
				address = append(address, addressValue)
			}
		}

	} else if street != "" {
		for _, addressValue := range *db.connection {
			if strings.Contains(addressValue.Logradouro, street) && !addressValue.DataRemocao.Valid {
				address = append(address, addressValue)
			}
		}
	} else if neighborhood != "" {
		for _, addressValue := range *db.connection {
			if strings.Contains(addressValue.Bairro, neighborhood) && !addressValue.DataRemocao.Valid {
				address = append(address, addressValue)
			}
		}
	} else if number != "" {
		numberConverted, _ := strconv.Atoi(number)

		for _, addressValue := range *db.connection {
			if addressValue.Numero == numberConverted && !addressValue.DataRemocao.Valid {
				address = append(address, addressValue)
			}
		}
	} else {
		for _, addressValue := range *db.connection {
			if !addressValue.DataRemocao.Valid {
				address = append(address, addressValue)
			}
		}
	}

	return address
}

// NewAddressRepositoryFake cria uma nova instancia de AddressRepository para os testes.
func NewAddressRepositoryFake(database *[]entities.Endereco) repositories.AddressRepository {
	return &addressConnectionFake{
		connection: database,
	}
}
