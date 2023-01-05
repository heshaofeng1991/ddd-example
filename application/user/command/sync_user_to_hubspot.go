package command

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	domainUser "github.com/heshaofeng1991/ddd-johnny/domain/user"
	ent "github.com/heshaofeng1991/entgo/ent/gen"
	"github.com/heshaofeng1991/entgo/ent/viewer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const HubspotContactApi = "https://api.hubapi.com/crm/v3/objects/contacts"

var HubspotContactExisted = "Hubspot contact existed"

type HubspotErrorResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Category string `json:"category"`
}

type HubspotContactCreatedResponse struct {
	Properties struct {
		HsObjectID string `json:"hs_object_id"` // nolint: tagliatelle
		Email      string `json:"email"`
		Firstname  string `json:"firstname"`
		Lastname   string `json:"lastname"`
		Phone      string `json:"phone"`
		Website    string `json:"website"`
	} `json:"properties"`
}

type SyncUserToHubspotHandler struct {
	userRepo domainUser.Repository
}

func NewSyncUserToHubspotHandler(userRepo domainUser.Repository) SyncUserToHubspotHandler {
	if userRepo == nil {
		panic("userRepo cannot be nil")
	}
	return SyncUserToHubspotHandler{userRepo: userRepo}
}

func (h *SyncUserToHubspotHandler) Handle(userID int64) error {
	ctx := viewer.NewContext(context.Background(), viewer.UserViewer{T: &ent.Tenant{ID: -1}})

	user, err := h.userRepo.GetUserInfo(ctx, userID)
	if err != nil {
		return err
	}

	contact, err := domainUser.NewHubspotContact(user)
	if err != nil {
		return err
	}

	if user.GuideInfo().HsObjectID == "" {
		contactID, err := h.createHubspotContact(contact)
		if err != nil {
			if err.Error() == HubspotContactExisted {
				contact.HsObjectID = contactID
				if err := h.updateHubspotContact(contact); err != nil {
					return err
				}
				return nil
			}
			return err
		}
		user.GuideInfo().HsObjectID = contactID
		if _, err := h.userRepo.Save(ctx, user, nil); err != nil {
			return err
		}
	}

	if err := h.updateHubspotContact(contact); err != nil {
		return err
	}

	return nil
}

func (h *SyncUserToHubspotHandler) createHubspotContact(contact *domainUser.HubspotContact) (string, error) {
	header := map[string]string{
		"authorization": "Bearer " + os.Getenv("HUBSPOT_API_KEY"),
		"Content-Type":  "application/json",
	}

	resp, code, err := HTTPReq(HubspotContactApi, http.MethodPost, header, domainUser.HubspotPayload{Properties: contact.Properties})
	if err != nil {
		return "", errors.Wrap(err, "post create contacts fail")
	}

	logrus.Infof("POST /crm/v3/objects/contacts resp: %s %d", string(resp), code)

	if code == http.StatusCreated {
		var createResp HubspotContactCreatedResponse

		err = json.Unmarshal(resp, &createResp)
		if err != nil {
			return "", errors.Wrap(err, "json.Unmarshal error")
		}
		return createResp.Properties.HsObjectID, nil
	}

	if code == http.StatusConflict {
		var createResp HubspotErrorResponse
		err = json.Unmarshal(resp, &createResp)
		if err != nil {
			return "", errors.Wrap(err, "json.Unmarshal error")
		}

		if strings.Contains(createResp.Message, "Existing ID") {
			hsObjectIDArr := strings.Split(createResp.Message, ":")
			if len(hsObjectIDArr) >= 2 {
				hsObjectID := hsObjectIDArr[1]
				return strings.TrimSpace(hsObjectID), errors.New(HubspotContactExisted)
			}
		}
	}

	return "", errors.New("create hubspot contact fail")
}

func (h *SyncUserToHubspotHandler) updateHubspotContact(contact *domainUser.HubspotContact) error {
	header := map[string]string{
		"authorization": "Bearer " + os.Getenv("HUBSPOT_API_KEY"),
		"Content-Type":  "application/json",
	}

	resp, code, err := HTTPReq(HubspotContactApi+"/"+contact.HsObjectID, http.MethodPatch,
		header, domainUser.HubspotPayload{
			Properties: contact.Properties,
		})
	if err != nil {
		return errors.Wrap(err, "post update contacts properties fail")
	}

	logrus.Infof("PATCH /crm/v3/objects/contacts/%s %d resp: %s", contact.HsObjectID, code, string(resp))

	if code != http.StatusOK {
		var createResp HubspotErrorResponse

		err = json.Unmarshal(resp, &createResp)
		if err != nil {
			return errors.Wrap(err, "json.Unmarshal error")
		}

		return errors.Wrap(errors.New(createResp.Message), createResp.Message)
	}
	return nil
}

func HTTPReq(url, method string, header map[string]string, reqData interface{}) ([]byte, int, error) {
	client := &http.Client{}

	var r io.Reader

	value, ok := reqData.(string)
	if ok {
		r = strings.NewReader(value)
	} else {
		data, err := json.Marshal(reqData)
		if err != nil {
			return nil, 0, errors.Wrap(err, "")
		}

		r = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, url, r)
	if err != nil {
		return nil, 0, errors.Wrap(err, "")
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, errors.Wrap(err, "")
	}

	return b, resp.StatusCode, nil
}
