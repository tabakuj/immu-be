package persistance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	cerror "immudb/internal/errors"
	"immudb/internal/models"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ImmmuDB struct {
	url            string
	apiKey         string
	searchUrl      string
	DefaultHeaders map[string]string
	mx             sync.Mutex
}

func NewImmmuDB(url, apiKey, searchUrl string) *ImmmuDB {
	return &ImmmuDB{
		url:       url,
		apiKey:    apiKey,
		searchUrl: searchUrl,
		DefaultHeaders: map[string]string{
			"accept":       "application/json",
			"Content-Type": "application/json",
			"X-API-Key":    apiKey,
		},
		mx: sync.Mutex{},
	}
}

func (db *ImmmuDB) doCreateHttpCall(ctx context.Context, input interface{}) (*CreateResponse, error) {
	// usually this needs to be taken from the configurations but for this sample application i am leaving it here
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	jsonData, err := json.Marshal(input)
	if err != nil {
		logrus.WithError(err).Error("json marshal failed")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", db.url, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err).Error("http request failed")
		return nil, err
	}
	for key, value := range db.DefaultHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("http request failed")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("http response body read failed")
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var result CreateResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			logrus.WithError(err).Error("json unmarshal failed")
			return nil, err
		}
		return &result, nil
	}
	logrus.WithField("response", string(body)).Error("http response body read failed")
	return nil, fmt.Errorf("http call failed with status code: %d", resp.StatusCode)
}

func (db *ImmmuDB) doGetAllHttpCall(ctx context.Context, input interface{}) (*GetAllResponse, error) {
	// usually this needs to be taken from the configurations but for this sample application i am leaving it here
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	jsonData, err := json.Marshal(input)
	if err != nil {
		logrus.WithError(err).Error("json marshal failed")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", db.searchUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err).Error("http request failed")
		return nil, err
	}
	for key, value := range db.DefaultHeaders {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("http request failed")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).Error("http response body read failed")
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		var result GetAllResponse
		err = json.Unmarshal(body, &result)
		if err != nil {
			logrus.WithError(err).Error("json unmarshal failed")
			return nil, err
		}
		return &result, nil
	}
	logrus.WithField("response", string(body)).Error("http response body read failed")
	return nil, fmt.Errorf("http call failed with status code: %d", resp.StatusCode)
}

func (db *ImmmuDB) CreateAccountInfo(ctx context.Context, data models.AccountInfo) (*models.AccountInfo, error) {
	db.mx.Lock()
	defer db.mx.Unlock()
	data.Id = db.GetId()
	result, err := db.doCreateHttpCall(ctx, data)
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseUint(result.TransactionID, 10, 64)
	if err != nil {
		return nil, err
	}
	if id == 0 {
		return nil, errors.New("invalid transaction-id")
	}
	return &data, nil
}

func (db *ImmmuDB) GetAllAccountInfos(ctx context.Context, pageNr, pageSize int) ([]*models.AccountInfo, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// check if we have any cancellation before continuing
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	req := GetAllSimpleRequest{Page: pageNr, PerPage: pageSize}
	result, err := db.doGetAllHttpCall(ctx, req)
	if err != nil {
		return nil, err
	}

	var output []*models.AccountInfo
	if result != nil && result.Revisions != nil {
		for _, value := range result.Revisions {
			tmp := db.convertModelToAccountInfo(value.Document)
			if tmp != nil {
				output = append(output, tmp)
			}
		}
	}

	return output, nil
}

func (db *ImmmuDB) GetAccountInfoById(ctx context.Context, Id uint) (*models.AccountInfo, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// check if we have any cancellation before continuing
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	requestData := GetAllRequest{
		Query: Query{
			Expressions: []Expression{
				{
					FieldComparisons: []FieldComparison{
						{
							Field:    "id",
							Operator: "EQ",
							Value:    Id,
						},
					},
				},
			},
			Limit: 0,
			OrderBy: []OrderBy{
				{
					Desc:  true,
					Field: "id",
				},
			},
		},
		Page:    1,
		PerPage: 1,
	}
	result, err := db.doGetAllHttpCall(ctx, requestData)
	if err != nil {
		return nil, err
	}

	if result != nil && len(result.Revisions) > 0 {
		return db.convertModelToAccountInfo(result.Revisions[0].Document), nil
	}
	// normally this doesn't need to have acesss to this but i am leaving it for simplity
	return nil, cerror.NewServiceError("account info not found", http.StatusNotFound)
}

// GetId this does not guarantee uniqueness in case of vertical scaling but for purposes of this looks ok
func (db *ImmmuDB) GetId() uint {
	return uint(time.Now().UnixNano())
}

// this is to have optional converting i tried to do this with reflect but i couldn't not get to work as i wanted
func (db *ImmmuDB) convertModelToAccountInfo(sc interface{}) *models.AccountInfo {
	jsonString, err := json.Marshal(sc)
	if err != nil {
		return nil
	}
	var result models.AccountInfo
	err = json.Unmarshal(jsonString, &result)
	if err != nil {
		return nil
	}
	return &result
}
