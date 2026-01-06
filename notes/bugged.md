
** Tree View Object (TVO)

When you are in the GoLand debugger and you have broke at a source line where one of the variables is obj, defined as *object.Object.
If you would like to display the entire object, execute the following:
- Position into the "Evaluate expression" windows.
- Type in this: bugged.TVO(obj)
- Press return. This yields a new varible "result" set = a long complex Go string with embedded newline characters.
- Right-click on "result" and copy its value.
- Alt-tab to a text editor and paste into a blank window.

** Convert a Java String object to a viewable Go string

When you are in the GoLand debugger and you have broke at a source line where one of the variables is jstr, defined as *object.Object but you know it's a Java String object.
If you would like to display the Go string equivalent, execute the following:
- Position into the "Evaluate expression" windows.
- Type in this: bugged.STR(jstr)
- Press return. This yields a new varible "result" with the expected result.

