// File containing all the handler functions and their logics separated from other code

package controllers

import (
	"context"
	"my-rest-api/configs"
	"my-rest-api/models"
	"my-rest-api/responses"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// variable to the collection
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

// special validator variable
var validate = validator.New()

// We are going to validate the request body and check whether the fields/attributes are properly set are not to avoid inconsistency
// We are going test for this in CreateUser and EditUser handlers where we receive json in request body

// function responsible for creating a new user in the database
func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var user models.User
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	// filling details in the user model
	// the createdAt attribute is set only at the time of user creation
	// it is not manipulated elsewhere (Updating)
	newUser := models.User{
		Name:        user.Name,
		DOB:         user.DOB,
		Address:     user.Address,
		Description: user.Description,
		CreatedAt:   time.Now().String(),
	}

	// query to insert a user
	result, err := userCollection.InsertOne(ctx, newUser)

	// checking whether an error occured while updating
	// sending an error response to the user if error exists
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// sending correct response upon success
	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

// function responsible for retrieving a user from the database based on UserID
func GetAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// extracting userId from params
	userId := c.Params("userId")

	// user model to store fetched data
	var user models.User

	defer cancel()

	// converting userId from string to ObjectID
	objId, _ := primitive.ObjectIDFromHex(userId)

	// query to fetch an existing users from collection
	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)

	// checking whether an error occured while fetching
	// sending an error response to the user if error exists
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// sending correct response upon success
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": user}})
}

// function responsible for editing a user from the database based on UserID
func EditAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// extracting userId from params
	userId := c.Params("userId")

	// user model to store fetched data
	var user models.User
	defer cancel()

	// converting userId from string to ObjectID
	objId, _ := primitive.ObjectIDFromHex(userId)

	//validate the request body
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&user); validationErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}

	// variable which stores the new user attributes after fetching to be updated
	update := bson.M{"name": user.Name, "dob": user.DOB, "address": user.Address, "description": user.Description}

	// query to update a user based on the "_id" value passed
	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": objId}, bson.M{"$set": update})

	// checking whether an error occured while updating
	// sending an error response to the user if error exists
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	//get updated user details
	var updatedUser models.User

	// After updating the user, fetching back the same user and returning it to the user as a response
	// this code is similar to the fetching a single user code
	if result.MatchedCount == 1 {
		err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&updatedUser)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}
	}

	// sending correct response upon success
	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": updatedUser}})
}

// function responsible for deleting a user from the database based on UserID
func DeleteAUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// extracting userId from params
	userId := c.Params("userId")
	defer cancel()

	// converting userId from string to ObjectID
	objId, _ := primitive.ObjectIDFromHex(userId)

	// query to delete o user based on the "_id" value passed
	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})

	// checking whether an error occured while deleting
	// sending an error response to the user if error exists
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// if deleted users are less than 1 -> No user deleted -> Invalid userId
	// sending error response to the user
	if result.DeletedCount < 1 {
		return c.Status(http.StatusNotFound).JSON(
			responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: &fiber.Map{"data": "User with specified ID not found!"}},
		)
	}

	// sending correct response upon success
	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": "User successfully deleted!"}},
	)
}

// function responsible for retrieving all the user from the database
func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// slice to store all the retrieved users
	var users []bson.M
	defer cancel()

	// query to fetch all existing users from collection
	results, err := userCollection.Find(ctx, bson.M{})

	// checking whether an error occured while fetching
	// sending an error response to the user if error exists
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	defer results.Close(ctx)

	// reading from the db in an optimal way
	// fetching an individual user using a curson and appending it to the users slice
	for results.Next(ctx) {
		var singleUser bson.M

		// sending back error response if error exists
		if err = results.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
		}

		users = append(users, singleUser)
	}

	// sending correct response upon success
	return c.Status(http.StatusOK).JSON(
		responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": users}},
	)
}
