//<editor-fold desc="Changeable Configuration Block">
window.onload = function() {
    window.ui = SwaggerUIBundle({
        url: "swagger.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true,
        displayRequestDuration: true,
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
        ],
        plugins: [
            SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        defaultModelsExpandDepth: 1,
        tryItOutEnabled: true
    });
};
//</editor-fold>