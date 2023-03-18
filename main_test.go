package main

import (
	"testing"

	"github.com/shurcooL/githubv4"
)

func TestCheckAllSuccess(t *testing.T) {
	tests := []struct {
		name string
		args []Context
		want bool
	}{
		{
			name: "all_success",
			args: []Context{
				{
					CheckRun: struct {
						Conclusion githubv4.CheckConclusionState
						DetailsUrl githubv4.URI
						Name       string
					}{
						githubv4.CheckConclusionStateSuccess,
						githubv4.URI{},
						"check_run_1",
					},
					StatusContext: struct {
						State     githubv4.StatusState
						TargetUrl githubv4.URI
						Context   string
					}{
						githubv4.StatusStateSuccess,
						githubv4.URI{},
						"status_context_1",
					},
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkAllSuccess(tt.args)
			if err != nil {
				t.Fatalf("error occured: %v", err)
			}
			if got != tt.want {
				t.Fatalf("checkAllSuccess() =%vm want %v", got, tt.want)
			}
		})
	}
}
