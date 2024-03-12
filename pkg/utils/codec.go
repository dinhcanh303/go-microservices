package utils

import (
	"encoding/json"

	structpb "github.com/golang/protobuf/ptypes/struct"
	"google.golang.org/protobuf/encoding/protojson"
)

func MapToProtobufStruct(m map[string]interface{}) (*structpb.Struct, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	s := &structpb.Struct{}
	err = protojson.Unmarshal(b, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func ProtobufStructToMap(s *structpb.Struct) (map[string]interface{}, error) {
	b, err := protojson.Marshal(s)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func StructToProtobufStruct(s interface{}) (*structpb.Struct, error) {
	return MapToProtobufStruct(s.(map[string]interface{}))
}

func ProtobufStructToStruct(s *structpb.Struct) (interface{}, error) {
	return ProtobufStructToMap(s)
}
