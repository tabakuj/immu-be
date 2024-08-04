package persistance

import (
	"context"
	"github.com/stretchr/testify/assert"
	"immudb/internal/models"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestImmmuDB_GetId(t *testing.T) {
	tests := []struct {
		name         string
		numberOfJobs int
	}{
		{
			name:         "Test_Validity",
			numberOfJobs: 1,
		},
		{
			name:         "Test_5_Parallel_Jobs",
			numberOfJobs: 5,
		},
		{
			name:         "Test_50_Parallel_Jobs",
			numberOfJobs: 50,
		},
		{
			name:         "Test_500_Parallel_Jobs",
			numberOfJobs: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			wg := sync.WaitGroup{}
			mx := sync.Mutex{}
			db := NewImmmuDB("", "", "")
			result := make([]uint, tt.numberOfJobs)

			// action
			for i := 0; i < tt.numberOfJobs; i++ {
				wg.Add(1)
				go func(counter int) {
					defer wg.Done()
					tmpID := db.GetId()
					mx.Lock()
					result[counter] = tmpID
					mx.Unlock()
				}(i)
			}

			wg.Wait()

			// assertions
			t.Logf("result: %v", result)
			tmp := make(map[uint]bool)
			for i := 0; i < tt.numberOfJobs; i++ {
				if _, ok := tmp[result[i]]; ok {
					t.Errorf("GetId() returned the same number more than once")
				}
				tmp[result[i]] = false
			}
		})
	}
}

func TestImmmuDB_CreateAccountInfo(t *testing.T) {
	tests := []struct {
		name               string
		ImmuReturnResponse string
		ImmuReturnStatus   int
		request            models.AccountInfo
		expectedResult     models.AccountInfo
	}{
		{
			name:             "Test_Validity_Sending_Type",
			ImmuReturnStatus: http.StatusOK,
			ImmuReturnResponse: `{
									"documentId": "66af5b390000000000000007f0afc792",
									"transactionId": "8"
			}`,
			request: models.AccountInfo{
				Name:    "test",
				Iban:    "testIban",
				Address: AddPointer("test Address"),
				Amount:  300,
				Type:    AddPointer(models.Sending),
			},
			expectedResult: models.AccountInfo{
				Name:    "test",
				Iban:    "testIban",
				Address: AddPointer("test Address"),
				Amount:  300,
				Type:    AddPointer(models.Sending),
			},
		},
		{
			name:             "Test_Validity_Receiving_Type",
			ImmuReturnStatus: http.StatusOK,
			ImmuReturnResponse: `{
									"documentId": "66af5b390000000000000007f0afc792",
									"transactionId": "8"
			}`,
			request: models.AccountInfo{
				Name:    "test",
				Iban:    "testIban",
				Address: AddPointer("test Address"),
				Type:    AddPointer(models.Receiving),
			},
			expectedResult: models.AccountInfo{
				Name:    "test",
				Iban:    "testIban",
				Address: AddPointer("test Address"),
				Type:    AddPointer(models.Receiving),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.ImmuReturnStatus)
				_, err := w.Write([]byte(tt.ImmuReturnResponse))
				if err != nil {
					t.Error(err, "failed to write response")
				}
			}))
			defer ts.Close()
			db := NewImmmuDB(ts.URL, "", "")

			// action
			result, err := db.CreateAccountInfo(context.Background(), tt.request)
			if err != nil {
				t.Error(err)
			}
			// bc this is set on db level
			tt.expectedResult.Id = result.Id

			// assertions
			assert.Equal(t, &tt.expectedResult, result, "invalid result returned from CreateAccountInfo")
		})
	}
}

