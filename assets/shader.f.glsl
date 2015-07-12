#version 100

precision mediump float;

uniform vec3 lightPosition;
uniform vec3 lightIntensities;

varying vec3 fragNormal;
varying vec3 fragCoord;

const vec4 surfaceColor = vec4(0.3, 0.3, 0.9, 1.0);

void main() {
    float lightAngle = dot(fragNormal, lightPosition);

    // use the gamma corrected color in the fragment
    gl_FragColor = surfaceColor * lightAngle;
}
