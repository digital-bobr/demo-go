package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"service-status-reporter/pkg/data"
	"service-status-reporter/pkg/model"
	"strconv"
)

var router *gin.Engine

func init() {
	router = gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/static", "./static")
	registerViews()
}

func GetInstance() *gin.Engine {
	return router
}

func registerViews() {
	getIndex()
	getEditServiceForm()
	updateService()
	getNewServiceForm()
	getNewContactForm()
	getNewTeamForm()
	addNewService()
	deleteService()
	addNewContact()
	addNewTeam()
	getBlanc()
	getAllContacts()
	deleteContact()
	getAllTeams()
	deleteTeam()
	getEditContactForm()
	updateContact()
	getServiceHealth()
}

// service

func getBlanc() {
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/index")
	})
}

func getIndex() {
	router.GET("/index", func(c *gin.Context) {
		services := data.GetAllServices()
		if len(services) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load services"})
			return
		}
		c.HTML(http.StatusOK, "index.html", gin.H{
			"services": services,
		})
	})
}

func addNewService() {
	router.POST("/service/create", func(c *gin.Context) {
		service := model.Service{}
		populateServiceDetails(&service, c)

		if err := data.CreateService(&service); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create service", "details": err.Error()})
			return
		}

		c.Redirect(http.StatusSeeOther, "/index")
	})
}

func getEditServiceForm() {
	router.GET("/service/edit", func(c *gin.Context) {
		service := validateServiceById(c, "query")
		c.HTML(http.StatusOK, "service.html", gin.H{
			"title":      "Service Status Reporter - Edit Service",
			"buttonText": "Save",
			"endpoint":   "/service/update",
			"service":    service,
			"steOwners":  data.GetAllSTEs(),
			"developers": data.GetAllDevelopers(),
			"teams":      data.GetAllTeams(),
		})
	})
}

func getNewServiceForm() {
	router.GET("/service/add", func(c *gin.Context) {

		c.HTML(http.StatusOK, "service.html", gin.H{
			"title":      "Service Status Reporter - Add New Service",
			"buttonText": "Add",
			"endpoint":   "/service/create",
			"steOwners":  data.GetAllSTEs(),
			"developers": data.GetAllDevelopers(),
			"teams":      data.GetAllTeams(),
		})
	})
}

func updateService() {
	router.POST("/service/update", func(c *gin.Context) {
		idStr := c.PostForm("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID", "details": err.Error()})
			return
		}

		service := data.GetServiceById(id)
		if service.ID != id {
			c.JSON(http.StatusNotFound, gin.H{"error": "Service with id '" + idStr + "' not found"})
			return
		}
		populateServiceDetails(&service, c)

		if err := data.UpdateService(&service); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update service", "details": err.Error()})
			return
		}

		c.Redirect(http.StatusSeeOther, "/index")
	})
}

func deleteService() {
	router.DELETE("/service/delete/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service ID", "details": err.Error()})
			return
		}
		if err := data.DeleteService(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete service", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully deleted",
		})
	})
}

func getServiceHealth() {
	router.GET("/service/health/:id", func(c *gin.Context) {
		service := validateServiceById(c, "param")
		serviceHealth := getServiceHealthInfo(service)
		c.JSON(http.StatusOK, serviceHealth)
	})
}

// contact

func getAllContacts() {
	router.GET("/contact/all", func(c *gin.Context) {
		contacts := data.GetAllContacts()
		if len(contacts) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contacts"})
			return
		}
		c.HTML(http.StatusOK, "contacts.html", gin.H{
			"contacts": contacts,
		})
	})
}

func getNewContactForm() {
	router.GET("/contact/add", func(c *gin.Context) {

		c.HTML(http.StatusOK, "contact.html", gin.H{
			"title":      "Service Status Reporter - Add New Contact",
			"buttonText": "Add",
			"endpoint":   "/contact/create",
			"roles":      data.GetAllRoles(),
		})
	})
}

func addNewContact() {
	router.POST("/contact/create", func(c *gin.Context) {
		contact := model.Contact{}
		populateContactDetails(&contact, c)
		if err := data.CreateContact(&contact); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact", "details": err.Error()})
			return
		}
		if c.PostForm("is_form") == "true" {
			c.Redirect(http.StatusSeeOther, "/index")
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"id":      contact.ID,
			})
		}
	})
}

func deleteContact() {
	router.DELETE("/contact/delete/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID", "details": err.Error()})
			return
		}
		if err := data.DeleteContact(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully deleted",
		})
	})
}

