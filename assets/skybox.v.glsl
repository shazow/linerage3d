#version 100

uniform mat4 projection;
uniform mat4 view;

attribute vec3 vertCoord;

varying vec3 fragCoord;

void main(){
    gl_Position = projection * view * vec4(vertCoord, 1.0);
    fragCoord = vertCoord;
}
