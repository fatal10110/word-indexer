package main

import "testing"

func Test_batchingJob_Execute(t *testing.T) {
	type fields struct {
		inputType InputType
		input     string
		batchSize int64
		broker    Broker
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &batchingJob{
				inputType: tt.fields.inputType,
				input:     tt.fields.input,
				batchSize: tt.fields.batchSize,
				broker:    tt.fields.broker,
			}
			if err := j.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("batchingJob.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_batchIndexerJob_Execute(t *testing.T) {
	type fields struct {
		batch  []byte
		broker Broker
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &batchIndexerJob{
				batch:  tt.fields.batch,
				broker: tt.fields.broker,
			}
			if err := j.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("batchIndexerJob.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_reportStatisticJob_Execute(t *testing.T) {
	type fields struct {
		statistic IndexResults
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &reportStatisticJob{
				statistic: tt.fields.statistic,
			}
			if err := j.Execute(); (err != nil) != tt.wantErr {
				t.Errorf("reportStatisticJob.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
