package charm

import "github.com/vegidio/mediaorient"

type spinnerDoneMsg struct {
	result []mediaorient.Media
	err    error
}
