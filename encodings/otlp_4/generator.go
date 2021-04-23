package otlp_4

import (
	"math/rand"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/tigrannajaryan/exp-otelproto/core"
)

// Generator allows to generate a ExportRequest.
type Generator struct {
	random     *rand.Rand
	tracesSent uint64
	spansSent  uint64
}

func NewGenerator() *Generator {
	return &Generator{
		random: rand.New(rand.NewSource(99)),
	}
}

func (g *Generator) genRandByteString(len int) string {
	b := make([]byte, len)
	for i := range b {
		b[i] = byte(g.random.Intn(10) + 33)
	}
	return string(b)
}

func GenResource() *Resource {
	return &Resource{
		Attributes: []*KeyValue{
			{Key: "StartTimeUnixnano", Value: &AnyValue{Value: &AnyValue_IntValue{IntValue: 12345678}}},
			{Key: "Pid", Value: &AnyValue{Value: &AnyValue_IntValue{IntValue: 1234}}},
			{Key: "HostName", Value: &AnyValue{Value: &AnyValue_StringValue{StringValue: "fakehost"}}},
			{Key: "ServiceName", Value: &AnyValue{Value: &AnyValue_StringValue{StringValue: "generator"}}},
		},
	}
}

func (g *Generator) GenerateSpanBatch(spansPerBatch int, attrsPerSpan int, timedEventsPerSpan int) core.ExportRequest {
	traceID := atomic.AddUint64(&g.tracesSent, 1)
	batchStartTime := time.Date(2019, 10, 31, 10, 11, 12, 13, time.UTC)

	il := &InstrumentationLibrarySpans{
		InstrumentationLibrary: &InstrumentationLibrary{Name: "io.opentelemetry"},
	}
	batch := &ExportTraceServiceRequest{
		ResourceSpans: []*ResourceSpans{
			{
				Resource:                    GenResource(),
				InstrumentationLibrarySpans: []*InstrumentationLibrarySpans{il},
			},
		},
	}

	for i := 0; i < spansPerBatch; i++ {
		startTime := batchStartTime.Add(time.Duration(i) * time.Millisecond)

		spanID := atomic.AddUint64(&g.spansSent, 1)

		// Create a span.
		span := &Span{
			TraceId:           core.GenerateTraceID(traceID),
			SpanId:            core.GenerateSpanID(spanID),
			Name:              "load-generator-span",
			Kind:              Span_CLIENT,
			StartTimeUnixNano: core.TimeToTimestamp(startTime),
			EndTimeUnixNano:   core.TimeToTimestamp(startTime.Add(time.Duration(i) * time.Millisecond)),
		}

		if attrsPerSpan >= 0 {
			// Append attributes.
			span.Attributes = []*KeyValue{}

			if attrsPerSpan >= 2 {
				span.Attributes = append(span.Attributes,
					&KeyValue{
						Key:   "load_generator.span_seq_num",
						Value: &AnyValue{Value: &AnyValue_IntValue{IntValue: int64(spanID)}},
					})
				span.Attributes = append(span.Attributes,
					&KeyValue{
						Key:   "load_generator.trace_seq_num",
						Value: &AnyValue{Value: &AnyValue_IntValue{IntValue: int64(traceID)}},
					})
			}

			for j := len(span.Attributes); j < attrsPerSpan; j++ {
				attrName := core.GenRandAttrName(g.random)
				span.Attributes = append(span.Attributes,
					&KeyValue{
						Key:   attrName,
						Value: &AnyValue{Value: &AnyValue_StringValue{StringValue: g.genRandByteString(g.random.Intn(20) + 1)}},
					})
			}
		}

		if timedEventsPerSpan > 0 {
			for i := 0; i < timedEventsPerSpan; i++ {
				span.Events = append(span.Events, &Span_Event{
					TimeUnixNano: core.TimeToTimestamp(startTime.Add(time.Duration(i) * time.Millisecond)),
					// TimeStartDeltaNano: (time.Duration(i) * time.Millisecond).Nanoseconds(),
					Attributes: []*KeyValue{
						{Key: "te", Value: &AnyValue{Value: &AnyValue_IntValue{IntValue: int64(spanID)}}},
					},
				})
			}
		}

		il.Spans = append(il.Spans, span)
	}
	return batch
}

func (g *Generator) GenerateLogBatch(logsPerBatch int, attrsPerLog int) core.ExportRequest {
	return nil
}

func GenInt64Timeseries(startTime time.Time, offset int, valuesPerTimeseries int) []*Int64DataPoint {
	var timeseries []*Int64DataPoint
	for j := 0; j < 5; j++ {
		var points []*Int64DataPoint

		for k := 0; k < valuesPerTimeseries; k++ {
			pointTs := core.TimeToTimestamp(startTime.Add(time.Duration(j*k) * time.Millisecond))

			point := Int64DataPoint{
				TimeUnixNano: pointTs,
				Value:        int64(offset * j * k),
				Labels: []*StringKeyValue{
					{
						Key:   "label1",
						Value: "val1",
					},
					{
						Key:   "label2",
						Value: "val2",
					},
				},
			}

			if k == 0 {
				point.StartTimeUnixNano = pointTs
			}

			points = append(points, &point)
		}

		timeseries = append(timeseries, points...)
	}

	return timeseries
}

