window.onload = function() {
  const ui = SwaggerUIBundle({
    url: "swagger.json", // adjust if your spec is elsewhere
    dom_id: '#swagger-ui',
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    layout: "StandaloneLayout",
    requestInterceptor: (req) => {
      // Fetch domain from localStorage and set as header
      const domain = localStorage.getItem('domain');
      if (domain) {
        req.headers['domain'] = domain;
      }
      return req;
    },
    responseInterceptor: (res) => {
      if (
        res.url.includes('/login') && // adjust to your login endpoint
        res.status === 200 &&
        res.body &&
        res.body.accessToken // adjust to your response structure
      ) {
        localStorage.setItem('accessToken', res.body.accessToken);
        if (window.ui && window.ui.preauthorizeApiKey) {
          window.ui.preauthorizeApiKey('ApiKeyAuth', 'Bearer ' + res.body.accessToken);
        }
        // Parse JWT to extract domain
        try {
          const payload = JSON.parse(atob(res.body.accessToken.split('.')[1]));
          if (payload.domains[0].UUID) {
            localStorage.setItem('domain', payload.domains[0].UUID);
          }
        } catch (e) {
          console.error('Failed to parse JWT for domain:', e);
        }
        
      }
      return res;
    },
  });
  window.ui = ui;
};
