//<editor-fold desc="Changeable Configuration Block">
window.onload = function() {
    window.ui = SwaggerUIBundle({
        url: "swagger.yaml",  // Изменили путь
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
        defaultModelsExpandDepth: 1
    });
};
//</editor-fold>