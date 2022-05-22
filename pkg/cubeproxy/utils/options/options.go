package options

const (
	NginxFile     = "/etc/cubernetes/cubeproxy/nginx/"
	CopySrc       = "/etc/nginx/"
	SiteEnabled   = "sites-enabled/"
	SiteAvailable = "sites-available/"
	Suffix        = ".conf"
	DefaultConf   = "nginx.conf"

	DefaultFile = "default"

	DockerNamePrefix = "C8-DNS-Nginx-"
	DockerName       = "Nginx"

	DefaultConfContent = `
user www-data;
worker_processes auto;
pid /run/nginx.pid;

events {
	worker_connections 768;
}

http {
	sendfile on;
	tcp_nopush on;
	tcp_nodelay on;
	keepalive_timeout 65;
	types_hash_max_size 2048;

	include /etc/nginx/mime.types;
	default_type application/octet-stream;

	ssl_protocols TLSv1 TLSv1.1 TLSv1.2 TLSv1.3; # Dropping SSLv3, ref: POODLE
	ssl_prefer_server_ciphers on;

	access_log /var/log/nginx/access.log;
	error_log /var/log/nginx/error.log;

	gzip on;

	include /etc/nginx/conf.d/*.conf;
	include /etc/nginx/sites-enabled/*;
}
`
)
