package services

import (
	"context"
	"errors"
	"time"
	"university-management-api/src/data"
	"university-management-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
)

func CheckFaculty(faculty_name string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = nil

	count, _ := data.FacultyCollection.CountDocuments(ctx, bson.M{"name": faculty_name})

	if count > 0 {
		err = errors.New("faculty already exist")
		return err
	}
	return err
}

func CheckFacultyById(faculty_id string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = nil

	count, _ := data.FacultyCollection.CountDocuments(ctx, bson.M{"facultyid": faculty_id})

	if count == 0 {
		err = errors.New("faculty does not exist")
		return err
	}
	return err
}

func CheckFacultyDean(dean_id string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = nil

	count, _ := data.FacultyCollection.CountDocuments(ctx, bson.M{"deanid": dean_id})

	if count > 0 {
		err = errors.New("faculty already have this dean")
		return err
	}
	return err
}

func GetFacultyById(faculty_id string) (faculty models.Faculty, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = nil

	foundErr := data.FacultyCollection.FindOne(ctx, bson.M{"facultyid": faculty_id}).Decode(&faculty)

	if foundErr != nil {
		err = foundErr
		return faculty, err
	}
	return faculty, err
}
