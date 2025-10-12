package ui

type MenuUI interface {
	Select(title string, options []string) (index int, err error)
}

type Fake struct {
	Choices []int
	Err     error
	i       int
}

func (f *Fake) Select(_ string, _ []string) (int, error) {
	if f.Err != nil {
		return 0, f.Err
	}
	if f.i >= len(f.Choices) {
		return 0, nil
	}
	v := f.Choices[f.i]
	f.i++
	return v, nil
}
