package json

import (
	"encoding/json"
	"github.com/imulab/go-scim/src/core"
	"github.com/imulab/go-scim/src/core/prop"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

func TestSerialize(t *testing.T) {
	s := new(JSONSerializeTestSuite)
	s.resourceBase = "../tests/json_serialize_test_suite"
	suite.Run(t, s)
}

type JSONSerializeTestSuite struct {
	suite.Suite
	resourceBase string
}

func (s *JSONSerializeTestSuite) TestSerialize() {
	tests := []struct{
		name		string
		getResource	func(t *testing.T) *prop.Resource
		getOption	func() *options
		expect		func(t *testing.T, raw []byte, err error)
	}{
		{
			name: 	"default serialize",
			getResource: func(t *testing.T) *prop.Resource {
				_ = s.mustSchema("/user_schema.json")
				resource := prop.NewResourceOf(s.mustResourceType("/user_resource_type.json"), map[string]interface{}{
					"schemas": []interface{}{
						"urn:ietf:params:scim:schemas:core:2.0:User",
					},
					"id": "3cc032f5-2361-417f-9e2f-bc80adddf4a3",
					"meta": map[string]interface{}{
						"resourceType": "User",
						"created": "2019-11-20T13:09:00",
						"lastModified": "2019-11-20T13:09:00",
						"location": "https://identity.imulab.io/Users/3cc032f5-2361-417f-9e2f-bc80adddf4a3",
						"version": "W/\"1\"",
					},
					"userName": "imulab",
					"name": map[string]interface{}{
						"formatted": "Mr. Weinan Qiu",
						"familyName": "Qiu",
						"givenName": "Weinan",
						"honorificPrefix": "Mr.",
					},
					"displayName": "Weinan",
					"profileUrl": "https://identity.imulab.io/profiles/3cc032f5-2361-417f-9e2f-bc80adddf4a3",
					"userType": "Employee",
					"preferredLanguage": "zh_CN",
					"locale": "zh_CN",
					"timezone": "Asia/Shanghai",
					"active": true,
					"emails": []interface{}{
						map[string]interface{}{
							"value": "imulab@foo.com",
							"type": "work",
							"primary": true,
							"display": "imulab@foo.com",
						},
						map[string]interface{}{
							"value": "imulab@bar.com",
							"type": "home",
							"display": "imulab@bar.com",
						},
					},
					"phoneNumbers": []interface{}{
						map[string]interface{}{
							"value": "123-45678",
							"type": "work",
							"primary": true,
							"display": "123-45678",
						},
						map[string]interface{}{
							"value": "123-45679",
							"type": "work",
							"display": "123-45679",
						},
					},
				})
				require.NotNil(t, resource)
				return resource
			},
			getOption: func() *options {
				return Options()
			},
			expect: func(t *testing.T, raw []byte, err error) {
				assert.Nil(t, err)
				expect := `
{
   "schemas":[
      "urn:ietf:params:scim:schemas:core:2.0:User"
   ],
   "id":"3cc032f5-2361-417f-9e2f-bc80adddf4a3",
   "meta":{
      "resourceType":"User",
      "created":"2019-11-20T13:09:00",
      "lastModified":"2019-11-20T13:09:00",
      "location":"https://identity.imulab.io/Users/3cc032f5-2361-417f-9e2f-bc80adddf4a3",
      "version":"W/\"1\""
   },
   "userName":"imulab",
   "name":{
      "formatted":"Mr. Weinan Qiu",
      "familyName":"Qiu",
      "givenName":"Weinan",
      "honorificPrefix":"Mr."
   },
   "displayName":"Weinan",
   "profileUrl":"https://identity.imulab.io/profiles/3cc032f5-2361-417f-9e2f-bc80adddf4a3",
   "userType":"Employee",
   "preferredLanguage":"zh_CN",
   "locale":"zh_CN",
   "timezone":"Asia/Shanghai",
   "active":true,
   "emails":[
      {
         "value":"imulab@foo.com",
         "type":"work",
         "primary":true,
         "display":"imulab@foo.com"
      },
      {
         "value":"imulab@bar.com",
         "type":"home",
         "display":"imulab@bar.com"
      }
   ],
   "phoneNumbers":[
      {
         "value":"123-45678",
         "type":"work",
         "primary":true,
         "display":"123-45678"
      },
      {
         "value":"123-45679",
         "type":"work",
         "display":"123-45679"
      }
   ]
}
`
				assert.JSONEq(t, expect, string(raw))
			},
		},
	}

	for _, test := range tests {
		s.T().Run(test.name, func(t *testing.T) {
			raw, err := Serialize(test.getResource(t), test.getOption())
			test.expect(t, raw, err)
		})
	}
}

func (s *JSONSerializeTestSuite) mustResourceType(filePath string) *core.ResourceType {
	f, err := os.Open(s.resourceBase + filePath)
	s.Require().Nil(err)

	raw, err := ioutil.ReadAll(f)
	s.Require().Nil(err)

	rt := new(core.ResourceType)
	err = json.Unmarshal(raw, rt)
	s.Require().Nil(err)

	return rt
}

func (s *JSONSerializeTestSuite) mustSchema(filePath string) *core.Schema {
	f, err := os.Open(s.resourceBase + filePath)
	s.Require().Nil(err)

	raw, err := ioutil.ReadAll(f)
	s.Require().Nil(err)

	sch := new(core.Schema)
	err = json.Unmarshal(raw, sch)
	s.Require().Nil(err)

	core.SchemaHub.Put(sch)

	return sch
}