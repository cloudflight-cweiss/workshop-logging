package workshop

import "testing"

func TestClient(t *testing.T) {
	c := NewClient("localhost:8080")

	if err := c.Login(); err != nil {
		t.Fatal(err)
	}
	if err := c.Logout(); err != nil {
		t.Fatal(err)
	}
	if err := c.GetProject(); err != nil {
		t.Fatal(err)
	}
	if err := c.UpdateProject(); err != nil {
		t.Fatal(err)
	}
}
