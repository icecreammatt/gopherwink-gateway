#!/bin/sh

openssl req -new -nodes -x509 -out client.pem -keyout client.key -days 365
