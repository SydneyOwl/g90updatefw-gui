# g90updatefw-gui

## Introduction
This is the GUI version of [g90updatefw](https://github.com/DaleFarnsworth/g90updatefw), It uploads firmware updates to the Xiegu G90 and Xiegu G106 radios.

In other words it is same as g90updatefw but with a gui interface.

![](./assets/img.png)

## Usage
Really simple, Just:
+ Connect cable to your computer
+ Choose firmware and port
+ Click "start"

Should work on windows/linux/macOS but only tried on ubuntu 24.04

but don't forget to run it with `sudo -E` or `sudo` on Linux!

![](./assets/usagegif.gif)

## Build
Just run `go mod tidy && go build`.

On unix platform maybe you need `gcc` and `libgtk-3-dev` as well

## License
```text
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <https://unlicense.org>
```