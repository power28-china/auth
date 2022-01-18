package domain

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/power28-china/auth/config"
	"github.com/power28-china/auth/database/mongo"
	"github.com/power28-china/auth/utils/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	err      error
	reqJSON  []byte
	response *http.Response
)

// AuthApp represents authentication object for application.
type AuthApp struct {
	AppID           string `json:"appId"`
	CorpAccessToken string `json:"corpAccessToken"`
	CorpID          string `json:"corpId"`
	ExpiresIn       int    `json:"expiresIn"`
}

// An InvalidPtrError describes an invalid argument passed to Unmarshal.
// (The argument to Unmarshal must be a non-nil pointer.)
type InvalidPtrError struct {
	Type reflect.Type
}

// AppAuthResponse represents the response of App authentication from fenxiang open api.
type AppAuthResponse struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	AuthApp
}

func (e *InvalidPtrError) Error() string {
	if e.Type == nil {
		return "json: Unmarshal(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "json: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "json: Unmarshal(nil " + e.Type.String() + ")"
}

// Auth get corperation access token and ID from AppAuthResponse object.
func (authApp *AuthApp) Auth() error {

	request := make(map[string]interface{})

	request["appId"] = config.Config("APP_ID")
	request["appSecret"] = config.Config("APP_SECRET")
	request["permanentCode"] = config.Config("PERMANENT_CODE")

	var appAuthResponse AppAuthResponse

	if err := query("POST", "/cgi/corpAccessToken/get/V2", request, &appAuthResponse); err != nil {
		return err
	}

	if appAuthResponse.ErrorCode != 0 {
		errMessage := fmt.Sprintf("Error Message for App Authentication: %s", appAuthResponse.ErrorMessage)
		logger.Sugar.Errorf(errMessage)
		err = errors.New(errMessage)
		return err
	}

	authApp.AppID = config.Config("APP_ID")
	authApp.CorpAccessToken = appAuthResponse.AuthApp.CorpAccessToken
	authApp.CorpID = appAuthResponse.AuthApp.CorpID
	authApp.ExpiresIn = appAuthResponse.AuthApp.ExpiresIn

	// save result in database.
	filter := bson.D{primitive.E{Key: "appid", Value: authApp.AppID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "corpaccesstoken", Value: authApp.CorpAccessToken},
		primitive.E{Key: "corpid", Value: authApp.CorpID},
		primitive.E{Key: "expiresin", Value: authApp.ExpiresIn}}}}
	mongo.Update(config.Config("AUTH_COLLECTION"), filter, update)

	// Create TTL Index for the collection.
	if _, err := mongo.CreateTTLIndex(config.Config("AUTH_COLLECTION"), int32(1)); err != nil {
		return err
	}

	return nil
}

// query fenxiang API service by given uri,responses and request parameters
func query(method string, uri string, request map[string]interface{}, responseObject interface{}) error {

	apiURL := config.Config("API_HOST") + uri

	transport := &http.Transport{
		// This is the insecure setting, it should be set to false.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}

	switch method {
	case "POST":
		{
			// responseObject must be a `Pointer` and should not be nil
			rv := reflect.ValueOf(responseObject)
			if rv.Kind() != reflect.Ptr || rv.IsNil() {
				return &InvalidPtrError{reflect.TypeOf(responseObject)}
			}

			// convert request object to json format.
			reqJSON, err = json.Marshal(request)
			if err != nil {
				return err
			}

			// logger.Sugar.Debugf("Request JSON:%s \n", string(reqJSON))

			if response, err = client.Post(apiURL, "application/json", bytes.NewBuffer(reqJSON)); err != nil {
				return err
			}

			defer response.Body.Close()

			responseData, err := ioutil.ReadAll(response.Body)

			// logger.Sugar.Debugf("Response JSON:%s \n", string(responseData))

			if err != nil {
				return err
			}

			json.Unmarshal(responseData, &responseObject)
		}
	case "GET":
		{
			if response, err = client.Get(apiURL); err != nil {
				return err
			}

			defer response.Body.Close()
		}
	}

	return nil
}
