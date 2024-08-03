package persistance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"immudb/internal/models"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type ImmmuDB struct {
	url     string
	apiKey  string
	search  string
	headers map[string]string
	mx      sync.Mutex
}

func NewImmmuDB(url, apiKey, search string) *ImmmuDB {
	return &ImmmuDB{
		url:    url,
		apiKey: apiKey,
		search: search,
		headers: map[string]string{
			"accept":       "application/json",
			"X-API-Key":    apiKey,
			"Content-Type": "application/json",
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
	for key, value := range db.headers {
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

func (db *ImmmuDB) doGetAllHttpCall(ctx context.Context, input GetAllRequest) (*GetAllResponse, error) {
	// usually this needs to be taken from the configurations but for this sample application i am leaving it here
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	jsonData, err := json.Marshal(input)
	if err != nil {
		logrus.WithError(err).Error("json marshal failed")
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", db.url, bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err).Error("http request failed")
		return nil, err
	}
	for key, value := range db.headers {
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
	req := GetAllRequest{Page: pageNr, PerPage: pageSize}
	result, err := db.doGetAllHttpCall(ctx, req)
	if err != nil {
		return nil, err
	}

	var output []*models.AccountInfo
	if result != nil && result.Revisions != nil {
		for _, value := range result.Revisions {
			output = append(output, &value.Document)
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
		Query: &Query{
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
		return &result.Revisions[0].Document, nil
	}
	return nil, fmt.Errorf("account info not found")
}

// GetId this does not guarantee uniqueness in case of vertical scaling but for purposes of this looks ok
func (db *ImmmuDB) GetId() uint {
	db.mx.Lock()
	defer db.mx.Unlock()
	return uint(time.Now().UnixNano())
}
