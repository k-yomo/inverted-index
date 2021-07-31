package index

import (
	"github.com/google/go-cmp/cmp"
	"reflect"
	"testing"
)

func TestNewIndex(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want *Index
	}{
		{
			name: "create inverted index",
			args: args{word: "AAA b&bb AAA A^Aa."},
			want: &Index{PostingMap: map[string][]int{
				"aaa": {0, 2, 3},
				"bbb": {1},
			}},
		},
		{
			name: "don't create postings for empty string(after normalized)",
			args: args{word: "     $   -  "},
			want: &Index{PostingMap: map[string][]int{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIndex(tt.args.word)
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("NewIndex() = (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIndex_First(t *testing.T) {
	t.Parallel()

	type args struct {
		term string
	}
	tests := []struct {
		name  string
		index *Index
		args  args
		want  int
	}{
		{
			name:  "return first position",
			index: NewIndex("abc test test"),
			args:  args{term: "test"},
			want:  1,
		},
		{
			name:  "return -∞ if the term is not found",
			index: NewIndex("abc test test"),
			args:  args{term: "zzz"},
			want:  NegativeInf,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.index.First(tt.args.term); got != tt.want {
				t.Errorf("First() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_Last(t *testing.T) {
	type args struct {
		term string
	}
	tests := []struct {
		name  string
		index *Index
		args  args
		want  int
	}{
		{
			name:  "return last position",
			index: NewIndex("abc test test"),
			args:  args{term: "test"},
			want:  2,
		},
		{
			name:  "return ∞ if the term is not found",
			index: NewIndex("abc test test"),
			args:  args{term: "zzz"},
			want:  Inf,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.index.Last(tt.args.term); got != tt.want {
				t.Errorf("Last() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_Prev(t *testing.T) {
	t.Parallel()

	type args struct {
		term    string
		current int
	}
	tests := []struct {
		name  string
		index *Index
		args  args
		want  int
	}{
		{
			name:  "return prev position",
			index: NewIndex("word word word"),
			args:  args{term: "word", current: 2},
			want:  1,
		},
		{
			name:  "return last when current is ∞",
			index: NewIndex("word word word"),
			args:  args{term: "word", current: Inf},
			want:  2,
		},
		{
			name:  "return -∞ when current is first",
			index: NewIndex("word word word"),
			args:  args{term: "word", current: 0},
			want:  NegativeInf,
		},
		{
			name:  "return -∞ when current is -∞",
			index: NewIndex("word word word"),
			args:  args{term: "word", current: NegativeInf},
			want:  NegativeInf,
		},
		{
			name:  "return -∞ when not found",
			index: NewIndex("word word word"),
			args:  args{term: "test", current: Inf},
			want:  NegativeInf,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.index.Prev(tt.args.term, tt.args.current); got != tt.want {
				t.Errorf("Prev() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_Next(t *testing.T) {
	t.Parallel()

	type args struct {
		term    string
		current int
	}
	tests := []struct {
		name  string
		index *Index
		args  args
		want  int
	}{
		{
			name:  "return next position",
			index: NewIndex("word word word"),
			args:  args{term: "word", current: 1},
			want:  2,
		},
		{
			name:  "return ∞ when current is ∞",
			index: NewIndex("word word word"),
			args:  args{term: "word", current: Inf},
			want:  Inf,
		},
		{
			name:  "return ∞ when current is last",
			index: NewIndex("word word word"),
			args:  args{term: "word", current: 2},
			want:  Inf,
		},
		{
			name:  "return first when current is -∞",
			index: NewIndex("word word word"),
			args:  args{term: "word", current: NegativeInf},
			want:  0,
		},
		{
			name:  "return ∞ when not found",
			index: NewIndex("word word word"),
			args:  args{term: "test", current: Inf},
			want:  Inf,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.index.Next(tt.args.term, tt.args.current); got != tt.want {
				t.Errorf("Next() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndex_NextPhrase(t *testing.T) {
	type args struct {
		phrase   string
		position int
	}
	tests := []struct {
		name  string
		index *Index
		args  args
		want  *Range
	}{
		{
			name:  "return first range of the phrase appearance when given position is -∞",
			index: NewIndex("test word test word"),
			args: args{
				phrase:   "test word",
				position: NegativeInf,
			},
			want: &Range{
				From: 0,
				To:   1,
			},
		},
		{
			name:  "return next range of the phrase appearance",
			index: NewIndex("test word test abc word test word"),
			args: args{
				phrase:   "test word",
				position: 1,
			},
			want: &Range{
				From: 5,
				To:   6,
			},
		},
		{
			name:  "return range of ∞ if the phrase is not found after the position",
			index: NewIndex("test word test word"),
			args: args{
				phrase:   "test abc",
				position: 1,
			},
			want: &Range{
				From: Inf,
				To:   Inf,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.index.NextPhrase(tt.args.phrase, tt.args.position); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NextPhrase() = %v, want %v", got, tt.want)
			}
		})
	}
}
