package main

import (
	"embed"
	_ "embed"
	"github.com/gin-gonic/gin"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	//go:embed static
	staticTemplateFS  embed.FS
	minDelayMs        = 200
	maxDelayMs        = 1000
	handlerIterations = atomic.Int64{}

	delayRand *rand.Rand
	probRand  *rand.Rand

	returnCodes = []ReturnCodeProbability{
		{
			Probability: 5,
			Code:        500,
			Message:     "You found one of the internal server errors, good job :)",
		},
		{
			Probability: 10,
			Code:        401,
			Message:     "This will happen when users input wrong credentials",
		},
	}
	defaultReturnCode = ReturnCodeProbability{
		Probability: 0,
		Code:        200,
	}
)

type ReturnCodeProbability struct {
	Probability int
	Code        int
	Message     string
}

func init() {
	delayRand = rand.New(rand.NewSource(1337))
	probRand = rand.New(rand.NewSource(1337))
}

func waitDelay() {
	waitDelayScaled(1)
}

func waitDelayScaled(scale int) {
	delay := minDelayMs + delayRand.Intn(maxDelayMs-minDelayMs)*scale
	//fmt.Printf("waiting forced delay %d ms", delay)
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func selectReturnCode() ReturnCodeProbability {
	for _, prob := range returnCodes {
		if probRand.Intn(100) <= prob.Probability {
			return prob
		}
	}
	return defaultReturnCode
}

func handleProbabilityResponse(context *gin.Context, okMessage string) {
	rc := selectReturnCode()
	if rc.Code != http.StatusOK {
		context.JSON(rc.Code, rc.Message)
	} else {
		context.JSON(http.StatusOK, okMessage)
	}
}

func addEndpointHandlers(g *gin.Engine) error {
	if t, err := template.ParseFS(staticTemplateFS, "static/*"); err != nil {
		return err
	} else {
		g.SetHTMLTemplate(t)
	}

	g.Handle("GET", "/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})

	g.Handle("GET", "/admin", func(context *gin.Context) {
		context.HTML(http.StatusOK, "admin.html", nil)
	})

	g.Handle("POST", "/admin/setdelay", func(context *gin.Context) {
		minDelay, _ := context.GetPostForm("mindelay")
		maxDelay, _ := context.GetPostForm("maxdelay")

		if delay, err := strconv.Atoi(minDelay); err == nil {
			minDelayMs = delay
		} else {
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}
		if delay, err := strconv.Atoi(maxDelay); err == nil {
			maxDelayMs = delay
		} else {
			context.AbortWithError(http.StatusBadRequest, err)
			return
		}
		context.HTML(http.StatusOK, "index.html", nil)
	})

	g.Handle("GET", "/health", func(context *gin.Context) {
		context.JSON(http.StatusOK, "ok")
	})

	g.Handle("GET", "/api/logout", func(context *gin.Context) {
		waitDelay()
		handleProbabilityResponse(context, "logout successful")
	})
	g.Handle("GET", "/api/login", func(context *gin.Context) {
		waitDelayScaled(2)
		handleProbabilityResponse(context, "login successful")
	})
	g.Handle("GET", "/api/project", func(context *gin.Context) {
		waitDelayScaled(3)
		handleProbabilityResponse(context, "get project")
	})
	g.Handle("POST", "/api/project", func(context *gin.Context) {
		waitDelayScaled(10)
		handleProbabilityResponse(context, "post project")
	})

	return nil
}
