#version 100

uniform samplerCube skybox;

varying vec3 fragCoord;

void main() {
    gl_FragColor = textureCube(skybox, fragTexCoord);
}
