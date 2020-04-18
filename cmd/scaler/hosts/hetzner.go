package hosts

import (
	"context"
	"io/ioutil"
	"net/http"

	"github.com/gamedb/gamedb/pkg/config"
	"github.com/gamedb/gamedb/pkg/helpers"
	"github.com/hetznercloud/hcloud-go/hcloud"
)

var (
	hetznerClient  = hcloud.NewClient(hcloud.WithToken(config.Config.HetznerAPIToken.Get()))
	hetznerContext = context.TODO()
)

type Hetzner struct {
}

func (h Hetzner) ListConsumers() (consumers []Consumer, err error) {

	servers, err := hetznerClient.Server.All(hetznerContext)
	if err != nil {
		return consumers, err
	}

	for _, server := range servers {

		var labels []string
		for k := range server.Labels {
			labels = append(labels, k)
		}

		consumers = append(consumers, Consumer{
			ID:        server.ID,
			Name:      server.Name,
			IP:        server.PublicNet.IPv4.IP.String(),
			Tags:      labels,
			CreatedAt: server.Created.Unix(),
			Locked:    server.Locked,
		})
	}

	return consumers, nil
}

func (h Hetzner) CreateConsumer() (c Consumer, err error) {

	gh, ctx := helpers.GetGithub()
	ghResponse, _, _, err := gh.Repositories.GetContents(ctx, "Jleagle", "infrastructure-gamedb", "scaler/cloud-config.yaml", nil)
	if err != nil {
		return c, err
	}

	fileResponse, err := http.Get(ghResponse.GetDownloadURL())
	if err != nil {
		return c, err
	}
	defer fileResponse.Body.Close()

	b, err := ioutil.ReadAll(fileResponse.Body)
	if err != nil {
		return c, err
	}

	_, _, err = hetznerClient.Server.Create(hetznerContext, hcloud.ServerCreateOpts{
		Name:       "gamedb-consumer-" + helpers.RandString(5, helpers.Letters),
		ServerType: &hcloud.ServerType{Name: "cx11"},
		Image:      &hcloud.Image{Name: "debian-10"},
		SSHKeys:    []*hcloud.SSHKey{{ID: config.Config.HetznerSSHKeyID.GetInt()}},
		Datacenter: &hcloud.Datacenter{Name: "nbg1-dc3"},
		UserData:   string(b),
		Labels:     map[string]string{"consumer": "", ConsumerTag: ""},
		Networks:   []*hcloud.Network{{ID: config.Config.HetznerNetworkID.GetInt()}},
	})

	return c, err
}

func (h Hetzner) DeleteConsumer(id int) (err error) {

	_, err = hetznerClient.Server.Delete(hetznerContext, &hcloud.Server{ID: id})
	return err
}