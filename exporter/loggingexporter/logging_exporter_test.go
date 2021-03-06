// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package loggingexporter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/internal/data/testdata"
)

func TestLoggingTraceExporterNoErrors(t *testing.T) {
	lte, err := newTraceExporter(&configmodels.ExporterSettings{}, "debug", zap.NewNop())
	require.NotNil(t, lte)
	assert.NoError(t, err)

	assert.NoError(t, lte.ConsumeTraces(context.Background(), testdata.GenerateTraceDataEmpty()))
	assert.NoError(t, lte.ConsumeTraces(context.Background(), testdata.GenerateTraceDataOneEmptyOneNilResourceSpans()))
	assert.NoError(t, lte.ConsumeTraces(context.Background(), testdata.GenerateTraceDataOneEmptyOneNilInstrumentationLibrary()))
	assert.NoError(t, lte.ConsumeTraces(context.Background(), testdata.GenerateTraceDataOneSpanOneNil()))
	assert.NoError(t, lte.ConsumeTraces(context.Background(), testdata.GenerateTraceDataTwoSpansSameResourceOneDifferent()))

	assert.NoError(t, lte.Shutdown(context.Background()))
}

func TestLoggingMetricsExporterNoErrors(t *testing.T) {
	lme, err := newMetricsExporter(&configmodels.ExporterSettings{}, "debug", zap.NewNop())
	require.NotNil(t, lme)
	assert.NoError(t, err)

	assert.NoError(t, lme.ConsumeMetrics(context.Background(), testdata.GenerateMetricsEmpty()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), testdata.GenerateMetricsOneEmptyOneNilResourceMetrics()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), testdata.GenerateMetricsOneEmptyOneNilInstrumentationLibrary()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), testdata.GenerateMetricsOneMetricOneNil()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), testdata.GenerateMetricsWithCountersHistograms()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), testdata.GenerateMetricsAllTypesNilDataPoint()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), testdata.GenerateMetricsAllTypesEmptyDataPoint()))
	assert.NoError(t, lme.ConsumeMetrics(context.Background(), testdata.GenerateMetricsMetricTypeInvalid()))

	assert.NoError(t, lme.Shutdown(context.Background()))
}

func TestLoggingLogsExporterNoErrors(t *testing.T) {
	lle, err := newLogsExporter(&configmodels.ExporterSettings{}, "debug", zap.NewNop())
	require.NotNil(t, lle)
	assert.NoError(t, err)

	assert.NoError(t, lle.ConsumeLogs(context.Background(), testdata.GenerateLogDataEmpty()))
	assert.NoError(t, lle.ConsumeLogs(context.Background(), testdata.GenerateLogDataOneEmptyResourceLogs()))
	assert.NoError(t, lle.ConsumeLogs(context.Background(), testdata.GenerateLogDataOneEmptyOneNilResourceLogs()))
	assert.NoError(t, lle.ConsumeLogs(context.Background(), testdata.GenerateLogDataNoLogRecords()))
	assert.NoError(t, lle.ConsumeLogs(context.Background(), testdata.GenerateLogDataOneEmptyLogs()))

	assert.NoError(t, lle.Shutdown(context.Background()))
}

func TestNestedArraySerializesCorrectly(t *testing.T) {
	av := pdata.NewAnyValueArray()
	ava := pdata.NewAttributeValueArray()
	av.Append(pdata.NewAttributeValueString("foo"))
	av.Append(pdata.NewAttributeValueInt(42))

	av2 := pdata.NewAnyValueArray()
	av2.Append(pdata.NewAttributeValueString("bar"))
	ava2 := pdata.NewAttributeValueArray()
	ava2.SetArrayVal(av2)

	av.Append(ava2)
	ava.SetArrayVal(av)

	assert.Equal(t, 3, ava.ArrayVal().Len())
	assert.Equal(t, "[foo, 42, [bar]]", attributeValueToString(ava))
}
