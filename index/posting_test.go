package index

import "testing"

func Test_postings_NextDocIndex(t *testing.T) {
	type args struct {
		docID int
	}
	tests := []struct {
		name string
		ps   postings
		args args
		want int
	}{
		{
			ps:   postings{{docID: 1}, { docID: 2}},
			args: args{docID: NegativeInf},
			want: 0,
		},
		{
			ps:   postings{{docID: 1}, { docID: 2}},
			args: args{docID: 0},
			want: 0,
		},
		{
			ps:   postings{{docID: 1}, { docID: 2}},
			args: args{docID: 1},
			want: 1,
		},
		{
			ps:   postings{{docID: 1}, { docID: 3}},
			args: args{docID: 2},
			want: 1,
		},
		{
			ps:   postings{{docID: 1}, { docID: 3}},
			args: args{docID: 3},
			want: Inf,
		},
		{
			ps:   postings{{docID: 1}, { docID: 3}},
			args: args{docID: 4},
			want: Inf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ps.NextDocIndex(tt.args.docID); got != tt.want {
				t.Errorf("NextDocIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}
