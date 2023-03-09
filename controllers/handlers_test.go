package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockAdminService struct {
	registerMap map[string][]string
	commonStudents []string
	suspendedStudent string
	recipients []string
}

func (m *MockAdminService) RegisterStudents(teacher string, students []string) error {
	m.registerMap[teacher] = append(m.registerMap[teacher], students...)
	return nil
}

func TestRegisterStudents(t *testing.T) {
	t.Run("can register students", func(t *testing.T) {
		teacher := "teacher1@gmail.com"
		students := []string{"student1@gmail.com", "student2@gmail.com", "student3@gmail.com"}
		input := &HandleRegisterInput{ Teacher: teacher, Students: students }
		jsonBody, _ := json.Marshal(input)

		mockAdminService := &MockAdminService{
			registerMap: make(map[string][]string),
		}
		adminServer := NewAdminServer(mockAdminService)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/register", strings.NewReader(string(jsonBody)))
		
		adminServer.HandleRegister(c)

		if c.Writer.Status() != http.StatusNoContent {
			t.Errorf("expected status code %d, got %d", http.StatusNoContent, c.Writer.Status())
		}
		if len(mockAdminService.registerMap[teacher]) != len(students) {
			t.Errorf("expected %d students to be registered, got %d", len(students), len(mockAdminService.registerMap))
		}
	})

	t.Run("returns 422 when invalid json is sent", func(t *testing.T) {
		mockAdminService := &MockAdminService{
			registerMap: make(map[string][]string),
		}
		adminServer := NewAdminServer(mockAdminService)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/register", strings.NewReader(`{teacher: "teacher1@gmail.com"}`))
		
		adminServer.HandleRegister(c)

		if c.Writer.Status() != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, c.Writer.Status())
		}
	})
}

func (m *MockAdminService) GetCommonStudents(teachers []string) (commonStudents []string, err error) {
	commonStudents = m.commonStudents;
	return
}

func TestCommonStudents(t *testing.T) {
	t.Run("returns empty list when no query parameters is passed", func(t *testing.T) {
		mockAdminService := &MockAdminService{
			commonStudents: []string{},
		}
		adminServer := NewAdminServer(mockAdminService)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/commonstudents", nil)
		
		adminServer.HandleCommonStudents(c)

		if c.Writer.Status() != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, c.Writer.Status())
		}		
	})
}

func (m *MockAdminService) SuspendStudent(student string) error {
	m.suspendedStudent = student
	return nil
}

func TestSuspendStudent(t *testing.T) {
	t.Run("can suspend student", func(t *testing.T) {
		student := "student1@gmail.com"
		input := &HandleSuspendInput{ Student: student }
		jsonBody, _ := json.Marshal(input)

		mockAdminService := &MockAdminService{}
		adminServer := NewAdminServer(mockAdminService)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/suspend", strings.NewReader(string(jsonBody)))

		adminServer.HandleSuspend(c)

		if c.Writer.Status() != http.StatusNoContent {
			t.Errorf("expected status code %d, got %d", http.StatusNoContent, c.Writer.Status())
		}

		if mockAdminService.suspendedStudent != student {
			t.Errorf("expected student %s to be suspended, got %s", student, mockAdminService.suspendedStudent)
		}
	})

	t.Run("returns 422 when invalid json is sent", func(t *testing.T) {
		mockAdminService := &MockAdminService{}
		adminServer := NewAdminServer(mockAdminService)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/suspend", strings.NewReader(`{students: ["student1@gmail.com", "student2@gmail.com"]}`))

		adminServer.HandleSuspend(c)

		if c.Writer.Status() != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, c.Writer.Status())
		}
	})
}

func (m *MockAdminService) RetrieveForNotifications(teacher string, notification string) (recipients []string, err error) {
	recipients = m.recipients
	return
}

func TestRetrieveForNotifications(t *testing.T) {
	t.Run("returns 422 when invalid json is sent", func(t *testing.T) {
		mockAdminService := &MockAdminService{}
		adminServer := NewAdminServer(mockAdminService)

		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request, _ = http.NewRequest(http.MethodPost, "/api/suspend", strings.NewReader(`{students: "student1@gmail.com"}`))

		adminServer.HandleSuspend(c)

		if c.Writer.Status() != http.StatusUnprocessableEntity {
			t.Errorf("expected status code %d, got %d", http.StatusUnprocessableEntity, c.Writer.Status())
		}
	})
}
