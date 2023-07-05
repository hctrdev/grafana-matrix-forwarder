package formatter

import (
	"grafana-matrix-forwarder/model"
	"testing"
)

func TestGenerateMessage(t *testing.T) {
	type args struct {
		alert          model.AlertData
		metricRounding int
	}
	tests := []struct {
		name                 string
		args                 args
		wantPlainMessage     string
		wantFormattedMessage string
		wantErr              bool
	}{
		{
			name: "alertingStateTest",
			args: args{
				alert: model.AlertData{
					State:    "alerting",
					RuleURL:  "http://example.com",
					RuleName: "sample",
					Message:  "sample message",
				}, metricRounding: 0},
			wantFormattedMessage: "💔 <b>ALERT</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p>",
			wantPlainMessage:     "💔 ALERT Rule: sample | sample message",
			wantErr:              false,
		},
		{
			name: "alertingStateWithEvalMatchesTest",
			args: args{
				alert: model.AlertData{
					State:    "alerting",
					RuleURL:  "http://example.com",
					RuleName: "sample",
					Message:  "sample message",
					EvalMatches: []struct {
						Value  float64
						Metric string
						Tags   map[string]string
					}{
						{
							Value:  10.65124,
							Metric: "sample",
							Tags:   map[string]string{},
						},
					},
				},
				metricRounding: 5,
			},
			wantFormattedMessage: "💔 <b>ALERT</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p><ul><li><b>sample</b>: 10.65124</li></ul>",
			wantPlainMessage:     "💔 ALERT Rule: sample | sample message sample: 10.65124",
			wantErr:              false,
		},
		{
			name: "alertingStateWithEvalMatchesAndTagsTest",
			args: args{
				alert: model.AlertData{
					State:    "alerting",
					RuleURL:  "http://example.com",
					RuleName: "sample",
					Message:  "sample message",
					EvalMatches: []struct {
						Value  float64
						Metric string
						Tags   map[string]string
					}{
						{
							Value:  10.65124,
							Metric: "sample",
						},
					},
					Tags: map[string]string{"key1": "value1", "key2": "value2"},
				},
				metricRounding: 5,
			},
			wantFormattedMessage: "💔 <b>ALERT</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p><ul><li><b>sample</b>: 10.65124</li></ul><p>Tags:</p><ul><li><b>key1</b>: value1</li><li><b>key2</b>: value2</li></ul>",
			wantPlainMessage:     "💔 ALERT Rule: sample | sample message sample: 10.65124 Tags: key1: value1key2: value2",
			wantErr:              false,
		},
		{
			name: "okStateTest",
			args: args{
				alert: model.AlertData{
					State:    "ok",
					RuleURL:  "http://example.com",
					RuleName: "sample",
					Message:  "sample message",
				},
			},
			wantFormattedMessage: "💚 <b>RESOLVED</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p>",
			wantPlainMessage:     "💚 RESOLVED Rule: sample | sample message",
			wantErr:              false,
		},
		{
			name: "noDataStateTest",
			args: args{
				alert: model.AlertData{
					State:    "no_data",
					RuleURL:  "http://example.com",
					RuleName: "sample",
					Message:  "sample message",
				},
			},
			wantFormattedMessage: "❓ <b>NO DATA</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p>",
			wantPlainMessage:     "❓ NO DATA Rule: sample | sample message",
			wantErr:              false,
		},
		{
			name: "unknownStateTest",
			args: args{
				alert: model.AlertData{
					State:    "invalid state",
					RuleURL:  "http://example.com",
					RuleName: "sample",
					Message:  "sample message",
				},
			},
			wantFormattedMessage: "❓ <b>UNKNOWN</b><p>Rule: <a href=\"http://example.com\">sample</a> | sample message</p>",
			wantPlainMessage:     "❓ UNKNOWN Rule: sample | sample message",
			wantErr:              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMessage, err := GenerateMessage(tt.args.alert, tt.args.metricRounding)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotMessage.TextBody != tt.wantPlainMessage {
				t.Errorf("GenerateMessage() gotPlainMessage = %v, want %v", gotMessage.TextBody, tt.wantPlainMessage)
			}
			if gotMessage.HtmlBody != tt.wantFormattedMessage {
				t.Errorf("GenerateMessage() gotFormattedMessage = %v, want %v", gotMessage.HtmlBody, tt.wantFormattedMessage)
			}
		})
	}
}
