# url_shortener
A URL shortener service that not only shortens URLs but can handle high volumes of requests concurrently



# URL Shortening Flow:
 - Client sends POST request with original URL
 - Handler parses request and validates
 - Service generates short code
 - Saves to PostgreSQL 
 - Caches in Redis
 - Returns short URL to client

URL Redirection Flow:
 - Client requests original URL using short code
 - Handler extracts short code 
 - Checks Redis cache first
 - Falls back to PostgreSQL if not in cache
 - Redirects to original URL
