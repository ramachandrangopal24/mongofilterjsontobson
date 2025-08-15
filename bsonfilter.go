package filter

import (
	"encoding/json"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Request is an optional wrapper if you want sort/limit/projection.
// If you only want raw filter, you can just pass the body to ParseFilter.
type Request struct {
	Filter     json.RawMessage `json:"filter"`     // Extended JSON or plain JSON
	Sort       map[string]int  `json:"sort"`       // e.g. {"createdAt": -1}
	Projection map[string]int  `json:"projection"` // e.g. {"secret": 0}
	Limit      *int64          `json:"limit"`      // e.g. 50
	Skip       *int64          `json:"skip"`       // e.g. 100
}

// ParseFilter takes JSON (prefer Mongo Extended JSON) and returns a BSON document.
// If extJSON parsing fails, it falls back to plain JSON with light coercions (e.g., _id hex → ObjectID).
func ParseFilter(body []byte) (bson.D, error) {
	if len(body) == 0 {
		return bson.D{}, nil
	}
	var d bson.D
	if err := bson.UnmarshalExtJSON(body, true, &d); err == nil {
		return d, nil
	}

	// Fallback to plain JSON -> bson.M with _id coercion
	m := map[string]any{}
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, errors.Join(errors.New("extended json parse failed"), err)
	}
	coerceSpecials(m)
	// keep order-stable by converting to bson.D
	return toDoc(m), nil
}

// ParseRequest parses a wrapper {filter, sort, limit, projection, skip}.
func ParseRequest(body []byte) (f bson.D, sortDoc bson.D, projDoc bson.D, limit, skip *int64, err error) {
	var req Request
	if err = json.Unmarshal(body, &req); err != nil {
		return
	}
	f, err = ParseFilter(req.Filter)
	if err != nil {
		return
	}
	sortDoc = toDocInt(req.Sort)
	projDoc = toDocInt(req.Projection)
	limit, skip = req.Limit, req.Skip
	return
}

func coerceSpecials(v any) {
	switch t := v.(type) {
	case map[string]any:
		for k, vv := range t {
			// Common: _id hex -> ObjectID
			if k == "_id" {
				if s, ok := vv.(string); ok && len(s) == 24 {
					if oid, err := primitive.ObjectIDFromHex(s); err == nil {
						t[k] = oid
						continue
					}
				}
			}
			coerceSpecials(vv)
		}
	case []any:
		for i := range t {
			coerceSpecials(t[i])
		}
	}
}

func toDoc(m map[string]any) bson.D {
	d := bson.D{}
	for k, v := range m {
		d = append(d, bson.E{Key: k, Value: v})
	}
	return d
}

func toDocInt(m map[string]int) bson.D {
	if len(m) == 0 {
		return nil
	}
	d := bson.D{}
	for k, v := range m {
		d = append(d, bson.E{Key: k, Value: v})
	}
	return d
}