func getEditContactForm() {
	router.GET("/contact/edit", func(c *gin.Context) {
		idStr := c.Query("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			c.HTML(http.StatusBadRequest, "contact.html", gin.H{"error": "Invalid contact ID"})
			return
		}

		var contact *model.Contact
		contacts := data.GetAllContacts()
		for _, c := range contacts {
			if c.ID == id {
				contact = &c
				break
			}
		}

		if contact == nil {
			c.HTML(http.StatusNotFound, "contact.html", gin.H{"error": "Contact not found"})
			return
		}

		c.HTML(http.StatusOK, "contact.html", gin.H{
			"title":      "Service Status Reporter - Edit Contact",
			"buttonText": "Save",
			"endpoint":   "/contact/update",
			"contact":    contact,
			"roles":      data.GetAllRoles(),
		})
	})
}

func updateContact() {
	router.POST("/contact/update", func(c *gin.Context) {
		idStr := c.PostForm("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contact ID", "details": err.Error()})
			return
		}

		contact := data.GetContactById(id)
		if contact.ID != id {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contact with id '" + idStr + "' not found"})
			return
		}
		populateContactDetails(&contact, c)

		if err := data.UpdateContact(&contact); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact", "details": err.Error()})
			return
		}

		c.Redirect(http.StatusSeeOther, "/contact/all")
	})
}

// team

func getAllTeams() {
	router.GET("/team/all", func(c *gin.Context) {
		teams := data.GetAllTeams()
		if len(teams) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load teams"})
			return
		}
		c.HTML(http.StatusOK, "teams.html", gin.H{
			"teams": teams,
		})
	})
}

func getNewTeamForm() {
	router.GET("/team/add", func(c *gin.Context) {

		c.HTML(http.StatusOK, "team.html", gin.H{
			"title":      "Service Status Reporter - Add New Team",
			"buttonText": "Add",
			"endpoint":   "/team/create",
		})
	})
}

func addNewTeam() {
	router.POST("/team/create", func(c *gin.Context) {
		team := model.Team{}
		team.Name = c.PostForm("name")
		if err := data.CreateTeam(&team); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team", "details": err.Error()})
			return
		}
		if c.PostForm("is_form") == "true" {
			c.Redirect(http.StatusSeeOther, "/index")
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"id":      team.ID,
			})
		}
	})
}

func deleteTeam() {
	router.DELETE("/team/delete/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID", "details": err.Error()})
			return
		}
		if err := data.DeleteTeam(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team", "details": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Successfully deleted",
		})
	})
}

// private usage

func validateServiceById(c *gin.Context, idSource string) model.Service {
	var idStr string
	switch idSource {
	case "param":
		idStr = c.Param("id")
	case "form":
		idStr = c.PostForm("id")
	case "query":
		idStr = c.Query("id")
	default:
		idStr = ""
	}

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.HTML(http.StatusBadRequest, "service.html", gin.H{"error": "Invalid service ID"})
		return model.Service{}
	}

	service := data.GetServiceById(id)
	if service.ID != id {
		c.HTML(http.StatusNotFound, "service.html", gin.H{"error": "Service not found"})
		return model.Service{}
	}
	return service
}

func populateServiceDetails(service *model.Service, c *gin.Context) {
	service.Name = c.PostForm("name")
	teamID, err1 := strconv.Atoi(c.PostForm("team_id"))
	if err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID", "details": err1.Error()})
		return
	}
	service.TeamID = &teamID
	steOwnerId, err2 := strconv.Atoi(c.PostForm("ste_owner_id"))
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ste owner ID", "details": err2.Error()})
		return
	}
	service.SteOwnerID = &steOwnerId
	developerId, err3 := strconv.Atoi(c.PostForm("developer_id"))
	if err3 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid developer ID", "details": err3.Error()})
		return
	}
	service.DeveloperID = &developerId
	service.ProdEnv = c.PostForm("prod_env")
	service.PreprodEnv = c.PostForm("preprod_env")
	service.NonprodEnv = c.PostForm("nonprod_env")
}

func populateContactDetails(contact *model.Contact, c *gin.Context) {
	roleId, err := strconv.Atoi(c.PostForm("role_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID", "details": err.Error()})
		return
	}
	contactRole := data.GetRoleById(roleId)
	contact.Name = c.PostForm("name")
	contact.Email = c.PostForm("email")
	contact.Slack = c.PostForm("slack")
	contact.RoleID = &contactRole.ID
}
