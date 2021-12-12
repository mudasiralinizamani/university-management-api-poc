package services

import (
	"context"
	"errors"
	"time"
	"university-management-api/src/data"

	"go.mongodb.org/mongo-driver/bson"
)

func DoesDepartmentExist(departmentName string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, _ := data.DepartmentCollection.CountDocuments(ctx, bson.M{"name": departmentName})

	if count > 0 {
		err := errors.New("department already exists")
		return err
	}

	return nil
}

func CheckDepartmentById(department_id string) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	count, _ := data.DepartmentCollection.CountDocuments(ctx, bson.M{"departmentid": department_id})

	if count == 0 {
		return errors.New("Department not found")
	}
	return nil
}
