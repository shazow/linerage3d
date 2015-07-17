package main

import "golang.org/x/mobile/gl"

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
	shader.vertCoord = gl.GetAttribLocation(shader.program, "vertCoord")
	shader.vertNormal = gl.GetAttribLocation(shader.program, "vertNormal")

	shader.projection = gl.GetUniformLocation(shader.program, "projection")
	shader.view = gl.GetUniformLocation(shader.program, "view")
	shader.model = gl.GetUniformLocation(shader.program, "model")
	shader.normalMatrix = gl.GetUniformLocation(shader.program, "normalMatrix")

	shader.lightIntensities = gl.GetUniformLocation(shader.program, "lightIntensities")
	shader.lightPosition = gl.GetUniformLocation(shader.program, "lightPosition")
}
