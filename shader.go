package main

import "golang.org/x/mobile/gl"

// TODO: Need a ShaderRegistry of somekind, ideally with support for default
// scene values vs per-shape values and attribute checking.
// TODO: Should each NewShader be a struct embedding a Program?

type Shader interface {
	Use()
	Close()
	Attrib(string) gl.Attrib
	Uniform(string) gl.Uniform
}

func NewShader(vertAsset, fragAsset string) (Shader, error) {
	program, err := LoadProgram(vertAsset, fragAsset)
	if err != nil {
		return nil, err
	}

	return &shader{
		program:  program,
		attribs:  map[string]gl.Attrib{},
		uniforms: map[string]gl.Uniform{},
	}, nil
}

type shader struct {
	program gl.Program

	attribs  map[string]gl.Attrib
	uniforms map[string]gl.Uniform
}

func (shader *shader) Attrib(name string) gl.Attrib {
	v, ok := shader.attribs[name]
	if !ok {
		v = gl.GetAttribLocation(shader.program, name)
		shader.attribs[name] = v
	}
	return v
}

func (shader *shader) Uniform(name string) gl.Uniform {
	v, ok := shader.uniforms[name]
	if !ok {
		v = gl.GetUniformLocation(shader.program, name)
		shader.uniforms[name] = v
	}
	return v
}

func (shader *shader) Use() {
	gl.UseProgram(shader.program)
}

func (shader *shader) Close() {
	gl.DeleteProgram(shader.program)
}
