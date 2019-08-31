package main

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_wordsSplitGenerator_Next(t *testing.T) {
	type fields struct {
		input 	  []byte
		delimiter byte
	}
	tests := []struct {
		name     string
		fields   fields
		wantWord [][]byte
		wantErr  []bool
	}{
		{name: "empty batch", fields: fields{input: []byte(""), delimiter: byte('-')}, wantWord: [][]byte{[]byte("")}, wantErr: []bool{true}},
		{name: "single word batch", fields: fields{input: []byte("asd"), delimiter: byte('-')}, wantWord: [][]byte{[]byte("asd")}, wantErr: []bool{true}},
		{name: "single word with delimeter suffix", fields: fields{input: []byte("asd-"), delimiter: byte('-')}, wantWord: [][]byte{[]byte("asd"), []byte("")}, wantErr: []bool{false, true}},
		{name: "single word with delimeter prefix", fields: fields{input: []byte("-asd"), delimiter: byte('-')}, wantWord: [][]byte{[]byte("asd")}, wantErr: []bool{true}},
		{name: "single word with multiple delimeter prefix", fields: fields{input: []byte("---asd"), delimiter: byte('-')}, wantWord: [][]byte{[]byte("asd")}, wantErr: []bool{true}},
		{name: "multiple words", fields: fields{input: []byte("123-456--789"), delimiter: byte('-')}, wantWord: [][]byte{[]byte("123"), []byte("456"), []byte("789")}, wantErr: []bool{false, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewWordSplitGenerator(bytes.NewBuffer(tt.fields.input), tt.fields.delimiter)
			var i int
			
			for ;; i++ {
				gotWord, err := gen.Next()

				wantErr := tt.wantErr[i]

				if (err != nil) != wantErr {
					t.Errorf("wordsSplitGenerator.Next() error = %v, wantErr %v", err, wantErr)
					return
				}

				want := tt.wantWord[i]

				
				if !(len(gotWord) == 0 && len(want) == 0) && !reflect.DeepEqual(gotWord, want) {
					t.Errorf("wordsSplitGenerator.Next() = %s want %s bytes: %v, want %v", string(gotWord), string(want), gotWord, want)
				}

				if err == io.EOF {
					break
				}
			}

			if i != len(tt.wantWord) - 1 || i != len(tt.wantErr) - 1 {
				t.Errorf("wordsSplitGenerator.Next() results number mismatch got %d expected results %d errors %d", i, len(tt.wantWord), len(tt.wantErr))
			}
		})
	}
}

func Test_bytesBatchGenerator_Next(t *testing.T) {
	type fields struct {
		input     []byte
		delimiter byte
		batchSize int64
	}
	tests := []struct {
		name    string
		fields  fields
		want    [][]byte
		wantErr []bool
	}{
		{name: "empty input", fields: fields{input: []byte(""), delimiter: byte('-'), batchSize: 2}, want: [][]byte{[]byte("")}, wantErr: []bool{true}},
		{name: "without delimiter", fields: fields{input: []byte("asd"), delimiter: byte('-'), batchSize: 2}, want: [][]byte{[]byte("asd")}, wantErr: []bool{true}},
		{name: "delimiter suffix", fields: fields{input: []byte("asd-"), delimiter: byte('-'), batchSize: 2}, want: [][]byte{[]byte("asd"), []byte("")}, wantErr: []bool{false, true}},
		{name: "multi delimiter suffix", fields: fields{input: []byte("asd---"), delimiter: byte('-'), batchSize: 2}, want: [][]byte{[]byte("asd"), []byte("")}, wantErr: []bool{false, true}},
		{name: "delimiter prefix", fields: fields{input: []byte("-asd"), delimiter: byte('-'), batchSize: 2}, want: [][]byte{[]byte("asd")}, wantErr: []bool{true}},
		{name: "multi delimiter prefix", fields: fields{input: []byte("---asd"), delimiter: byte('-'), batchSize: 2}, want: [][]byte{[]byte("asd")}, wantErr: []bool{true}},
		{name: "with delimiter", fields: fields{input: []byte("123-456-789"), delimiter: byte('-'), batchSize: 2}, want: [][]byte{[]byte("123"), []byte("456"), []byte("789")}, wantErr: []bool{false, false, true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := NewBatchGenerator(bytes.NewBuffer(tt.fields.input), tt.fields.delimiter, tt.fields.batchSize)
			var i int

			for ;; i++ {
				got, err := gen.Next()

				wantErr := tt.wantErr[i]

				if (err != nil) != wantErr {
					t.Errorf("bytesBatchGenerator.Next() error = %v, wantErr %v", err, wantErr)
					return
				}

				want := tt.want[i]

				if !(len(got) == 0 && len(want) == 0) && !reflect.DeepEqual(got, want) {
					t.Errorf("bytesBatchGenerator.Next() = %s, want %s as bytes %v want %v", string(got), string(want), got, want)
				}

				if err == io.EOF {
					break
				}
			}

			if i != len(tt.want) - 1 || i != len(tt.wantErr) - 1 {
				t.Errorf("bytesBatchGenerator.Next() results number mismatch got %d expected results %d errors %d", i, len(tt.want), len(tt.wantErr))
			}
		})
	}
}