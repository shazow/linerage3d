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
	"golang.org/x/mobile/gl"
)

type dimslice_float32 struct {
	dim   int
	slice []float32
}

func (o dimslice_float32) Slice(i, j int) interface{} { return o.slice[i:j] }
func (o dimslice_float32) Dim() int                   { return o.dim }
func (o dimslice_float32) String() string {
	return fmt.Sprintf("<float32 slice: len=%d dim=%d>", len(o.slice), o.dim)
}

type dimslice_uint8 struct {
	dim   int
	slice []uint8
}

func (o dimslice_uint8) Slice(i, j int) interface{} { return o.slice[i:j] }
func (o dimslice_uint8) Dim() int                   { return o.dim }
func (o dimslice_uint8) String() string {
	return fmt.Sprintf("<uint8 slice: len=%d dim=%d>", len(o.slice), o.dim)
}

func NewDimSlice(dim int, slice interface{}) DimSlicer {
	switch slice := slice.(type) {
	case []float32:
		return &dimslice_float32{dim, slice}
	case []uint8:
		return &dimslice_uint8{dim, slice}
	}
	panic(fmt.Sprintf("invalid slice type: %T", slice))
	return nil
}

type DimSlicer interface {
	Slice(int, int) interface{}
	Dim() int
	String() string
}

