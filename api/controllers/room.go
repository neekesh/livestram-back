package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoomControllers struct {
	// Room repository.RoomRepository
}

func NewRoomControllers() RoomControllers {
	return RoomControllers{}
}

func (rc RoomControllers) GetAllRoom(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "get all teh Room",
		"data": "Rooms",
	})
}

func (rc RoomControllers) PostRoom(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{
		"msg":  "new Room craete",
		"data": "",
	})
}

func (rc RoomControllers) DeleteRoom(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Room deleted",
	})
}

func (rc RoomControllers) UpdateRoom(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "Update Room data",
		"data": "",
	})
}
