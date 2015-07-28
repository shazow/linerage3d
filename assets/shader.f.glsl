#version 100

precision mediump float;

uniform vec3 lightPosition;
uniform vec3 lightIntensities;
uniform vec3 cameraCoord;
uniform vec4 surfaceColor;
uniform samplerCube tex;

varying vec3 fragNormal;
varying vec3 fragCoord;
varying vec3 fragPos;

const float shininess = 16.0;
const float screenGamma = 2.2; // Assume the monitor is calibrated to the sRGB color space
const float refractRatio = 1.0 / 1.52;
const vec3 specColor = vec3(1.0, 1.0, 1.0);


vec3 Light_BlinnPhong(vec4 surface, float shine, vec3 fragPos, vec3 lightPos, vec3 lightColor) {
    // Based on https://en.wikipedia.org/wiki/Blinn%E2%80%93Phong_shading_model

    vec3 lightDir = normalize(lightPos - fragPos);

    float lambertian = max(dot(lightDir, fragNormal), 0.0);
    float specular = 0.0;

    if (lambertian > 0.0) {
        vec3 viewDir = normalize(-fragPos);
        vec3 halfDir = normalize(lightDir + viewDir);
        float specAngle = max(dot(halfDir, fragNormal), 0.0);
        //specular = pow(specAngle, shine);
    }

    vec3 colorLinear = surface.rgb + lambertian * lightColor + specular * specColor;

    // apply gamma correction (assume ambientColor, diffuseColor and lightIntensities
    // have been linearized, i.e. have no gamma correction in them)
    vec3 colorGammaCorrected = pow(colorLinear, vec3(1.0/screenGamma));

    return colorGammaCorrected;
}

void main() {
    vec3 I = normalize(fragPos - cameraCoord);
    vec3 R = refract(I, fragNormal, refractRatio);

    //vec4 surface = surfaceColor;
    vec4 surface = vec4(0.1,0.15,0.4,1);
    if (surface.a < 1.0) {
        surface = mix(textureCube(tex, R), surface, surface.a);
    }

    vec3 lit = Light_BlinnPhong(surface, shininess, fragCoord, lightPosition, vec3(1,1,1));

    // use the gamma corrected color in the fragment
    gl_FragColor = vec4(lit, surface.a);
}
