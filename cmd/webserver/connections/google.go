package connections

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Jleagle/session-go/session"
	"github.com/gamedb/gamedb/pkg/config"
	"github.com/gamedb/gamedb/pkg/helpers"
	"github.com/gamedb/gamedb/pkg/log"
	"github.com/gamedb/gamedb/pkg/mongo"
	"golang.org/x/oauth2"
	goog "golang.org/x/oauth2/google"
)

type google struct {
}

func (g google) getID(r *http.Request, token *oauth2.Token) interface{} {

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Err(err)
		err = session.SetFlash(r, helpers.SessionBad, "Invalid token")
		log.Err(err)
		return nil
	}
	defer func(response *http.Response) {
		err := response.Body.Close()
		log.Err(err)
	}(response)

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Err(err)
		err = session.SetFlash(r, helpers.SessionBad, "An error occurred (1004)")
		log.Err(err)
		return nil
	}

	userInfo := struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Picture    string `json:"picture"`
		Locale     string `json:"locale"`
	}{}

	err = json.Unmarshal(b, &userInfo)
	if err != nil {
		log.Err(err)
		err = session.SetFlash(r, helpers.SessionBad, "An error occurred (1005)")
		log.Err(err)
		return nil
	}

	return userInfo.ID
}

func (g google) getName() string {
	return "Google"
}

func (g google) getEnum() connectionEnum {
	return ConnectionGoogle
}

func (g google) getConfig() oauth2.Config {
	return oauth2.Config{
		ClientID:     config.Config.GoogleOauthClientID.Get(),
		ClientSecret: config.Config.GoogleOauthClientSecret.Get(),
		Scopes:       []string{"profile"},
		RedirectURL:  config.Config.GameDBDomain.Get() + "/settings/google-callback",
		Endpoint:     goog.Endpoint,
	}
}

func (g google) getEmptyVal() interface{} {
	return ""
}

func (g google) LinkHandler(w http.ResponseWriter, r *http.Request) {
	linkOAuth(w, r, g)
}

func (g google) UnlinkHandler(w http.ResponseWriter, r *http.Request) {
	unlink(w, r, g, mongo.EventUnlinkGoogle)
}

func (g google) CallbackHandler(w http.ResponseWriter, r *http.Request) {

	callbackOAuth(w, r, g, mongo.EventLinkGoogle)

	err := session.Save(w, r)
	log.Err(err)

	http.Redirect(w, r, "/settings", http.StatusFound)
}
