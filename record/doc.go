// Copyright (C) 2014 Ronoaldo JLP <ronoaldo@gmail.com>
// Licensed under the terms of the Apache License 2.0

/*
Package record provides fixed width record encoding utilities.


Encoding

The Encoder type converts struct fields into fixed width records.
These records are commonly used for data transport
on legacy systems, or some financial institutions, to name a few.

The struct fields of type string and int
are encoded with left-padding by default,
respectivelly using spaces or zeroes.
All other types are ignored by the encoder.

The convenience function Marshal
encodes a value and returns the encoded record bytes.


Tags

The struct fields can have a comma separated list of options
in a struct tag named `record`.

The tag can start with a number, like `record:"1"`,
that determines the size of the field.
String fields are padded with spaces to fill up to size,
and are truncated if higher than size.
Int fields are padded with zeroes to fill up to size.

The tag can have a "nopadding" option,
that avoids zero or space padding,
and uses the raw value as in %s and %d to fmt.Printf.
Note that this can make the resulting record variable in length.

The tag can have a "upper" option,
that causes strings to be upper cased before encoding.

When the tag option "-" is present, the field is skipped.
*/
package record
