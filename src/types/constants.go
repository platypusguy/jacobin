/*
 * Jacobin VM - A Java virtual machine
 * Copyright (c) 2023 by  the Jacobin authors. Consult jacobin.org.
 * Licensed under Mozilla Public License 2.0 (MPL 2.0) All rights reserved.
 */

package types

// Grab bag of constants used in Jacobin

// ---- <clInit> status bytes ----
const NoClinit byte = 0x00
const ClInitNotRun byte = 0x01
const ClInitInProgress byte = 0x02
const ClInitRun byte = 0x03
