package protobq

import (
	"cloud.google.com/go/bigquery"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// InferSchema infers a BigQuery schema for the given proto.Message using default options.
func InferSchema(msg proto.Message) bigquery.Schema {
	return SchemaOptions{}.InferSchema(msg)
}

// SchemaOptions contains configuration options for BigQuery schema inference.
type SchemaOptions struct{}

// InferSchema infers a BigQuery schema for the given proto.Message using options in
// MarshalOptions.
func (o SchemaOptions) InferSchema(msg proto.Message) bigquery.Schema {
	return o.inferMessageSchema(msg.ProtoReflect().Descriptor())
}

// inferMessageSchema infers the BigQuery schema for the given protoreflect.MessageDescriptor.
func (o SchemaOptions) inferMessageSchema(msg protoreflect.MessageDescriptor) bigquery.Schema {
	schema := make(bigquery.Schema, 0, msg.Fields().Len())
	for i := 0; i < msg.Fields().Len(); i++ {
		field := msg.Fields().Get(i)
		if field.IsMap() {
			continue // TODO: support maps
		}
		fieldSchema := &bigquery.FieldSchema{
			Name:     string(field.Name()),
			Repeated: field.Cardinality() == protoreflect.Repeated,
		}
		switch field.Kind() {
		case protoreflect.DoubleKind,
			protoreflect.FloatKind:
			fieldSchema.Type = bigquery.FloatFieldType
		case protoreflect.Int64Kind,
			protoreflect.Uint64Kind,
			protoreflect.Int32Kind,
			protoreflect.Fixed64Kind,
			protoreflect.Fixed32Kind,
			protoreflect.Uint32Kind,
			protoreflect.Sfixed32Kind,
			protoreflect.Sfixed64Kind,
			protoreflect.Sint32Kind,
			protoreflect.Sint64Kind:
			fieldSchema.Type = bigquery.IntegerFieldType
		case protoreflect.BoolKind:
			fieldSchema.Type = bigquery.BooleanFieldType
		case protoreflect.StringKind:
			fieldSchema.Type = bigquery.StringFieldType
		case protoreflect.GroupKind:
			continue // ignore legacy proto2 group fields
		case protoreflect.BytesKind:
			fieldSchema.Type = bigquery.BytesFieldType
		case protoreflect.EnumKind:
			fieldSchema.Type = bigquery.StringFieldType
		case protoreflect.MessageKind:
			switch field.Message().FullName() {
			case "google.protobuf.Timestamp":
				fieldSchema.Type = bigquery.TimestampFieldType
			case "google.protobuf.Duration":
				fieldSchema.Type = bigquery.FloatFieldType
			case "google.protobuf.DoubleValue",
				"google.protobuf.FloatValue":
				fieldSchema.Type = bigquery.FloatFieldType
			case "google.protobuf.Int32Value",
				"google.protobuf.Int64Value",
				"google.protobuf.UInt32Value",
				"google.protobuf.UInt64Value":
				fieldSchema.Type = bigquery.IntegerFieldType
			case "google.protobuf.BoolValue":
				fieldSchema.Type = bigquery.BooleanFieldType
			case "google.protobuf.StringValue":
				fieldSchema.Type = bigquery.StringFieldType
			case "google.protobuf.BytesValue":
				fieldSchema.Type = bigquery.BytesFieldType
			case "google.protobuf.StructValue":
				fieldSchema.Type = bigquery.StringFieldType // JSON string
			default:
				fieldSchema.Type = bigquery.RecordFieldType
				fieldSchema.Schema = o.inferMessageSchema(field.Message())
			}
		}
		schema = append(schema, fieldSchema)
	}
	return schema
}