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
	"time"

	"github.com/power28-china/auth/config"
	"github.com/power28-china/auth/database/mongo"
	"github.com/power28-china/auth/utils/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
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

// Authentication represents the methods for the fenxiang OpenAPI authentication.
type Authentication interface {
	Auth() error
	GetAuth() error
}

// GetAuth get corperation access token and ID from database.
func (auth *AuthApp) GetAuth() error {
	if err := mongo.Find(config.Config("AUTH_COLLECTION"), "appid", config.Config("APP_ID")).Decode(auth); err != nil {
		if err == mongodb.ErrNoDocuments {
			// if authentication information for app is not available, call `AppAuth` method to create one.
			logger.Sugar.Infof("No APP authentication founded from database, will recreate it after 2 seconds.\n")
			time.Sleep(2 * time.Second)
			if err = auth.Auth(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	resp := map[string]interface{}{}
	// make sure the app authentication is valid.
	for {
		if err := auth.tryToQueryAPI(resp); err != nil {
			return err
		}

		// spew.Dump(resp)

		// if App authentication is invalid, delete app authentication in database and  retry to get new authentication.
		if resp["errorCode"].(float64) == 20016 {
			logger.Sugar.Infof("APP authentication is invalid, recreate it now.")
			mongo.Delete(config.Config("AUTH_COLLECTION"), "appid", config.Config("APP_ID"))
			// logger.Sugar.Debugf("old authentication: %v", appAuth.CorpAccessToken)
			if err = auth.Auth(); err != nil {
				return err
			}
			// logger.Sugar.Debugf("new authentication: %v", appAuth.CorpAccessToken)
			continue
		}

		if resp["errorCode"].(float64) == 20003 {
			logger.Sugar.Infof("Param illegal exception, try again after 2 seconds.")
			time.Sleep(2 * time.Second)
			continue
		}

		if resp["errorCode"].(float64) == 504 {
			logger.Sugar.Infof("Gateway Timeout, try again after 2 seconds.")
			time.Sleep(2 * time.Second)
			continue
		}

		if resp["errorCode"].(float64) != 0 {
			errMessage := fmt.Sprintf("Error Message for getAppAuthentication: %s(%f)", resp["errorMessage"], resp["errorCode"])
			logger.Sugar.Infof(errMessage)
			err := errors.New(errMessage)
			return err
		}

		if resp["errorCode"].(float64) == 0 {
			break
		}
	}

	logger.Sugar.Infof("Get App Authentication from database. token:%s", auth.CorpAccessToken)
	return nil
}

// Auth get corperation access token and ID from AppAuthResponse object.
func (auth *AuthApp) Auth() error {

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

	auth.AppID = config.Config("APP_ID")
	auth.CorpAccessToken = appAuthResponse.AuthApp.CorpAccessToken
	auth.CorpID = appAuthResponse.AuthApp.CorpID
	auth.ExpiresIn = appAuthResponse.AuthApp.ExpiresIn

	// save result in database.
	filter := bson.D{primitive.E{Key: "appid", Value: auth.AppID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "corpaccesstoken", Value: auth.CorpAccessToken},
		primitive.E{Key: "corpid", Value: auth.CorpID},
		primitive.E{Key: "expiresin", Value: auth.ExpiresIn}}}}
	mongo.Update(config.Config("AUTH_COLLECTION"), filter, update)

	// Create TTL Index for the collection.
	if _, err := mongo.CreateTTLIndex(config.Config("AUTH_COLLECTION"), int32(1)); err != nil {
		return err
	}

	return nil
}

func (auth *AuthApp) tryToQueryAPI(response map[string]interface{}) error {
	queryInfo := make(map[string]interface{})
	queryInfo["offset"] = 0
	queryInfo["limit"] = 1
	queryInfo["filters"] = []interface{}{}
	queryInfo["orders"] = []interface{}{}
	queryInfo["fieldProjection"] = []string{"name"}

	data := make(map[string]interface{})
	data["search_query_info"] = queryInfo
	data["dataObjectApiName"] = "AccountObj"

	request := make(map[string]interface{})
	request["corpAccessToken"] = auth.CorpAccessToken
	request["corpId"] = auth.CorpID
	request["currentOpenUserId"] = config.Config("CURRENT_OPENUSER_ID")
	request["data"] = data

	if err := query("POST", config.Config("API_QUERY_URL"), request, &response); err != nil {
		return err
	}

	logger.Sugar.Debugf("query API successfully with token: %v", auth.CorpAccessToken)
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
