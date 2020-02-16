package slug_test

import "fmt"

func Example() {
	for _, text := range []string{
		"JOHN Doe h@#$ $Saver žůžo this ---__ hahah",
		"Ths ias  $anso sthinr.     ",
		"  dsflskdfa -.x.cxv.>>>___",
		"影師",
	} {
		slug := slugify(text)
		fmt.Println(slug, len(slug))
	}
}
