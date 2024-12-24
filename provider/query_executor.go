package provider

import (
	"github.com/traceableai/terraform-provider-traceable/provider/common"
)

func ExecuteQuery(query string, meta interface{}) (string, error) {
	return common.CallExecuteQuery(query, meta)
}
