package progress

import "github.com/schollz/progressbar/v3"

func New(title string) *progressbar.ProgressBar {
	return progressbar.NewOptions(
		-1,
		progressbar.OptionSetDescription(title),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
	)
}
