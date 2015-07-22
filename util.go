package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	image_draw "image/draw"
	_ "image/png"
	"io/ioutil"
	"log"

	mgl "github.com/go-gl/mathgl/mgl32"

	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/gl"
)

type obj struct {
	slice []float32
	dim   int
}

func (o *obj) Slice() []float32 {
	return o.slice
}

func (o *obj) Dim() int {
	return o.dim
}

func NewObjectData(slice []float32, dim int) ObjectData {
	return &obj{
		slice: slice,
		dim:   dim,
	}
}

type ObjectData interface {
	Slice() []float32
	Dim() int
}

// EncodeObjects converts float32 vertices into a LittleEndian byte array.
// Offset and length are based on the number of rows per dimension.
func EncodeObjects(offset int, length int, objects ...ObjectData) []byte {
	log.Println("EncodeObjects:", offset, length, objects)

	// TODO: Pre-allocate?
	/*
		dimSum := 0 // yum!
		for _, obj := range objects {
			dimSum += obj.Dim()
		}
		v := make([]float32, dimSum*length)
	*/

	buf := bytes.Buffer{}

	for i := offset; i < length; i++ {
		v := []float32{}
		for _, obj := range objects {
			v = append(v, obj.Slice()[i*obj.Dim():(i+1)*obj.Dim()]...)
		}

		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			panic(fmt.Sprintln("binary.Write failed:", err))
		}
	}

	//fmt.Printf("Wrote %d vertices: %d to %d \t", shape.Len()-n, n, shape.Len())
	//fmt.Println(wrote)

	return buf.Bytes()
}

func loadAsset(name string) ([]byte, error) {
	f, err := asset.Open(name)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

// LoadProgram reads shader sources from the asset repository, compiles, and
// links them into a program.
func LoadProgram(vertexAsset, fragmentAsset string) (p gl.Program, err error) {
	log.Println("LoadProgram:", vertexAsset, fragmentAsset)

	vertexSrc, err := loadAsset(vertexAsset)
	if err != nil {
		return
	}

	fragmentSrc, err := loadAsset(fragmentAsset)
	if err != nil {
		return
	}

	p, err = glutil.CreateProgram(string(vertexSrc), string(fragmentSrc))
	return
}

// LoadTexture2D reads and decodes an image from the asset repository and creates
// a texture object based on the full dimensions of the image.
func LoadTexture2D(name string) (tex gl.Texture, err error) {
	imgFile, err := asset.Open(name)
	if err != nil {
		return
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return
	}

	rgba := image.NewRGBA(img.Bounds())
	image_draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, image_draw.Src)

	tex = gl.CreateTexture()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		rgba.Rect.Size().X,
		rgba.Rect.Size().Y,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		rgba.Pix)
	return
}

// LoadTextureCube reads and decodes an image from the asset repository and creates
// a texture cube map object based on the full dimensions of the image.
func LoadTextureCube(name string) (tex gl.Texture, err error) {
	imgFile, err := asset.Open(name)
	if err != nil {
		return
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return
	}

	rgba := image.NewRGBA(img.Bounds())
	image_draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, image_draw.Src)

	tex = gl.CreateTexture()
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_CUBE_MAP, tex)

	target := gl.TEXTURE_CUBE_MAP_POSITIVE_X
	for i := 0; i < 6; i++ {
		// TODO: Load atlas, not the same image
		gl.TexImage2D(
			gl.Enum(target+i),
			0,
			rgba.Rect.Size().X,
			rgba.Rect.Size().Y,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			rgba.Pix,
		)
	}

	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	// Not available in GLES 2.0 :(
	//gl.TexParameteri(gl.TEXTURE_CUBE_MAP, gl.TEXTURE_WRAP_R, gl.CLAMP_TO_EDGE)

	return
}

func Quad(a mgl.Vec3, b mgl.Vec3) []float32 {
	return []float32{
		// First triangle
		b.X(), b.Y(), b.Z(), // Top Right
		a.X(), b.Y(), a.Z(), // Top Left
		a.X(), a.Y(), a.Z(), // Bottom Left
		// Second triangle
		b.X(), a.Y(), b.Z(), // Bottom Right
		b.X(), b.Y(), b.Z(), // Top Right
		a.X(), a.Y(), a.Z(), // Bottom Left
	}
}
