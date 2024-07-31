package handlers

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func StreamVideo(c *gin.Context) {
	c.Header("Content-Type", "video/mp4")

	video_path, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not retrieve current directory")
	} else {
		print(video_path)
	}

	video, err := os.Open(video_path + "\\handlers\\person_walking.mp4")
	if err != nil {
		log.Fatal(err)
	}
	defer video.Close()

	http.ServeContent(c.Writer, c.Request, "video.mp4", time.Now(), video)
}
