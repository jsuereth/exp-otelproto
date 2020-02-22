package internal

import (
	"testing"

	"github.com/tigrannajaryan/exp-otelproto/encodings/experimental"

	"github.com/golang/protobuf/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tigrannajaryan/exp-otelproto/encodings/otlp"
)

const BatchCount = 1000

func BenchmarkFromOtlpToInternal(b *testing.B) {
	b.StopTimer()
	g := otlp.NewGenerator()

	var batch []*otlp.TraceExportRequest
	for i := 0; i < BatchCount; i++ {
		batch = append(batch,
			g.GenerateSpanBatch(100, 5, 0).(*otlp.TraceExportRequest))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < BatchCount; i++ {
			FromOtlp(batch[i])
		}
	}
}

func TestMarshalMap(t *testing.T) {
	r := &experimental.Resource{Labels: map[string]*experimental.AttributeKeyValue{
		"label1": &experimental.AttributeKeyValue{StringValue: "val1"},
	}}
	_, err := proto.Marshal(r)
	require.NoError(t, err)
	//log.Printf("%x", b)

	r2 := &Resource{Labels: map[string]*otlp.AttributeKeyValue{
		"label1": &otlp.AttributeKeyValue{StringValue: "val1"},
	}}
	buf := proto.NewBuffer([]byte{})
	err = MarshalResource(buf, r2)
	require.NoError(t, err)
	//log.Printf("%x", buf.Bytes())
}

func TestDecodeFromProto(t *testing.T) {
	rs := &otlp.TraceExportRequest{
		ResourceSpans: []*otlp.ResourceSpans{
			{
				Resource: &otlp.Resource{},
				Spans: []*otlp.Span{
					{
						Name: "spanA",
					},
					{},
				},
			},
		},
	}
	b, err := proto.Marshal(rs)
	//fmt.Printf("%x", b)
	require.NoError(t, err)
	buf := proto.NewBuffer(b)

	_, err = ResourceSpansFromBuf(buf, 0)
	assert.NoError(t, err)
}

func TestEncodeFromProto(t *testing.T) {
	tes := &TraceExportRequest{
		ResourceSpans: []*ResourceSpans{
			{
				Resource: &Resource{
					Labels: map[string]*otlp.AttributeKeyValue{
						"label1": {StringValue: "val1"},
						"label2": {StringValue: "val2"},
						"label3": {StringValue: "val3"},
						"label4": {StringValue: "val4"},
					},
				},
				Spans: []*Span{
					{
						Name: "spanA",
						Attributes: map[string]*otlp.AttributeKeyValue{
							"attribute1": {StringValue: "value1"},
							"attribute2": {StringValue: "value2"},
							"attribute3": {StringValue: "value3"},
							"attribute4": {StringValue: "value4"},
							"attribute5": {StringValue: "value5"},
						},
					},
					{},
				},
			},
		},
	}
	g := otlp.NewGenerator()
	otlpTes := g.GenerateSpanBatch(100, 3, 0).(*otlp.TraceExportRequest)
	//b1, err := proto.Marshal(otlpTes)
	//log.Printf("%x", b1)
	//require.NoError(t, err)
	//
	tes = FromOtlp(otlpTes)

	b2, err := Marshal(tes)
	//log.Printf("%x", b2)
	require.NoError(t, err)

	var tes2 experimental.TraceExportRequest
	err = proto.Unmarshal(b2, &tes2)
	assert.NoError(t, err)
}
