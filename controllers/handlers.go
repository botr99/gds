package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminService interface {
	RegisterStudents(teacher string, students []string) error
	GetCommonStudents(teachers []string) (commonStudents []string, err error)
	SuspendStudent(student string) error
	RetrieveForNotifications(teacher string, notification string) ([]string, error)
}

func NewAdminServer(service AdminService) *AdminServer {
	return &AdminServer{service: service}
}

type AdminServer struct {
	service AdminService
}

type HandleRegisterInput struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

type HandleCommonStudentsOutput struct {
	Students []string `json:"students"`
}

type HandleSuspendInput struct {
	Student string `json:"student"`
}

type HandleRetrieveForNotificationsInput struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}

type HandleRetrieveForNotificationsOutput struct {
	Recipients []string `json:"recipients"`
}

func (a *AdminServer) HandleRegister(c *gin.Context) {
	var HandleRegisterInput HandleRegisterInput;
	if err := c.ShouldBindJSON(&HandleRegisterInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Invalid JSON"})
		return
	}

	err := a.service.RegisterStudents(HandleRegisterInput.Teacher, HandleRegisterInput.Students)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}

func (a *AdminServer) HandleCommonStudents(c *gin.Context) {
	paramPairs := c.Request.URL.Query()
	if len(paramPairs) == 0 {
		c.JSON(http.StatusOK, HandleCommonStudentsOutput{Students: []string{}})
		return
	}
	teachers := make([]string, len(paramPairs))
	for _, values := range paramPairs {
		teachers = append(teachers, values...)
	}

	commonStudents, err := a.service.GetCommonStudents(teachers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, HandleCommonStudentsOutput{Students: commonStudents})
}

func (a *AdminServer) HandleSuspend(c *gin.Context) {
	var HandleSuspendInput HandleSuspendInput;
	if err := c.ShouldBindJSON(&HandleSuspendInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Invalid JSON"})
		return
	}

	err := a.service.SuspendStudent(HandleSuspendInput.Student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}
	
	c.JSON(http.StatusNoContent, gin.H{})
}

func (a *AdminServer) HandleRetrieveForNotifications(c *gin.Context) {
	var HandleRetrieveForNotificationsInput HandleRetrieveForNotificationsInput;
	if err := c.ShouldBindJSON(&HandleRetrieveForNotificationsInput); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Invalid JSON"})
		return
	}

	recipients, err := a.service.RetrieveForNotifications(HandleRetrieveForNotificationsInput.Teacher, HandleRetrieveForNotificationsInput.Notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	c.JSON(http.StatusOK, HandleRetrieveForNotificationsOutput{Recipients: recipients})
}
