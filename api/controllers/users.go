package controllers

import (
	"errors"

	dUser "github.com/GotBot-AI/users-api/models/user"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func GetAllUsers(c *fiber.Ctx) error {
	users, err := dUser.GetAll()
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = c.Status(200).JSON(users)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func GetUserByID(c *fiber.Ctx) error {
	userID := c.Params("userID", "")
	if userID == "" {
		return errors.New("not userID provided")
	}

	user, err := dUser.GetByID(userID)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = c.Status(200).JSON(user)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func CreateUser(c *fiber.Ctx) error {
	var user dUser.User

	err := c.BodyParser(&user)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = dUser.Validate(user)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// Uses pointer to get the updated ID and timestamps
	_, err = dUser.Create(&user)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = c.Status(200).JSON(user)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func UpdateUser(c *fiber.Ctx) error {
	var user dUser.User
	err := c.BodyParser(&user)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = dUser.Validate(user)
	if err != nil {
		logrus.Error(err)
		return err
	}

	// Uses pointer to get the updated ID and timestamps
	_, err = dUser.Update(&user)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = c.Status(200).JSON(user)
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func DeleteUser(c *fiber.Ctx) error {
	userID := c.Params("userID", "")
	if userID == "" {
		return errors.New("not userID provided")
	}

	err := dUser.Delete(userID)
	if err != nil {
		logrus.Error(err)
		return err
	}
	c.Status(200).JSON(fiber.Map{"success": true})
	return nil
}
