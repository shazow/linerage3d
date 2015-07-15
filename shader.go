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

func (shader *Shader) Draw(shape *Shape) {
	gl.UseProgram(shader.program)

	gl.BindBuffer(gl.ARRAY_BUFFER, shape.VBO)

	stride := shape.Stride()

	gl.EnableVertexAttribArray(shader.vertCoord)
	gl.VertexAttribPointer(shader.vertCoord, vertexDim, gl.FLOAT, false, stride, 0)

	if len(shape.normals) > 0 {
		gl.EnableVertexAttribArray(shader.vertNormal)
		gl.VertexAttribPointer(shader.vertNormal, normalDim, gl.FLOAT, false, stride, vertexDim*vecSize)
	}

	gl.DrawArrays(gl.TRIANGLES, 0, shape.Len())

	gl.DisableVertexAttribArray(shader.vertCoord)
	if len(shape.normals) > 0 {
		gl.DisableVertexAttribArray(shader.vertNormal)
	}
}
