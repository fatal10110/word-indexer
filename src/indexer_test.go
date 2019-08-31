package main

import (
	"errors"
	"io"
	"reflect"
	"testing"
)

type MockGenerator struct {
	res   [][]byte
	err   error
	index int
}

func (gen *MockGenerator) Next() ([]byte, error) {
	if len(gen.res) == 0 {
		return []byte{}, gen.err
	}

	word := gen.res[gen.index]
	gen.index++

	return word, gen.err
}

func Test_wordIndexer_Index(t *testing.T) {
	type fields struct {
		gen BatchesGenerator
	}
	tests := []struct {
		name    string
		fields  fields
		want    IndexResults
		wantErr bool
	}{
		{name: "Test no words", fields: fields{gen: &MockGenerator{}}, want: IndexResults{}, wantErr: false},
		{name: "Test no words with err", fields: fields{gen: &MockGenerator{err: errors.New("")}}, want: IndexResults{}, wantErr: true},
		{name: "Test no words with eof", fields: fields{gen: &MockGenerator{err: io.EOF}}, want: IndexResults{}, wantErr: false},
		{name: "Test index one word", fields: fields{gen: &MockGenerator{err: io.EOF, res: [][]byte{[]byte("123")}}}, want: IndexResults{"123": 1}, wantErr: false},
		{name: "Test index multiple words", fields: fields{gen: &MockGenerator{res: [][]byte{[]byte("123"), []byte("456"), []byte("")}}}, want: IndexResults{"123": 1, "456": 1}, wantErr: false},
		{name: "Test index multiple times", fields: fields{gen: &MockGenerator{res: [][]byte{[]byte("456"), []byte("123"), []byte("456"), []byte("")}}}, want: IndexResults{"123": 1, "456": 2}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wi := &wordIndexer{
				gen: tt.fields.gen,
			}
			got, err := wi.Index()
			if (err != nil) != tt.wantErr {
				t.Errorf("wordIndexer.Index() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wordIndexer.Index() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexResults_Add(t *testing.T) {
	type args struct {
		word  string
		value int
	}
	tests := []struct {
		name string
		ir   IndexResults
		args args
		want IndexResults
	}{
		{"Creates stats", IndexResults{}, args{"some", 199}, IndexResults{"some": 199}},
		{"Adds case insensetive", IndexResults{"some": 199}, args{"sOme", 1}, IndexResults{"some": 200}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.ir.Add(tt.args.word, tt.args.value)

			if !reflect.DeepEqual(tt.ir, tt.want) {
				t.Errorf("IndexResults.Add() = %v, want %v", tt.ir, tt.want)
			}
		})
	}
}

func TestIndexResults_Get(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name  string
		ir    IndexResults
		args  args
		want  int
		want1 bool
	}{
		{"Get when no key", IndexResults{}, args{"some"}, 0, false},
		{"Get case insensetive key", IndexResults{"some": 2}, args{"soMe"}, 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.ir.Get(tt.args.word)
			if got != tt.want {
				t.Errorf("IndexResults.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("IndexResults.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
