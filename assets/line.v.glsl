#version 100
precision highp float;

uniform mat4 projection;
uniform mat4 view;
uniform mat4 model;
uniform mat4 normalMatrix;

attribute vec3 vertCoord;
attribute vec3 vertNormal;

varying vec3 fragPos;
varying vec3 fragCoord;
varying vec3 fragNormal;


void main(){
    vec4 vertPos4 = view * model * vec4(vertCoord, 1.0);

    gl_Position = projection * vertPos4;

    fragPos = vec3(vertPos4) / vertPos4.w;
    fragCoord = vertCoord;
    fragNormal = normalize(vec3(model * normalMatrix * vec4(vertNormal, 0.0)));
}
