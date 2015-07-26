#version 100

uniform mat4 projection;
uniform mat4 view;

attribute vec3 vertCoord;

varying vec3 fragCoord;

void main(){
    vec4 pos = projection * view * vec4(vertCoord, 1.0);
    gl_Position = pos.xyww;
    fragCoord = vertCoord;
}
