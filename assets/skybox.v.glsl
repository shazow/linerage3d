#version 100

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;

attribute vec3 vertCoord;

void main(){
    vec4 vertPos4 = view * model * vec4(vertCoord, 1.0);
    gl_Position = projection * vertPos4;

    fragCoord = vertCoord;
}
