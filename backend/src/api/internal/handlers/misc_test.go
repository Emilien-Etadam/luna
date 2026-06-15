package handlers

import (
	"luna-backend/types"
	"testing"
)

func TestResolveCaldavResourceUrl(t *testing.T) {
	base, err := types.NewUrl("http://192.168.30.124/dav.php/")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		href string
		want string
	}{
		{
			href: "/dav.php/calendars/projest/",
			want: "http://192.168.30.124/dav.php/calendars/projest/",
		},
		{
			href: "http://example.com/dav.php/calendars/user/",
			want: "http://example.com/dav.php/calendars/user/",
		},
		{
			href: "",
			want: "http://192.168.30.124/dav.php/",
		},
	}

	for _, tc := range tests {
		got := resolveCaldavResourceUrl(base, tc.href)
		if got != tc.want {
			t.Fatalf("resolveCaldavResourceUrl(%q) = %q, want %q", tc.href, got, tc.want)
		}
	}
}
