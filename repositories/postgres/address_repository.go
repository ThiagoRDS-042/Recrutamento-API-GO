package repositories

import (
	"fmt"
	"log"

	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"gorm.io/gorm"
)

// AddressRepository representa o contracto de AddressRepository.
type AddressRepository interface {
	CreateAddress(address entities.Endereco) (entities.Endereco, error)
	UpdateAddress(address entities.Endereco) (entities.Endereco, error)
	FindAddressByID(addressID string) entities.Endereco
	FindAddressByFields(street string, neighborhood string, number int) entities.Endereco
	DeleteAddress(address entities.Endereco) error
	FindAddresses(street string, neighborhood string, number string) []entities.Endereco
}

type addressConnection struct {
	connection *gorm.DB
}

func (db *addressConnection) CreateAddress(address entities.Endereco) (entities.Endereco, error) {
	err := db.connection.Create(&address).Error
	if err != nil {
		return address, err
	}

	return address, nil
}

func (db *addressConnection) UpdateAddress(address entities.Endereco) (entities.Endereco, error) {
	err := db.connection.Save(&address).Error
	if err != nil {
		return address, err
	}

	return address, nil
}

func (db *addressConnection) FindAddressByID(addressID string) entities.Endereco {
	address := entities.Endereco{}

	err := db.connection.First(&address, "id = ?", addressID).Error
	if err != nil {
		log.Println(err.Error())
	}

	return address
}

func (db *addressConnection) FindAddressByFields(street string, neighborhood string, number int) entities.Endereco {
	address := entities.Endereco{}

	err := db.connection.Unscoped().First(&address, "logradouro = ? AND bairro = ? AND numero = ?",
		street, neighborhood, number).Error
	if err != nil {
		log.Println(err.Error())
	}

	return address
}

func (db *addressConnection) DeleteAddress(address entities.Endereco) error {
	err := db.connection.Delete(&address).Error
	if err != nil {
		return err
	}

	return nil
}

func (db *addressConnection) FindAddresses(street string, neighborhood string, number string) []entities.Endereco {
	addresses := []entities.Endereco{}

	var sqlQuery string

	if street != "" {
		street = fmt.Sprint("%", street, "%")
		sqlQuery += fmt.Sprintf("logradouro LIKE '%v' ", street)
	} else {
		sqlQuery += "NOT logradouro IS NULL "
	}

	if neighborhood != "" {
		neighborhood = fmt.Sprint("%", neighborhood, "%")
		sqlQuery += fmt.Sprintf("AND bairro LIKE '%v' ", neighborhood)
	} else {
		sqlQuery += "AND NOT bairro IS NULL "
	}

	if number != "" {
		sqlQuery += fmt.Sprintf("AND numero = '%v'", number)
	} else {
		sqlQuery += "AND NOT numero IS NULL"
	}

	err := db.connection.Find(&addresses, sqlQuery).Error
	if err != nil {
		log.Println(err.Error())
	}

	return addresses

}

// NewAddressRepository cria uma nova instancia de AddressRepository.
func NewAddressRepository(database *gorm.DB) AddressRepository {
	return &addressConnection{
		connection: database,
	}
}
