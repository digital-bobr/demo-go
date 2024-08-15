package data

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"service-status-reporter/pkg/model"
)

var instance *gorm.DB

func init() {
	var err error
	instance, err = gorm.Open(sqlite.Open("database/service.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database:\n" + err.Error())
	}
}

// service

func UpdateService(service *model.Service) error {
	if err := instance.Save(&service).Error; err != nil {
		log.Printf("failed to update service with ID: %d\n%v", service.ID, err)
		return err
	}
	return nil
}

func CreateService(service *model.Service) error {
	if err := instance.Create(&service).Error; err != nil {
		log.Printf("failed to create new service: \n%v", err)
		return err
	}
	return nil
}

func DeleteService(id int) error {
	if err := instance.Delete(&model.Service{}, id).Error; err != nil {
		log.Printf("failed to delete service: \n%v", err)
		return err
	}
	return nil
}

func GetServiceById(id int) model.Service {
	var service model.Service
	if err := instance.First(&service, id).Error; err != nil {
		log.Printf("failed to fetch service with ID: %d\n%v", id, err)
	}
	return service
}

func GetServicesByContactId(id int) []model.Service {
	var services []model.Service
	if err := instance.Where("ste_owner_id = ?", id).Or("developer_id", id).Find(&services); err != nil {
		log.Printf("failed to fetch any services with contact ID: %d\n%v", id, err)
	}
	return services
}

func GetServicesByTeamId(id int) []model.Service {
	var services []model.Service
	if err := instance.Where("team_id = ?", id).Find(&services); err != nil {
		log.Printf("failed to fetch any services with team ID: %d\n%v", id, err)
	}
	return services
}

func GetAllServices() []model.Service {
	var services []model.Service
	if err := instance.Preload("SteOwner").Preload("Developer").Preload("Team").Find(&services); err != nil {
		log.Printf("failed to fetch services: %v", err)
	}
	return services
}

// contact

func UpdateContact(contact *model.Contact) error {
	if err := instance.Save(&contact).Error; err != nil {
		log.Printf("failed to update contact with ID: %d\n%v", contact.ID, err)
		return err
	}
	return nil
}

func CreateContact(contact *model.Contact) error {
	if err := instance.Create(&contact).Error; err != nil {
		log.Printf("failed to create new contact:\n%v", err)
		return err
	}
	return nil
}

func DeleteContact(id int) error {
	services := GetServicesByContactId(id)
	if len(services) != 0 {
		return fmt.Errorf("contact with ID '%d' cannot be deleted due to being is use", id)
	}

	if err := instance.Delete(&model.Contact{}, id).Error; err != nil {
		log.Printf("failed to delete contact: \n%v", err)
		return err
	}
	return nil
}

func GetContactById(id int) model.Contact {
	var contact model.Contact
	if err := instance.First(&contact, id).Error; err != nil {
		log.Printf("failed to fetch contact with ID: %d\n%v", id, err)
	}
	return contact
}

func GetAllContacts() []model.Contact {
	var contacts []model.Contact
	if err := instance.Preload("Role").Find(&contacts); err != nil {
		log.Printf("failed to fetch contacts: %v", err)
	}
	return contacts
}

func GetAllSTEs() []model.Contact {
	var stes []model.Contact
	if err := instance.Joins("Role").Where("Role.name = ?", "ste").Find(&stes).Error; err != nil {
		log.Printf("failed to fetch ste contacts: %v", err)
	}
	return stes
}

func GetAllDevelopers() []model.Contact {
	var developers []model.Contact
	if err := instance.Joins("Role").Where("Role.name = ?", "developer").Find(&developers).Error; err != nil {
		log.Printf("failed to fetch developer contacts: %v", err)
	}
	return developers
}

// team

func CreateTeam(team *model.Team) error {
	if err := instance.Create(&team).Error; err != nil {
		log.Printf("failed to create team:\n%v", err)
		return err
	}
	return nil
}

func DeleteTeam(id int) error {
	services := GetServicesByTeamId(id)
	if len(services) != 0 {
		return fmt.Errorf("team with ID '%d' cannot be deleted due to being is use", id)
	}
	if err := instance.Delete(&model.Team{}, id).Error; err != nil {
		log.Printf("failed to delete team: \n%v", err)
		return err
	}
	return nil
}

func GetAllTeams() []model.Team {
	var teams []model.Team
	if err := instance.Find(&teams).Error; err != nil {
		log.Printf("failed to fetch teams: %v", err)
	}
	return teams
}

// role

func GetAllRoles() []model.Role {
	var roles []model.Role
	if err := instance.Find(&roles).Error; err != nil {
		log.Printf("failed to fetch roles: %v", err)
	}
	return roles
}

func GetRoleById(id int) model.Role {
	var role model.Role
	if err := instance.First(&role, id).Error; err != nil {
		log.Printf("failed to fetch role with ID: %d\n%v", id, err)
	}
	return role
}
