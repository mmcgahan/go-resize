Go-Resize
===========

A simple image resizing server using
[SmartCrop](https://github.com/muesli/smartcrop) and
[Imaging](https://github.com/disintegration/imaging).

This is a pure Go solution, which means it's not as performant as some
other C-based libraries like ImageMagick and OpenCV - it's pretty simple
to switch out the crop/resize libraries, but I couldn't use external
libraries on Heroku.

## API

Just put the desired width, height, and source image URL in the url,
e.g.:

```
http://localhost/200/400/www.example.com/theImage.jpg
```

**Limitations**

1. https sources _only_
2. Jpegs _only_

**TODO**

- reduce image.Image copying
- reference implementations of different imaging libraries
- performance testing
- manual crop position control (not just SmartCrop)