func TestImmmuDB_GetById(t *testing.T) {
	tests := []struct {
		name               string
		ImmuReturnResponse string
		ImmuReturnStatus   int
		requestId          uint
		expectedResult     models.AccountInfo
	}{
		{
			name:             "Test_Validity_GetById",
			ImmuReturnStatus: http.StatusOK,
			ImmuReturnResponse: `{
									"page": 1,
									"perPage": 100,
									"revisions": [
										{
											"document": {
												"Address": "1234 Elm Street, Springfield, USA",
												"Amount": 1500.75,
												"Iban": "GB82WEST12345698765432",
												"Type": 1,
												"_id": "66ae4d420000000000000004f0afc78f",
												"_vault_md": {
													"creator": "a:14ffea9d-4313-465e-9f61-ed4f4097ca87",
													"ts": 1722699074
												},
												"id": 1722699072077830000,
												"name": "John Doe 22"
											},
											"revision": "",
											"transactionId": ""
										}
									],
									"searchId": ""
			}`,
			requestId: 1722699072077830000,
			expectedResult: models.AccountInfo{
				Id:      1722699072077830000,
				Name:    "John Doe 22",
				Iban:    "GB82WEST12345698765432",
				Address: AddPointer("1234 Elm Street, Springfield, USA"),
				Amount:  1500.75,
				Type:    AddPointer(models.Sending),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// here basically we can also check the request sent to the server but i want to keep this a little simple
				w.WriteHeader(tt.ImmuReturnStatus)
				_, err := w.Write([]byte(tt.ImmuReturnResponse))
				if err != nil {
					t.Error(err, "failed to write response")
				}
			}))
			defer ts.Close()
			db := NewImmmuDB("", "", ts.URL)

			// action
			result, err := db.GetAccountInfoById(context.Background(), tt.requestId)
			if err != nil {
				t.Error(err)
			}

			// assertions
			assert.Equal(t, &tt.expectedResult, result, "invalid result returned from getById")
		})
	}
}

func TestImmmuDB_GetAll(t *testing.T) {
	tests := []struct {
		name               string
		ImmuReturnResponse string
		ImmuReturnStatus   int
		expectedResult     []*models.AccountInfo
	}{
		{
			name:             "Test_Validity_GetAll",
			ImmuReturnStatus: http.StatusOK,
			ImmuReturnResponse: `{
							"page": 1,
							"perPage": 10,
							"revisions": [
								{
									"document": {
										"Address": "1234 Elm Street, Springfield, USA",
										"Amount": 1500.75,
										"Iban": "GB82WEST12345698765432",
										"Type": 1,
										"_id": "66ae5d610000000000000005f0afc790",
										"_vault_md": {
											"creator": "a:14ffea9d-4313-465e-9f61-ed4f4097ca87",
											"ts": 1722703201
										},
										"id": 1722703201251299000,
										"name": "John Doe 22"
									},
									"revision": "",
									"transactionId": ""
								},
								{
									"document": {
										"Address": null,
										"Amount": 1,
										"Iban": "aa",
										"Type": 1,
										"_id": "66ae9a240000000000000006f0afc791",
										"_vault_md": {
											"creator": "a:14ffea9d-4313-465e-9f61-ed4f4097ca87",
											"ts": 1722718756
										},
										"id": 1722718756365078000,
										"name": "test"
									},
									"revision": "",
									"transactionId": ""
								}
							],
							"searchId": ""
			}`,
			expectedResult: []*models.AccountInfo{
				&models.AccountInfo{
					Id:      1722703201251299000,
					Name:    "John Doe 22",
					Iban:    "GB82WEST12345698765432",
					Address: AddPointer("1234 Elm Street, Springfield, USA"),
					Amount:  1500.75,
					Type:    AddPointer(models.Sending),
				},
				&models.AccountInfo{
					Id:     1722718756365078000,
					Name:   "test",
					Iban:   "aa",
					Amount: 1,
					Type:   AddPointer(models.Sending),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// here basically we can also check the request sent to the server but i want to keep this a little simple
				w.WriteHeader(tt.ImmuReturnStatus)
				_, err := w.Write([]byte(tt.ImmuReturnResponse))
				if err != nil {
					t.Error(err, "failed to write response")
				}
			}))
			defer ts.Close()
			db := NewImmmuDB("", "", ts.URL)

			// action
			result, err := db.GetAllAccountInfos(context.Background(), 1, 10)
			if err != nil {
				t.Error(err)
			}

			// assertions
			assert.Equal(t, len(tt.expectedResult), len(result), "invalid length returned by getAllAccountInfos")
			assert.Equal(t, tt.expectedResult, result, "invalid result returned from getAll")
		})
	}
}
