package plugin

import (
	"fmt"
	"github.com/hashicorp/vault/sdk/framework"
	"reflect"
	"strings"
	"time"
)

func trimHelp(s string) string {
	return strings.Trim(s, "\n\t ")
}

func getFromFieldData[T interface{}](data *framework.FieldData, attr string) (value T, _ error) {
	rawValue, ok := data.GetOk(attr)
	if !ok {
		return value, fmt.Errorf("missing data: %s not present in field data", attr)
	}

	value, ok = rawValue.(T)
	if !ok {
		return value, fmt.Errorf("incorrect data: wrong type %T for %s", reflect.TypeOf(value), attr)
	}

	return value, nil
}

func getTTLFromFieldData(data *framework.FieldData, attr string, defaultTTL time.Duration, maxTTL time.Duration) time.Duration {
	ttlSeconds, err := getFromFieldData[int](data, attr)
	var ttl time.Duration
	if err != nil {
		ttl = defaultTTL
	} else {
		ttl = time.Duration(ttlSeconds) * time.Second
	}

	if ttl > maxTTL {
		ttl = maxTTL
	}

	return ttl
}

func getFromData[T interface{}](data map[string]interface{}, attr string) (value T, _ error) {
	rawValue, ok := data[attr]
	if !ok {
		return value, fmt.Errorf("missing data: %s not present in secret internal data", attr)
	}

	value, ok = rawValue.(T)
	if !ok {
		return value, fmt.Errorf("incorrect data: wrong data type %T for %s", reflect.TypeOf(value), attr)
	}

	return value, nil
}

func toDurationSeconds(d time.Duration) int64 {
	return int64(d / time.Second)
}

func (token *projectToken) toResponseData() map[string]interface{} {
	return map[string]interface{}{
		fldID:              token.metadata.Id,
		fldProjectName:     token.metadata.ProjectName,
		fldProjectRoleName: token.metadata.RoleName,
		fldToken:           token.token,
	}
}

func (token *accountToken) toResponseData() map[string]interface{} {
	return map[string]interface{}{
		fldID:          token.metadata.Id,
		fldAccountName: token.metadata.AccountName,
		fldToken:       token.token,
	}
}

func (token *accountToken) toLeaseData() map[string]interface{} {
	return map[string]interface{}{
		fldID:          token.metadata.Id,
		fldAccountName: token.metadata.AccountName,
	}
}

func (token *projectToken) toLeaseData() map[string]interface{} {
	return map[string]interface{}{
		fldID:              token.metadata.Id,
		fldProjectName:     token.metadata.ProjectName,
		fldProjectRoleName: token.metadata.RoleName,
	}
}

func newTokenSecret(secretType string, duration time.Duration) *framework.Secret {
	return &framework.Secret{
		Type:            secretType,
		DefaultDuration: duration,
	}
}
