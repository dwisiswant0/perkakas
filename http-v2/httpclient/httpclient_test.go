package httpclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gojektech/heimdall/httpclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"gopkg.in/h2non/gock.v1"
)

type httpClientResponse struct {
	Test string "json:`test`"
}

type testCustomHttp struct {
	client *http.Client
}

func (t *testCustomHttp) Do(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth("some-user", "password")
	return t.client.Do(request)
}

type LogTestSuite struct {
	suite.Suite
	HttpClient   *httpclient.Client
	host         string
	endpoint     string
	url          string
	mockResponse map[string]interface{}
}

func (suite *LogTestSuite) SetupTest() {
	suite.HttpClient = NewHttpClient(nil).Client // Default configuration
	suite.host = "http://some-test-url.com"
	suite.endpoint = "test"
	suite.url = fmt.Sprintf("%s/%s", suite.host, suite.endpoint)
	suite.mockResponse = map[string]interface{}{
		"test": "a1234-abcd",
	}
}

func (suite *LogTestSuite) TestGetRequest() {
	defer gock.Off()

	gock.New(suite.host).
		Get(suite.endpoint).
		Reply(200).
		JSON(suite.mockResponse)

	resp, err := suite.HttpClient.Get(suite.url, nil)
	assert.Nil(suite.T(), err, "Nil expected")

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	result := &httpClientResponse{}
	err = json.Unmarshal(bodyByte, result)
	if err != nil {
		suite.FailNow(err.Error())
	}

	assert.Equal(suite.T(), true, gock.IsDone(), "Must be equal")
	assert.Equal(suite.T(), "a1234-abcd", result.Test, "Result is not same")
}

func (suite *LogTestSuite) TestGetRequestWithCustomClient() {
	defer gock.Off()

	gock.New(suite.host).
		Get(suite.endpoint).
		Reply(200).
		JSON(suite.mockResponse)

	suite.HttpClient = NewHttpWithCustomClient(nil, &testCustomHttp{client: http.DefaultClient}).Client

	header := http.Header{}
	resp, err := suite.HttpClient.Get(suite.url, header)
	assert.Nil(suite.T(), err, "Nil expected")

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		suite.FailNow(err.Error())
	}

	result := &httpClientResponse{}
	err = json.Unmarshal(bodyByte, result)
	if err != nil {
		suite.FailNow(err.Error())
	}

	assert.Equal(suite.T(), true, gock.IsDone(), "Must be equal")
	assert.Equal(suite.T(), "a1234-abcd", result.Test, "Result is not same")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(LogTestSuite))
}