// EncodeObjects converts float32 vertices into a LittleEndian byte array.
// Offset and length are based on the number of rows per dimension.
// TODO: Replace with https://github.com/lunixbochs/struc?
func EncodeObjects(offset int, length int, objects ...DimSlicer) []byte {
	//log.Println("EncodeObjects:", offset, length, objects)
	// TODO: Pre-allocate? Use a SyncPool?
	/*
		dimSum := 0 // yum!
		for _, obj := range objects {
			dimSum += obj.Dim()
		}
		v := make([]float32, dimSum*length)
	*/

	buf := bytes.Buffer{}

	for i := offset; i < length; i++ {
		for _, obj := range objects {
			data := obj.Slice(i*obj.Dim(), (i+1)*obj.Dim())
			if err := binary.Write(&buf, binary.LittleEndian, data); err != nil {
				panic(fmt.Sprintln("binary.Write failed:", err))
			}
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

func loadShader(shaderType gl.Enum, assetName string) (gl.Shader, error) {
	// Borrowed from golang.org/x/mobile/exp/gl/glutil
	src, err := loadAsset(assetName)
	if err != nil {
		return gl.Shader{}, err
	}

	shader := gl.CreateShader(shaderType)
	if shader.Value == 0 {
		return gl.Shader{}, fmt.Errorf("failed to create shader (type %v)", shaderType)
	}
	gl.ShaderSource(shader, string(src))
	gl.CompileShader(shader)
	if gl.GetShaderi(shader, gl.COMPILE_STATUS) == 0 {
		defer gl.DeleteShader(shader)
		return gl.Shader{}, fmt.Errorf("shader compile: %s", gl.GetShaderInfoLog(shader))
	}
	return shader, nil
}

func LoadShaders(program gl.Program, vertexAsset, fragmentAsset string) error {
	vertexShader, err := loadShader(gl.VERTEX_SHADER, vertexAsset)
	if err != nil {
		return err
	}
	fragmentShader, err := loadShader(gl.FRAGMENT_SHADER, fragmentAsset)
	if err != nil {
		gl.DeleteShader(vertexShader)
		return err
	}

	if gl.GetProgrami(program, gl.ATTACHED_SHADERS) > 0 {
		for _, shader := range gl.GetAttachedShaders(program) {
			gl.DetachShader(program, shader)
		}
	}

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	// Flag shaders for deletion when program is unlinked.
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	if gl.GetProgrami(program, gl.LINK_STATUS) == 0 {
		defer gl.DeleteProgram(program)
		return fmt.Errorf("LoadShaders: %s", gl.GetProgramInfoLog(program))
	}
	return nil
}

// LoadProgram reads shader sources from the asset repository, compiles, and
// links them into a program.
func LoadProgram(vertexAsset, fragmentAsset string) (program gl.Program, err error) {
	log.Println("LoadProgram:", vertexAsset, fragmentAsset)

	program = gl.CreateProgram()
	if program.Value == 0 {
		return gl.Program{}, fmt.Errorf("glutil: no programs available")
	}

	err = LoadShaders(program, vertexAsset, fragmentAsset)
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

func AppendIndexed(slice []float32, idx *[]int, vertices ...float32) []float32 {
	// FIXME: This is the wrong algo.
	idxMap := map[[vertexDim]float32]int{}
	r := slice

	var vert [3]float32
	for _, pos := range *idx {
		vert[0] = slice[pos]
		vert[1] = slice[pos+1]
		vert[2] = slice[pos+2]
		idxMap[vert] = pos
	}

	for i := 0; i < len(vertices); i += vertexDim {
		vert[0] = vertices[i]
		vert[1] = vertices[i+1]
		vert[2] = vertices[i+2]
		pos, ok := idxMap[vert]
		if ok {
			*idx = append(*idx, pos)
			continue
		}
		pos = len(r) / vertexDim
		r = append(r, vert[:]...)
		idxMap[vert] = pos
		*idx = append(*idx, pos)
	}

	return r
}

// MultiMul multiplies every non-nil Mat4 reference and returns the result. If
// none are given, then it returns the identity matrix.
func MultiMul(matrices ...*mgl.Mat4) mgl.Mat4 {
	var r mgl.Mat4
	ok := false
	for _, m := range matrices {
		if m == nil {
			continue
		}
		if !ok {
			r = *m
			ok = true
			continue
		}
		r = r.Mul4(*m)
	}
	if ok {
		return r
	}
	return mgl.Ident4()
}

func Quad(a mgl.Vec3, b mgl.Vec3) []float32 {
	return []float32{
		// First triangle
		b[0], b[1], b[2], // Top Right
		a[0], b[1], a[2], // Top Left
		a[0], a[1], a[2], // Bottom Left
		// Second triangle
		a[0], a[1], a[2], // Bottom Left
		b[0], b[1], b[2], // Top Right
		b[0], a[1], b[2], // Bottom Right
	}
}

func Upvote(tip mgl.Vec3, size float32) []float32 {
	a := tip.Add(mgl.Vec3{-size / 2, -size * 2, 0})
	b := tip.Add(mgl.Vec3{size / 2, -size, 0})
	return []float32{
		tip[0], tip[1], tip[2], // Top
		tip[0] - size, tip[1] - size, tip[2], // Bottom left
		tip[0] + size, tip[1] - size, tip[2], // Bottom right

		// Arrow handle
		b[0], b[1], b[2], // Top Right
		a[0], b[1], a[2], // Top Left
		a[0], a[1], a[2], // Bottom Left
		a[0], a[1], a[2], // Bottom Left
		b[0], b[1], b[2], // Top Right
		b[0], a[1], b[2], // Bottom Right
	}
}

// IsBoundingBox returns true if a box intercepts b box.
func IsBoxCollision(a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y float32) bool {
	return a1_x <= b2_x && a2_x >= b1_x && a1_y <= b2_y && a2_y >= b1_y
}

// IsCollision returns true if segment a1->a2 intersects segment b1->b2.
func IsCollision2D(a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y float32) bool {
	/*
		if !IsBoxCollision(a1_x, a1_y, a2_x, a2_y, b1_x, b1_y, b2_x, b2_y) {
			// Short circuit check
			return false
		}
	*/

	// Based on http://stackoverflow.com/a/1968345/187878
	s1_x := a2_x - a1_x
	s1_y := a2_y - a1_y
	s2_x := b2_x - b1_x
	s2_y := b2_y - b1_y

	denom := s1_x*s2_y - s2_x*s1_y
	if denom == 0 {
		// Collinear
		return a1_x == b1_x || a1_x == b2_x
	}
	denomPositive := denom > 0

	s3_x := a1_x - b1_x
	s3_y := a1_y - b1_y

	s_numer := s1_x*s3_y - s1_y*s3_x
	if (s_numer < 0) == denomPositive {
		return false
	}

	t_numer := s2_x*s3_y - s2_y*s3_x
	if (t_numer < 0) == denomPositive {
		return false
	}

	if ((s_numer > denom) == denomPositive) || ((t_numer > denom) == denomPositive) {
		return false
	}

	// Do we want the intersecting point?
	// i_x = a1_x + (t * s1_x)
	// i_y = a1_x + (t * s1_y)
	return true
}
