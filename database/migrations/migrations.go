package migrations

import (
	"github.com/ThiagoRDS-042/Recrutamento-API-GO/entities"
	"gorm.io/gorm"
)

// RunMigrations cria as tabelas do banco de dados baseado nas entidades definidas na pasta entities.
func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(
		entities.Cliente{},
		entities.Endereco{},
		entities.Ponto{},
		entities.Contrato{},
		entities.ContratoEvento{},
	)
}
