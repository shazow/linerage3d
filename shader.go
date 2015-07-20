package main

import "golang.org/x/mobile/gl"

func NewShader(vertAsset, fragAsset string) (*Shader, error) {
	program, err := LoadProgram("shader.v.glsl", "shader.f.glsl")
	if err != nil {
		return nil, err
	}

	return &Shader{
		program: program,
	}, nil
}

type Shader struct {
	program      gl.Program
	vertCoord    gl.Attrib
	vertNormal   gl.Attrib
	vertTexCoord gl.Attrib

	projection   gl.Uniform
	view         gl.Uniform
	model        gl.Uniform
	normalMatrix gl.Uniform

	lightPosition    gl.Uniform
	lightIntensities gl.Uniform
}

func (shader *Shader) Bind() {
	gl.UseProgram(shader.program)

	shader.vertCoord = gl.GetAttribLocation(shader.program, "vertCoord")
	shader.vertNormal = gl.GetAttribLocation(shader.program, "vertNormal")

	shader.projection = gl.GetUniformLocation(shader.program, "projection")
	shader.view = gl.GetUniformLocation(shader.program, "view")
	shader.model = gl.GetUniformLocation(shader.program, "model")
	shader.normalMatrix = gl.GetUniformLocation(shader.program, "normalMatrix")

	shader.lightIntensities = gl.GetUniformLocation(shader.program, "lightIntensities")
	shader.lightPosition = gl.GetUniformLocation(shader.program, "lightPosition")
}

func (shader *Shader) Close() {
	gl.DeleteProgram(shader.program)
}
