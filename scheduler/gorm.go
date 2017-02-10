package scheduler

import (
	log "github.com/Sirupsen/logrus"
	"github.com/bewiwi/mta/models"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"time"
	"encoding/json"
)

type DB2 struct {
	db *gorm.DB
}

type dbService struct {
	Id          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Description string
}

type dbCheck struct {
	Id        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
	Type      string     `json:"type"`
	Freq      int        `json:"freq"`
	ServiceId int
	Param     string     `sql:"type:json"`
}

func (d *DB2) Init() {
	db, err := gorm.Open(viper.GetString("SCHEDULER.DB.DRIVER"), viper.GetString("SCHEDULER.DB.DATASOURCE"))
	if err != nil {
		log.WithError(err).Fatal("Can't connect to db")
	}
	d.db = db
	db.AutoMigrate(&dbService{})
	db.AutoMigrate(&dbCheck{})
}

func (d *DB2) Close() {
	d.db.Close()
}

func (d *DB2) serviceToDb(service models.Service) dbService {
	return dbService{
		Id:          uint(service.Id),
		Description: service.Description,
	}
}

func (d *DB2) db2Service(service dbService) models.Service {
	return models.Service{
		Id:          int(service.Id),
		Description: service.Description,
	}
}

func (d *DB2) checkToDb(check models.CheckV1) dbCheck {
	param, err := check.Param.MarshalJSON()
	if err != nil {
		log.WithError(err).Fatal("Can't convert")
	}
	return dbCheck{
		Id:        uint(check.Metadata.Id),
		Type:      check.Metadata.Type,
		Freq:      check.Metadata.Freq,
		ServiceId: check.Metadata.ServiceId,
		Param:     string(param),
	}
}

func (d *DB2) db2Check(check dbCheck) models.CheckV1 {
	param := json.RawMessage{}
	param.UnmarshalJSON([]byte(check.Param))
	return models.CheckV1{
		Metadata: models.CheckMetadaV1{
			Id:        int(check.Id),
			Type:      check.Type,
			Freq:      check.Freq,
			ServiceId: check.ServiceId,
		},
		Param: &param,
	}
}

func (d *DB2) dbs2Checks(checks []dbCheck) []models.CheckV1 {
	models := make([]models.CheckV1, len(checks))
	for i, check := range checks {
		models[i] = d.db2Check(check)
	}
	return models
}

func (d *DB2) CreateService(service models.Service) (models.Service, error) {
	dbService := d.serviceToDb(service)
	err := d.db.Create(&dbService).Error
	return d.db2Service(dbService), err
}

func (d *DB2) GetService(id int) (models.Service, error) {
	service := dbService{Id: uint(id)}
	err := d.db.First(&service).Error
	return d.db2Service(service), err
}

func (d *DB2) DeleteService(id int) error {
	service := dbService{Id: uint(id)}
	err := d.db.Delete(&service).Error
	return err
}

func (d *DB2) UpdateService(id int, service models.Service) (models.Service, error) {
	// Create
	service.Id = id
	dbService := d.serviceToDb(service)
	err := d.db.Save(&dbService).Error
	return d.db2Service(dbService), err
}

func (d *DB2) CreateCheck(check models.CheckV1) (models.CheckV1, error) {
	// Create
	dbCheck := d.checkToDb(check)
	err := d.db.Create(&dbCheck).Error
	return d.db2Check(dbCheck), err
}

func (d *DB2) DeleteCheck(id int) error {
	check := dbCheck{Id: uint(id)}
	err := d.db.Delete(&check).Error
	return err
}

func (d *DB2) GetChecks(serviceId int) ([]models.CheckV1, error) {
	var checks []dbCheck
	err := d.db.Where(dbCheck{ServiceId: serviceId}).Find(&checks).Error
	return d.dbs2Checks(checks), err
}

func (d *DB2) GetAllChecks() ([]models.CheckV1,error) {
	var checks []dbCheck
	err := d.db.Find(&checks).Error
	return d.dbs2Checks(checks), err
}

func (d *DB2) GetCheck(checkId int) (models.CheckV1, error) {
	check := dbCheck{Id: uint(checkId)}
	err := d.db.First(&check).Error
	return d.db2Check(check), err
}

func init() {
	viper.SetDefault("SCHEDULER.DB.DRIVER", "postgres")
	viper.SetDefault("SCHEDULER.DB.DATASOURCE", "user=mta dbname=mta")
}
