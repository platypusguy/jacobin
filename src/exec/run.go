/* Jacobin VM -- A Java virtual machine
 * © Copyright 2021 by Andrew Binstock. All rights reserved
 * Licensed under Mozilla Public License 2.0 (MPL-2.0)
 */

package exec

import "errors"

func StartExec(className string) error {
	m, err := fetchMethod(className, "main")
	if err != nil {
		return errors.New("Class not found: " + className + ".main()")
	}
	f := frame{}
	for i := 0; i < len(m.CodeAttr.Code); i++ {
		f.meth = append(f.meth, m.CodeAttr.Code[i])
	}

	t := CreateThread(0)
	pushFrame(&t.stack, f)

	return nil
}
