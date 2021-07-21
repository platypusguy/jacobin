/*
 * jacobin - JVM written in Swift
 *
 * Copyright (c) 2021 Andrew Binstock. All rights reserved.
 * Licensed under Mozilla Public License, v. 2.0. http://mozilla.org/MPL/2.0/.
 */

import Foundation

/// The field attribute that holds the ConstantValue data of a field, if specified
/// The type can be: I (int), L (long), F (float), D (double), s (string, note: lower-case)
/// The value is an Any type that is then downcast when accessed, based on the type variable
class FieldInitAttribute {
    let name = "FieldInit"
    var type = ""
    var value : Any

    init( type: String, value: Any ) {
        self.type = type
        self.value = value
    }
}
