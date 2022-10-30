package plugin

import (
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestTrimHelp(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "trim white space from before and after",
			fn: func(t *testing.T) {
				text := `

help line 1
help line 2

`
				expected := "help line 1" +
					"\nhelp line 2"

				actual := trimHelp(text)

				a := assert.New(t)
				a.EqualValues(expected, actual)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

func TestGetFromFieldData(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "get from data correct types and properties",
			fn: func(t *testing.T) {
				data := &framework.FieldData{
					Raw: map[string]interface{}{
						fldProjectName:     "p1",
						fldProjectRoleName: "r1",
						fldTTL:             "2h",
					},
					Schema: getProjectTokenSchema,
				}
				a := assert.New(t)

				projectName, err := getFromFieldData[string](data, fldProjectName)
				a.EqualValues("p1", projectName)
				require.NoError(t, err)

				projectRoleName, err := getFromFieldData[string](data, fldProjectRoleName)
				a.EqualValues("r1", projectRoleName)
				require.NoError(t, err)

				ttl, err := getFromFieldData[int](data, fldTTL)
				a.EqualValues(7200, ttl)
				require.NoError(t, err)
			},
		},
		{
			name: "get from data incorrect type",
			fn: func(t *testing.T) {
				data := &framework.FieldData{
					Raw: map[string]interface{}{
						fldProjectName:     "p1",
						fldProjectRoleName: "r1",
						fldTTL:             "2h",
					},
					Schema: getProjectTokenSchema,
				}
				a := assert.New(t)

				projectName, err := getFromFieldData[int](data, fldProjectName)
				a.EqualValues(0, projectName)
				require.ErrorContains(t, err, "wrong type")
			},
		},
		{
			name: "get from data incorrect property name",
			fn: func(t *testing.T) {
				data := &framework.FieldData{
					Raw: map[string]interface{}{
						fldProjectName:     "p1",
						fldProjectRoleName: "r1",
						fldTTL:             "2h",
					},
					Schema: getProjectTokenSchema,
				}
				a := assert.New(t)

				projectName, err := getFromFieldData[string](data, "incorrect_property_name")
				a.EqualValues("", projectName)
				require.ErrorContains(t, err, "missing data")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

func TestGetTTLFromFieldData(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "get TTL from data correct type and property",
			fn: func(t *testing.T) {
				data := &framework.FieldData{
					Raw: map[string]interface{}{
						fldTTL: "2h",
					},
					Schema: getProjectTokenSchema,
				}
				a := assert.New(t)

				ttl := getTTLFromFieldData(data, fldTTL, 1*time.Hour, 3*time.Hour)
				a.EqualValues(2*time.Hour, ttl)
			},
		},
		{
			name: "get TTL from data - capped",
			fn: func(t *testing.T) {
				data := &framework.FieldData{
					Raw: map[string]interface{}{
						fldTTL: "5h",
					},
					Schema: getProjectTokenSchema,
				}
				a := assert.New(t)

				ttl := getTTLFromFieldData(data, fldTTL, 1*time.Hour, 3*time.Hour)
				a.EqualValues(3*time.Hour, ttl)
			},
		},
		{
			name: "get TTL from data - default",
			fn: func(t *testing.T) {
				data := &framework.FieldData{
					Raw:    map[string]interface{}{},
					Schema: getProjectTokenSchema,
				}
				a := assert.New(t)

				ttl := getTTLFromFieldData(data, fldTTL, 1*time.Hour, 3*time.Hour)
				a.EqualValues(1*time.Hour, ttl)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

func TestGetFromData(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "get from data correct type - string",
			fn: func(t *testing.T) {
				data := map[string]interface{}{
					"p1": "d1",
				}
				expected := "d1"
				actual, err := getFromData[string](data, "p1")
				a := assert.New(t)
				a.EqualValues(expected, actual)
				require.NoError(t, err)
			},
		},
		{
			name: "get from data incorrect type - string",
			fn: func(t *testing.T) {
				data := map[string]interface{}{
					"p1": 10,
				}
				actual, err := getFromData[string](data, "p1")
				a := assert.New(t)
				a.True(actual == "")
				require.ErrorContains(t, err, "wrong data type")
			},
		},
		{
			name: "get from data incorrect property - string",
			fn: func(t *testing.T) {
				data := map[string]interface{}{
					"p1": 10,
				}
				actual, err := getFromData[string](data, "z1")
				a := assert.New(t)
				a.True(actual == "")
				require.ErrorContains(t, err, "missing data")
			},
		},
		{
			name: "get from data correct type - struct",
			fn: func(t *testing.T) {
				atm := accountTokenMetadata{
					Id:          "i1",
					AccountName: "a1",
					TTL:         1 * time.Hour,
				}
				data := map[string]interface{}{
					"atm": atm,
				}
				expected := atm
				actual, err := getFromData[accountTokenMetadata](data, "atm")
				a := assert.New(t)
				a.EqualValues(expected, actual)
				require.NoError(t, err)
			},
		},
		{
			name: "get from data incorrect type - struct",
			fn: func(t *testing.T) {
				atm := accountTokenMetadata{
					Id:          "i1",
					AccountName: "a1",
					TTL:         1 * time.Hour,
				}
				data := map[string]interface{}{
					"atm": atm,
				}
				actual, err := getFromData[projectTokenMetadata](data, "atm")
				expected := projectTokenMetadata{}
				a := assert.New(t)
				a.EqualValues(expected, actual)
				require.ErrorContains(t, err, "wrong data type")
			},
		},
		{
			name: "get from data incorrect property - struct",
			fn: func(t *testing.T) {
				atm := accountTokenMetadata{
					Id:          "i1",
					AccountName: "a1",
					TTL:         1 * time.Hour,
				}
				data := map[string]interface{}{
					"atm": atm,
				}
				actual, err := getFromData[accountTokenMetadata](data, "ptm")
				expected := accountTokenMetadata{}
				a := assert.New(t)
				a.EqualValues(expected, actual)
				require.ErrorContains(t, err, "missing data")
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

func TestToDurationSeconds(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "from seconds",
			fn: func(t *testing.T) {
				duration := 2 * time.Second
				expected := 2
				actual := toDurationSeconds(duration)
				a := assert.New(t)
				a.EqualValues(expected, actual)
			},
		},
		{
			name: "from hours",
			fn: func(t *testing.T) {
				duration := 2 * time.Hour
				expected := 7200
				actual := toDurationSeconds(duration)
				a := assert.New(t)
				a.EqualValues(expected, actual)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}

func TestTokensUtils(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "project token to lease data",
			fn: func(t *testing.T) {
				token := &projectToken{
					metadata: projectTokenMetadata{
						Id:          "some-id",
						ProjectName: "p1",
						RoleName:    "r1",
						TTL:         1 * time.Hour,
					},
					token: "some-token",
				}
				expected := map[string]interface{}{
					fldID:              "some-id",
					fldProjectName:     "p1",
					fldProjectRoleName: "r1",
				}
				actual := token.toLeaseData()
				a := assert.New(t)
				a.EqualValues(expected, actual)
			},
		},
		{
			name: "project token to response data",
			fn: func(t *testing.T) {
				token := &projectToken{
					metadata: projectTokenMetadata{
						Id:          "some-id",
						ProjectName: "p1",
						RoleName:    "r1",
						TTL:         1 * time.Hour,
					},
					token: "some-token",
				}
				expected := map[string]interface{}{
					fldID:              "some-id",
					fldToken:           "some-token",
					fldProjectName:     "p1",
					fldProjectRoleName: "r1",
				}
				actual := token.toResponseData()
				a := assert.New(t)
				a.EqualValues(expected, actual)
			},
		},
		{
			name: "account token to lease data",
			fn: func(t *testing.T) {
				token := &accountToken{
					metadata: accountTokenMetadata{
						Id:          "some-id",
						AccountName: "a1",
						TTL:         1 * time.Hour,
					},
					token: "some-token",
				}
				expected := map[string]interface{}{
					fldID:          "some-id",
					fldAccountName: "a1",
				}
				actual := token.toLeaseData()
				a := assert.New(t)
				a.EqualValues(expected, actual)
			},
		},
		{
			name: "account token to response data",
			fn: func(t *testing.T) {
				token := &accountToken{
					metadata: accountTokenMetadata{
						Id:          "some-id",
						AccountName: "a1",
						TTL:         1 * time.Hour,
					},
					token: "some-token",
				}
				expected := map[string]interface{}{
					fldID:          "some-id",
					fldToken:       "some-token",
					fldAccountName: "a1",
				}
				actual := token.toResponseData()
				a := assert.New(t)
				a.EqualValues(expected, actual)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, test.fn)
	}
}
