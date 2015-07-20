#version 100

precision mediump float;

uniform samplerCube tex;

varying vec3 fragCoord;

void main() {
    gl_FragColor = textureCube(tex, fragCoord);
}
