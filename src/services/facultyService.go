package services

import (
	"context"
	"errors"
	"time"
	"university-management-api/src/data"

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
