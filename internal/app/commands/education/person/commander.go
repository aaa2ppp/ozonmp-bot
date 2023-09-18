package person

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ozonmp/omp-bot/internal/app/path"
	"github.com/ozonmp/omp-bot/internal/model/education"
	"log"
)

// TODO: Зачем коммандеру интерфейс Service?
//  Чтобы предоставить командам, т.к. они его методы.
//  Imho, по фэн-шую, каждая команда должна требовать свой интерфейс:
//  Describer, Listener и т.д. Команды должны быть отдельными
//  сущностями?

type Service interface {
	Describe(PersonID uint64) (*education.Person, error)
	List(cursor uint64, limit uint64) ([]education.Person, error)
	Create(education.Person) (uint64, error)
	Update(PersonID uint64, Person education.Person) error
	Remove(PersonID uint64) (bool, error)
}

type commander struct {
	domain    string
	subdomain string
	bot       *tgbotapi.BotAPI
	service   Service
}

func newCommander(domain, subdomain string, bot *tgbotapi.BotAPI, service Service) commander {
	return commander{
		domain:    domain,
		subdomain: subdomain,
		bot:       bot,
		service:   service,
	}
}

func (c commander) checkPath(op string, domain, subdomain string) bool {

	if domain != c.domain || subdomain != c.subdomain {
		log.Printf("%s: unknown path - %s/%s", op, domain, subdomain)
		return false
	}

	return true
}

func (c commander) HandleCallback(callback *tgbotapi.CallbackQuery, callbackPath path.CallbackPath) {
	const op = "commander.HandleCallback"

	if !c.checkPath(op, callbackPath.Domain, callbackPath.Subdomain) {
		return
	}

	switch callbackPath.CallbackName {
	case "list":
		c.ListCallback(callback, callbackPath)
	default:
		log.Printf("%s: unknown callback name: %s", op, callbackPath.CallbackName)
	}
}

func (c commander) HandleCommand(inputMsg *tgbotapi.Message, commandPath path.CommandPath) {
	const op = "commander.HandleCommand"

	if !c.checkPath(op, commandPath.Domain, commandPath.Subdomain) {
		return
	}

	switch commandPath.CommandName {
	case "help":
		c.Help(inputMsg)
	case "list":
		c.List(inputMsg)
	case "get":
		c.Get(inputMsg)
	case "delete":
		c.Delete(inputMsg)
	case "new":
		c.New(inputMsg)
	case "edit":
		c.Edit(inputMsg)
	default:
		c.Default(inputMsg)
	}
}
