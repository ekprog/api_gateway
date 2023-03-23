package tools

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"io"
	"strings"
)

// EXAMPLE
//k, v := builder.BuildFor("is_active")
//query := fmt.Sprintf("UPDATE services SET %s, updated_at=now() WHERE id=$1", k)
//_, err := r.db.ExecContext(ctx, query, v...)

type UpdateReq struct {
	id     int32
	fields map[string]interface{}
}

func NewUpdateReq(id int32, fields map[string]interface{}) *UpdateReq {
	return &UpdateReq{id: id, fields: fields}
}

func NewUpdateReqBytes(data []byte) (*UpdateReq, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, errors.Wrap(err, "validation error")
	}
	var id int32
	if n, ok := m["id"].(float64); ok {
		if n != float64(int32(n)) {
			return nil, errors.Errorf("id should be int32")
		}
		id = int32(n)
	}
	delete(m, "id")

	return NewUpdateReq(id, m), nil
}

func NewUpdateReqReader(reader io.ReadCloser) (*UpdateReq, error) {
	body, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read all in UpdateReq")
	}
	return NewUpdateReqBytes(body)
}

func (b *UpdateReq) BuildFor(allows ...string) (string, []interface{}) {
	update := lo.PickByKeys(b.fields, allows)

	// KEYS
	keys := lo.Map(lo.Keys(update), func(x string, i int) string {
		return fmt.Sprintf("%s=$%d", x, i+2)
	})
	setQuery := strings.Join(keys, ",")

	//VALUES
	values := lo.Values(update)
	values = append([]interface{}{b.id}, values...)

	return setQuery, values
}
