package utils

import (
	"testing"

	"github.com/power28-china/auth/utils/logger"
	"github.com/stretchr/testify/require"
)

func TestConvertToTableName(t *testing.T) {
	name := "test"
	logger.Sugar.Debugf("table name: %s", ConvertToTableName(name))
}

func TestSnakeCase(t *testing.T) {
	for _, test := range []struct {
		in   string
		want string
	}{
		{
			in:   "",
			want: "",
		},
		{
			in:   "IsDigit",
			want: "is_digit",
		},
		{
			in:   "Is",
			want: "is",
		},
		{
			in:   "IsID",
			want: "is_id",
		},
		{
			in:   "IsSQL",
			want: "is_sql",
		},
		{
			in:   "LongSQL",
			want: "long_sql",
		},
		{
			in:   "Float64Val",
			want: "float64_val",
		},
		{
			in:   "XMLName",
			want: "xml_name",
		},
	} {
		require.Equal(t, test.want, CamelCaseToSnakeCase(test.in))
	}
}

func BenchmarkCamelCaseToSnakeCase(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CamelCaseToSnakeCase("getHTTPResponseCode")
	}
}
