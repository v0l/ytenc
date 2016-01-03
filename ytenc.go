package main

import (
	"fmt"
	"os"
)

var (
	bs = 20
	w  = 1920
	h  = 1080
	r  = 5
)

type COLOUR struct {
	r, g, b float32
}

//http://stackoverflow.com/questions/7706339/grayscale-to-red-green-blue-matlab-jet-color-scale
func GetColour(v float32, vmin float32, vmax float32) COLOUR {
	c := COLOUR{1.0, 1.0, 1.0}
	var dv float32

	if v < vmin {
		v = vmin
	}
	if v > vmax {
		v = vmax
	}
	dv = vmax - vmin

	if v < (vmin + 0.25*dv) {
		c.r = 0.0
		c.g = 4.0 * (v - vmin) / dv
	} else if v < (vmin + 0.5*dv) {
		c.r = 0.0
		c.b = 1.0 + 4.0*(vmin+0.25*dv-v)/dv
	} else if v < (vmin + 0.75*dv) {
		c.r = 4.0 * (v - vmin - 0.5*dv) / dv
		c.b = 0.0
	} else {
		c.g = 1.0 + 4.0*(vmin+0.75*dv-v)/dv
		c.b = 0.0
	}

	return c
}

func MakeFrame(data []byte) []byte {
	ret := make([]byte, w*h*3)

	if len(data) == (w/bs)*(h/bs) {
		for i := 0; i < len(data); i++ {
			rgb := GetColour(float32(data[i])/255.0, 0.0, 1.0)

			for y := 0; y < bs; y++ {
				for x := 0; x < bs; x++ {
					pos := ((w * bs) * (i / (w / bs))) + (bs * (i % (w / bs))) + x + w*y //pos is pixel position

					zpos := pos * 3
					ret[zpos] = byte(255 * rgb.r)
					ret[zpos+1] = byte(255 * rgb.g)
					ret[zpos+2] = byte(255 * rgb.b)
				}
			}
		}
	} else {
		fmt.Println("Data is too long!")
	}

	return ret
}

func main() {
	fmt.Println("ytenc © anon 2016")

	file, err := os.Open("test.raw")
	fout, erro := os.Create("test_out.rgb")
	if err != nil || erro != nil {
		fmt.Println("Something broke :/")
	} else {
		data := make([]byte, (w/bs)*(h/bs))

		for {
			_, err := file.Read(data)
			if err != nil {
				break
			}

			f := MakeFrame(data)
			for z := 0; z < r; z++ {
				fout.Write(f)
			}
			fout.Sync()
		}
	}
}
