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

const float screenGamma = 2.2; // Assume the monitor is calibrated to the sRGB color space

const int maxLights = 3;
struct Light
{
    vec3 color;
    vec3 position;
    float intensity;
};
uniform Light lights[maxLights];

struct Material
{
    vec3 ambient;
    vec3 diffuse;
    vec3 specular;
    float shininess;
    float refraction;
};
uniform Material material;


vec3 Light_BlinnPhong(Light light, vec3 fragPos)
{
    // Based on https://en.wikipedia.org/wiki/Blinn%E2%80%93Phong_shading_model

    vec3 lightDir = normalize(light.position - fragPos);

    float lambertian = max(dot(lightDir, fragNormal), 0.0);
    float specular = 0.0;

    if (lambertian > 0.0) {
        vec3 viewDir = normalize(-fragPos);
        vec3 halfDir = normalize(lightDir + viewDir);
        float specAngle = max(dot(halfDir, fragNormal), 0.0);
        specular = pow(specAngle, material.shininess);
    }

    // Linearized, before gamma correction
    return lambertian * light.color + specular * material.specular;
}

void main()
{
    vec3 fragColor = material.ambient;

    // Reflect
    if (material.refraction > 0.0) {
        vec3 I = normalize(fragPos - cameraCoord);
        vec3 R = refract(I, fragNormal, material.refraction);
        fragColor = vec3(mix(textureCube(tex, R), vec4(fragColor, 1.0), 0.5));
    }

    for (int i = 0; i < maxLights; i++)
    {
        fragColor += Light_BlinnPhong(lights[i], fragCoord);
    }

    // Gamma correct
    fragColor = pow(fragColor, vec3(1.0/screenGamma));

    gl_FragColor = vec4(fragColor, 1.0);
}
