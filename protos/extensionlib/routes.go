package extensionlib

import (
	"fmt"
	"github.com/NubeIO/rxlib/ginlib"
	"github.com/NubeIO/rxlib/protos/runtimebase/reactive"
	"github.com/NubeIO/rxlib/protos/runtimebase/runtime"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
	"net/http"
)

func (inst *Extensions) bootServer(opts *ginlib.Opts) {
	inst.server = ginlib.NewServer(opts)
}

func (inst *Extensions) infoRoute() {
	inst.server.AddGetRoute("/api/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, extensionInfo{
			Name:     inst.name,
			Version:  "1.0.0",
			BootTime: inst.bootTime,
		})
	})
}

func (inst *Extensions) pingRoute() {
	inst.server.AddGetRoute("/api/ping", func(c *gin.Context) {
		err := inst.sendMessageToServer("hello", "info")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to ping server"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("port: %d", inst.server.GetPort())})
	})
}

func (inst *Extensions) runtimeRoute() {
	inst.server.AddGetRoute("/api/runtime", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"runtime": reactive.ConvertObjects(inst.runtime)})
	})
}

func (inst *Extensions) messagesRoute() {
	inst.server.AddGetRoute("/api/messages", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"messages": messages})
	})
}

func (inst *Extensions) palletRoute() {
	inst.server.AddGetRoute("/api/pallet", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"pallet": reactive.ConvertObjects(maps.Values(inst.pallet))})
	})
}

// streamRoute Send a message to the server
func (inst *Extensions) streamRoute() {
	inst.server.AddGetRoute("/api/stream", func(c *gin.Context) {
		if err := inst.stream.Send(&runtime.MessageRequest{ExtensionUUID: inst.name}); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
}
