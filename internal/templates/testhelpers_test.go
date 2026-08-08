package templates

import "testing"

func requireErr(t *testing.T, err error, wantErr bool) {
	t.Helper()
	if wantErr {
		if err == nil {
			t.Fatal("expected error")
		}

		return
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
