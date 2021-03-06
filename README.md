WIP
---
This project is a *work in progress*. The implementation is *incomplete* and
subject to change. The documentation can be inaccurate.

Konstig
======

Konstig is a trigonometric strange attractor generator. It's name is simply weird, or strange, in Swedish.

Installation
------------

`go get github.com/karlek/konstig/cmd/konstig`

Generate an image
-----------------

```shell
$ konstig -a 4.663414 -b 1.502477 -c 3.351885 -d 2.894832 -i 1000000 -z 400 -w 5000 -h 5000 -f1 0.005 -f2 0.005 -f3 0.005 -p1 0 -p2 2 -p3 8 -o out.png
```

Equation
--------

x<sub>n+1</sub> = sin(y\*b) + c\*sin(x\*b)  
y<sub>n+1</sub> = sin(x\*a) + d\*sin(y\*a)

Flags:
------

* __a, b, c, d:__  
    A-,B-,C-,D-coefficient, see formula above.
* __i:__  
    Number of iterations performed.
* __f1, f2, f3:__  
    Color frequency.
* __p1, p2, p3:__  
    Color phases.
* __z:__  
    Zoom value. How many magnifications to make on the center point.
* __o:__  
    Output filename.
* __w, h:__  
    Width and height of created image.

Examples
--------

Some pretty examples. 

a = 1.364325  
b = 0.441972  
c = 3.868827  
d = 4.301396  

![A strange attractor which looks like a paint brush](https://github.com/karlek/konstig/blob/master/examples/1.364325_0.441972_3.868827_4.301396.png?raw=true)

a = 0.860540  
b = 1.602272  
c = 4.690712  
d = 2.123729  

![A strange attractor looking like a blanket of silk](https://github.com/karlek/konstig/blob/master/examples/0.860540_1.602272_4.690712_2.123729.png?raw=true)

All images were created with:
```shell
$ konstig
```

Public domain
-------------
I hereby release this code into the [public domain](https://creativecommons.org/publicdomain/zero/1.0/).
