package routers

// New Flutter-based firmware (v15+) uses a single JSON-RPC endpoint.
// All API calls go to /cgi/service.cgi with {"method":"...", "params":{...}}.
const serviceCGI = "/cgi/service.cgi"