func genInt64Gauge(startTime time.Time, i int, valuesPerTimeseries int) *Metric {
	descr := GenMetricDescriptor(i)
	metric1 := &Metric{
		MetricDescriptor: descr,
		Int64DataPoints:  GenInt64Timeseries(startTime, i, valuesPerTimeseries),
	}
	return metric1
}

func GenMetricDescriptor(i int) *MetricDescriptor {
	descr := &MetricDescriptor{
		Name:        "metric" + strconv.Itoa(i),
		Description: "some description: " + strconv.Itoa(i),
		Type:        MetricDescriptor_INT64,
	}
	return descr
}

func genHistogram(startTime time.Time, i int, valuesPerTimeseries int) *Metric {
	// Add Histogram
	descr := GenMetricDescriptor(i)
	descr.Type = MetricDescriptor_INT64

	var timeseries2 []*HistogramDataPoint
	for j := 0; j < 1; j++ {
		var points []*HistogramDataPoint

		//prevPointTs := int64(0)
		for k := 0; k < valuesPerTimeseries; k++ {
			pointTs := core.TimeToTimestamp(startTime.Add(time.Duration(j*k) * time.Millisecond))
			//diffTs := pointTs - prevPointTs
			//prevPointTs = pointTs
			val := float64(i * j * k)
			point := HistogramDataPoint{
				TimeUnixNano: pointTs,
				Count:        1,
				Sum:          val,
				Buckets: []*HistogramDataPoint_Bucket{
					{
						Count: 12,
						Exemplar: &HistogramDataPoint_Bucket_Exemplar{
							Value:        val,
							TimeUnixNano: pointTs,
						},
					},
					{
						Count: 345,
					},
				},
				ExplicitBounds: []float64{0, 1000000},
				Labels: []*StringKeyValue{
					{
						Key:   "label1",
						Value: "val1",
					},
					{
						Key:   "label2",
						Value: "val2",
					},
				},
			}
			if k == 0 {
				point.StartTimeUnixNano = pointTs
			}
			points = append(points, &point)
		}

		timeseries2 = append(timeseries2, points...)
	}

	metric2 := &Metric{
		MetricDescriptor:    descr,
		HistogramDataPoints: timeseries2,
	}

	return metric2
}

func genSummary(startTime time.Time, i int, valuesPerTimeseries int) *Metric {
	// Add Histogram
	descr := GenMetricDescriptor(i)

	var timeseries2 []*SummaryDataPoint
	for j := 0; j < 1; j++ {
		var points []*SummaryDataPoint

		for k := 0; k < valuesPerTimeseries; k++ {
			pointTs := core.TimeToTimestamp(startTime.Add(time.Duration(j*k) * time.Millisecond))
			val := float64(i * j * k)
			point := SummaryDataPoint{
				TimeUnixNano: pointTs,
				Count:        1,
				Sum:          val,
				PercentileValues: []*SummaryDataPoint_ValueAtPercentile{
					{
						Percentile: 99,
						Value:      val / 10,
					},
				},
				Labels: []*StringKeyValue{
					{
						Key:   "label1",
						Value: "val1",
					},
					{
						Key:   "label2",
						Value: "val2",
					},
				},
			}
			if k == 0 {
				point.StartTimeUnixNano = pointTs
			}
			points = append(points, &point)
		}

		timeseries2 = append(timeseries2, points...)
	}

	metric2 := &Metric{
		MetricDescriptor:  descr,
		SummaryDataPoints: timeseries2,
	}
	return metric2
}

func (g *Generator) GenerateMetricBatch(
	metricsPerBatch int,
	valuesPerTimeseries int,
	int64 bool,
	histogram bool,
	summary bool,
) core.ExportRequest {

	batchStartTime := time.Date(2019, 10, 31, 10, 11, 12, 13, time.UTC)

	il := &InstrumentationLibraryMetrics{}
	batch := &ExportMetricsServiceRequest{
		ResourceMetrics: []*ResourceMetrics{
			{
				Resource:                      GenResource(),
				InstrumentationLibraryMetrics: []*InstrumentationLibraryMetrics{il},
			},
		},
	}

	for i := 0; i < metricsPerBatch; i++ {
		startTime := batchStartTime.Add(time.Duration(i) * time.Millisecond)

		if int64 {
			il.Metrics = append(il.Metrics, genInt64Gauge(startTime, i, valuesPerTimeseries))
		}
		if histogram {
			il.Metrics = append(il.Metrics, genHistogram(startTime, i, valuesPerTimeseries))
		}
		if summary {
			il.Metrics = append(il.Metrics, genSummary(startTime, i, valuesPerTimeseries))
		}
	}
	return batch
}

func timeToTimestamp(t time.Time) *timestamp.Timestamp {
	nanoTime := t.UnixNano()
	return &timestamp.Timestamp{
		Seconds: nanoTime / 1e9,
		Nanos:   int32(nanoTime % 1e9),
	}
}
