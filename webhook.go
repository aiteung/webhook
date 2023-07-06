package webhook

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/aiteung/atmessage/iteung"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/go-playground/webhooks/v6/gitlab"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func Gitlab(c *fiber.Ctx, publickey string, iddebug string, APINotif string) error {
	openingmsg := "*WebHook Gitlab*\n"
	var msg string
	httpRequest := new(http.Request)
	err := fasthttpadaptor.ConvertRequest(c.Context(), httpRequest, true)
	if err != nil {
		log.Println("Error fasthttpadaptor", err)
	}
	hook, _ := gitlab.New(gitlab.Options.Secret(publickey))
	payload, err := hook.Parse(httpRequest, gitlab.PushEvents, gitlab.JobEvents, gitlab.MergeRequestEvents, gitlab.PipelineEvents)
	if err != nil {
		log.Println("Error parse to payload", err)
	}
	switch pyl := payload.(type) {
	case gitlab.PushEventPayload:
		komsg := ""
		for i, komit := range pyl.Commits {
			appd := strconv.Itoa(i+1) + ". " + komit.Message + "_" + komit.Author.Name + "_\n"
			komsg += appd
		}
		msg = pyl.UserName + "\n" + pyl.UserUsername + "\n" + pyl.Repository.Name + "\n" + pyl.Ref + "\n" + pyl.Project.PathWithNamespace + "\n" + komsg
		iteung.PostNotif(strings.ReplaceAll(openingmsg+msg, "eung", "*ng"), iddebug, APINotif)
	case gitlab.JobEventPayload:
		msg = pyl.ProjectName + "\n" + pyl.BuildStatus + "\nStart:" + pyl.BuildStartedAt.String() + "\nFinish" + pyl.BuildFinishedAt.String()
		iteung.PostNotif(strings.ReplaceAll(openingmsg+msg, "eung", "*ng"), iddebug, APINotif)
	case gitlab.MergeRequestEventPayload:
		msg = fmt.Sprintf("%+v", pyl)
		iteung.PostNotif(strings.ReplaceAll(openingmsg+msg, "eung", "*ng"), iddebug, APINotif)
	case gitlab.PipelineEventPayload:
		stmsg := strings.Join(pyl.ObjectAttributes.Stages, " - ")
		msg = pyl.Project.Name + "\n" + pyl.Project.PathWithNamespace + "\nPipeline Status : *" + pyl.ObjectAttributes.Status + "*\nStages : " + stmsg + "\nCreated : " + pyl.ObjectAttributes.CreatedAt.String() + "\nFinished : " + pyl.ObjectAttributes.FinishedAt.GoString()
		iteung.PostNotif(strings.ReplaceAll(openingmsg+msg, "eung", "*ng"), iddebug, APINotif)
	}
	return c.JSON(msg)

}

func Github(c *fiber.Ctx, publickey string, iddebug string, APINotif string) error {
	openingmsg := "*WebHook Github*\n"
	var msg string
	httpRequest := new(http.Request)
	err := fasthttpadaptor.ConvertRequest(c.Context(), httpRequest, true)
	if err != nil {
		log.Println("Error fasthttpadaptor", err)
	}
	hook, _ := github.New(github.Options.Secret(publickey))
	payload, err := hook.Parse(httpRequest, github.PushEvent, github.WorkflowJobEvent, github.WorkflowRunEvent, github.WorkflowDispatchEvent)
	if err != nil {
		log.Println("Error parse to payload", err)
	}
	switch pyl := payload.(type) {
	case github.PushPayload:
		komsg := ""
		for i, komit := range pyl.Commits {
			appd := strconv.Itoa(i+1) + ". " + komit.Message + "\n_" + komit.Author.Name + "_\n"
			komsg += appd
		}
		msg = pyl.Pusher.Name + "\n" + pyl.Sender.Login + "\n" + pyl.Repository.Name + "\n" + pyl.Ref + "\n" + pyl.Repository.URL + "\n" + komsg
		iteung.PostNotif(strings.ReplaceAll(openingmsg+msg, "eung", "*ng"), iddebug, APINotif)
	case github.WorkflowJobPayload:
		msg = pyl.Repository.FullName + "\n" + pyl.Action + "\nStatus:" + pyl.WorkflowJob.Status + "\nCompleted At : " + pyl.WorkflowJob.CompletedAt.String()
		iteung.PostNotif(strings.ReplaceAll(openingmsg+msg, "eung", "*ng"), iddebug, APINotif)
	case github.WorkflowRunPayload:
		msg = fmt.Sprintf("%+v", pyl)
		iteung.PostNotif(strings.ReplaceAll(openingmsg+msg, "eung", "*ng"), iddebug, APINotif)
	case github.WorkflowDispatchPayload:
		msg = pyl.Inputs.Name + "\n" + pyl.Repository.Name + "\nSender : *" + pyl.Sender.Login
		iteung.PostNotif(strings.ReplaceAll(openingmsg+msg, "eung", "*ng"), iddebug, APINotif)
	}
	return c.JSON(msg)

}
