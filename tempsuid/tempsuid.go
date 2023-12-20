package tempsuid

import (
	"github.com/dop251/goja"
)

func Get() (string, error) {
	vm := goja.New()
	_, err := vm.RunString(`
	function q0() {
    return "xxxxxxxxxxxxxxxx".replace(/[xy]/g, (function(e) {
        const t = 16 * Math.random() | 0;
        return ("x" == e ? t : 3 & t | 8).toString(16)
    }
    ))
}
	`)
	if err != nil {
		return ``, err
	}
	var js func() string
	err = vm.ExportTo(vm.Get("q0"), &js)
	if err != nil {
		return ``, err
	}
	return js(), nil
}
