package websockets

import (
	"strings"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/Jleagle/steam-go/steam"
	"github.com/gamedb/gamedb/pkg/helpers"
	"github.com/gamedb/gamedb/pkg/log"
	"github.com/gamedb/gamedb/pkg/sql"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type WebsocketPage string

const (
	PageAdmin    WebsocketPage = "admin"
	PageApp      WebsocketPage = "app"
	PageBundle   WebsocketPage = "bundle"
	PageBundles  WebsocketPage = "bundles"
	PageChanges  WebsocketPage = "changes"
	PageChat     WebsocketPage = "chat"
	PageGroup    WebsocketPage = "group"
	PageNews     WebsocketPage = "news"
	PagePackage  WebsocketPage = "package"
	PagePackages WebsocketPage = "packages"
	PagePrices   WebsocketPage = "prices"
	PagePlayer   WebsocketPage = "profile"
)

var (
	Pages = map[WebsocketPage]*Page{}
)

func init() {

	pagesSlice := []WebsocketPage{
		PageAdmin,
		PageApp,
		PageBundle,
		PageBundles,
		PageChanges,
		PageChat,
		PageGroup,
		PageNews,
		PagePackage,
		PagePackages,
		PagePrices,
		PagePlayer,
	}
	for _, v := range pagesSlice {
		Pages[v] = &Page{
			name:        v,
			connections: map[uuid.UUID]*websocket.Conn{},
		}
	}
}

func GetPage(page WebsocketPage) (ret *Page) {

	if val, ok := Pages[page]; ok {
		return val
	}

	return ret
}

type Page struct {
	name        WebsocketPage
	connections map[uuid.UUID]*websocket.Conn
	sync.Mutex
}

func (p Page) GetName() WebsocketPage {
	return p.name
}

func (p Page) CountConnections() int {
	return len(p.connections)
}

func (p *Page) AddConnection(conn *websocket.Conn) {

	p.Lock()
	defer p.Unlock()

	p.connections[uuid.NewV4()] = conn
}

func (p *Page) Send(data interface{}) {

	p.Lock()
	defer p.Unlock()

	if p.CountConnections() > 0 {

		var connsToDelete []uuid.UUID

		payload := WebsocketPayload{}
		payload.Page = p.name
		payload.Data = data

		for k, v := range p.connections {

			err := v.WriteJSON(payload)
			if err != nil {

				if strings.Contains(err.Error(), "broken pipe") {

					connsToDelete = append(connsToDelete, k)

				} else {

					log.Err(err)
				}
			}
		}

		for _, v := range connsToDelete {
			delete(p.connections, v)
		}
	}
}

type WebsocketPayload struct {
	Data  interface{}
	Page  WebsocketPage
	Error string
}

// Converts pubsub messages into websockets
func ListenToPubSub() {

	err := helpers.PubSubSubscribe(helpers.PubSubWebsockets, func(m *pubsub.Message) {

		// log.Info("PubSub (" + humanize.Bytes(uint64(len(m.Data))) + "): " + string(m.Data))

		pubSubMsg := PubSubBasePayload{}
		err := helpers.Unmarshal(m.Data, &pubSubMsg)
		if err != nil {
			log.Critical(err)
			return
		}

		for _, page := range pubSubMsg.Pages {

			wsPage := GetPage(page)

			if wsPage.CountConnections() == 0 {
				continue
			}

			switch page {
			case PageApp, PageBundle, PagePackage:

				idPayload := PubSubIDPayload{}

				err = helpers.Unmarshal(m.Data, &idPayload)
				log.Err(err)

				wsPage.Send(idPayload.ID)

			case PagePlayer:

				idPayload := PubSubIDStringPayload{} // int64 too large for js

				err = helpers.Unmarshal(m.Data, &idPayload)
				log.Err(err)

				wsPage.Send(idPayload.ID)

			case PageGroup:

				idPayload := PubSubIDStringPayload{} // ID too large for int64

				err = helpers.Unmarshal(m.Data, &idPayload)
				log.Err(err)

				wsPage.Send(idPayload.ID)

			case PageChanges:

				changePayload := PubSubChangesPayload{}

				err = helpers.Unmarshal(m.Data, &changePayload)
				log.Err(err)

				wsPage.Send(changePayload.Data)

			case PagePackages:

				idPayload := PubSubIDPayload{}

				err = helpers.Unmarshal(m.Data, &idPayload)
				log.Err(err)

				pack, err := sql.GetPackage(idPayload.ID)
				if err != nil {
					log.Err(err)
					continue
				}

				wsPage.Send(pack.OutputForJSON(steam.CountryUS))

			case PageBundles:

				idPayload := PubSubIDPayload{}

				err = helpers.Unmarshal(m.Data, &idPayload)
				log.Err(err)

				bundle, err := sql.GetBundle(idPayload.ID, nil)
				if err != nil {
					log.Err(err)
					continue
				}

				wsPage.Send(bundle.OutputForJSON())

			default:
				log.Err("no handler for page: " + string(page))
			}
		}
	})
	log.Err(err)
}
