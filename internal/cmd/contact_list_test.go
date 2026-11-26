package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestContactListsHavePagingFlags(t *testing.T) {
	for name, cmd := range map[string]*cobra.Command{
		"contact list":       contactListCmd,
		"contact-group list": contactGroupListCmd,
	} {
		for _, f := range []string{"page", "per-page"} {
			if cmd.Flags().Lookup(f) == nil {
				t.Errorf("%s に --%s が無い", name, f)
			}
		}
	}
}

func TestListQuerySendsPaging(t *testing.T) {
	defer func() { page, perPage = 0, 0 }()
	if err := contactListCmd.Flags().Parse([]string{"--page", "2", "--per-page", "100"}); err != nil {
		t.Fatal(err)
	}
	q := listQuery()
	if q.Get("page") != "2" || q.Get("per_page") != "100" {
		t.Errorf("query = %v, want page=2 per_page=100", q)
	}
}
