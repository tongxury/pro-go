package events

import (
	"github.com/go-mysql-org/go-mysql/canal"
	"store/pkg/sdk/conv"
)

// bigint类型的字段在传输过程中值会变化（不明原因），为避免数据问题，统一转换成string
type Row []string

type CanalEvent struct {
	TableName string
	Action    string
	RowBefore Row
	RowAfter  Row
}

func FromRowsEvent(event *canal.RowsEvent) CanalEvent {

	var rowBefore, rowAfter Row
	if len(event.Rows) == 1 {
		for _, c := range event.Rows[0] {
			rowAfter = append(rowAfter, conv.String(c))
		}
	} else if len(event.Rows) == 2 {
		for _, c := range event.Rows[0] {
			rowBefore = append(rowBefore, conv.String(c))
		}
		for _, c := range event.Rows[1] {
			rowAfter = append(rowAfter, conv.String(c))
		}
	}

	return CanalEvent{
		TableName: event.Table.Name,
		Action:    event.Action,
		RowBefore: rowBefore,
		RowAfter:  rowAfter,
	}
}
