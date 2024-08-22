/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2024 by  the Jacobin Authors. All rights reserved.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0)  Consult jacobin.org.
 */

package config

import (
	"fmt"
	"io/ioutil"
	"os"
)

var JacobinVersion = "0.5.001"

func GetJacobinVersion() string {
	file, err := os.Open("BUILDNO.txt")
	if err != nil {
		fmt.Println(err)
		return JacobinVersion
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return JacobinVersion
	}

	return fmt.Sprintf("%s Build %s", JacobinVersion, string(byteValue))
}
