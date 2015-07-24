#version 100

precision mediump float;

uniform vec3 lightPosition;
uniform vec3 lightIntensities;
uniform samplerCube tex;
uniform vec3 cameraCoord;
uniform vec4 surfaceColor;

varying vec3 fragNormal;
varying vec3 fragCoord;

const vec3 diffuseColor = vec3(0.7, 0.3, 0.3);
const float shininess = 16.0;
const float screenGamma = 2.2; // Assume the monitor is calibrated to the sRGB color space
const float refractRatio = 1.0 / 1.52;

void main() {
    vec3 normal = normalize(fragNormal);
    vec3 lightDir = normalize(lightPosition - fragCoord);

    float lambertian = max(dot(lightDir,normal), 0.0);
    float specular = 0.0;

    if (lambertian > 0.0) {
        vec3 viewDir = normalize(-fragCoord);
        vec3 halfDir = normalize(lightDir + viewDir);
        float specAngle = max(dot(halfDir, normal), 0.0);
        specular = pow(specAngle, shininess);
    }

    vec3 I = normalize(fragCoord - cameraCoord);
    vec3 R = refract(I, normal, refractRatio);

    vec4 surface = surfaceColor;
    if (surface.a < 1.0) {
        surface = normalize(textureCube(tex, R) + surface);
    }
    vec3 colorLinear = surface.rgb + lambertian * diffuseColor + specular * lightIntensities;

    // apply gamma correction (assume ambientColor, diffuseColor and lightIntensities
    // have been linearized, i.e. have no gamma correction in them)
    vec3 colorGammaCorrected = pow(colorLinear, vec3(1.0/screenGamma));

    // use the gamma corrected color in the fragment
    //gl_FragColor = vec4(colorGammaCorrected, surface.a);
    gl_FragColor = vec4(colorGammaCorrected, surface.a);
}
