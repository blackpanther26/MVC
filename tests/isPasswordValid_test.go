package tests

import (
	"testing"

	"github.com/blackpanther26/mvc/pkg/controllers"  
)

func TestIsPasswordValid(t *testing.T)  {
	got := controllers.IsPasswordValid("Pixel$62738982")
	want := true
	if got != want {
		t.Errorf("got %v, wanted %v", got, want)
	}
}
