package db

import (
	"database/sql"
	"os"
	"regexp"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type DbAdminService struct {
	db *sql.DB
}

func NewDbAdminService() *DbAdminService {
	return &DbAdminService{
		db: Init(),
	}
}

func (d *DbAdminService) RegisterStudents(teacher string, students []string) error {
	var teacherId int64;
	err := d.db.QueryRow("SELECT id FROM teachers WHERE email = ?", teacher).Scan(&teacherId)
	if err != nil {
		return err
	}

	for _, student := range students {
		var studentId int64
		err := d.db.QueryRow("SELECT id FROM students WHERE email = ?", student).Scan(&studentId)
		if err != nil {
			return err
		}
		_, err = d.db.Exec("INSERT INTO teachers_students (teacher_id, student_id) VALUES (?, ?)", teacherId, studentId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DbAdminService) GetCommonStudents(teachers []string) (commonStudents []string, err error) {
	args := make([]interface{}, len(teachers))
	for i, v := range teachers {
		args[i] = v
	}
	var rows *sql.Rows
	rows, err = d.db.Query(`SELECT s.email 
		FROM teachers_students ts, students s, teachers t 
		WHERE ts.student_id = s.id AND ts.teacher_id = t.id 
		AND t.email IN (?` + strings.Repeat(`,?`, len(args) - 1) + `)`, args...)		
	
	if err != nil {
		return
	}

	defer rows.Close()

	students := make(map[string]int)
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return
		}
		students[email]++
	}

	for email, count := range students {
		if count == len(teachers) - 1 {
			commonStudents = append(commonStudents, email)
		}
	}

	return commonStudents, nil
}

func (d *DbAdminService) SuspendStudent(student string) error {
	_, err := d.db.Exec("UPDATE students SET suspended = TRUE WHERE email = ?", student)
	return err
}

func (d *DbAdminService) RetrieveForNotifications(teacher string, notification string) ([]string, error) {		
	var teacherId int64;
	err := d.db.QueryRow("SELECT id FROM teachers WHERE email = ?", teacher).Scan(&teacherId)
	if err != nil {
		return nil, err
	}
	
	students := make(map[string]bool)
	
	var taggedStudents[]string;

	pattern := regexp.MustCompile(`@([a-zA-Z0-9+._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+)`)
	matches := pattern.FindAllString(notification, -1)
	for _, match := range matches {
		taggedStudents = append(taggedStudents, match[1:])
	}

	args := make([]interface{}, len(taggedStudents))
	for i, v := range taggedStudents {
		args[i] = v
	}
	
	if len(args) > 0 {
		registeredRows, err := d.db.Query(`SELECT email FROM students 
		WHERE suspended = FALSE AND 
		email IN (?` + strings.Repeat(`,?`, len(args) - 1) + `)`, args...)
		if err != nil {
			return nil, err
		}
		defer registeredRows.Close()

		for registeredRows.Next() {
			var email string
			err := registeredRows.Scan(&email)
			if err != nil {
				return nil, err
			}
			students[email] = true
		}
	}

	taggedRows, err := d.db.Query(`SELECT email FROM students
		WHERE suspended = FALSE AND 
		id IN (SELECT student_id FROM teachers_students WHERE teacher_id = ?)`, teacherId)	
	if err != nil {
		return nil, err
	}
	defer taggedRows.Close()

	for taggedRows.Next() {
		var email string
		err := taggedRows.Scan(&email)
		if err != nil {
			return nil, err
		}
		students[email] = true
	}
	
	recipients := make([]string, 0, len(students))
	for key := range students {
		recipients = append(recipients, key)
	}

	return recipients, err
}

func Init() *sql.DB {
	cfg := mysql.Config{
		User: os.Getenv("DB_USER"),
		Passwd: os.Getenv("DB_PASSWORD"),
		Net: "tcp",
		Addr: os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
	}
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		panic(pingErr)
	}
	return db;
}
