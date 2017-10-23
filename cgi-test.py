from http.server import CGIHTTPRequestHandler, HTTPServer


class Handler(CGIHTTPRequestHandler):
    cgi_directories = ["/"]


port = 9999
httpd = HTTPServer(("", port), Handler)
print('CGI server at localhost:{}.'.format(port))
httpd.serve_forever()
