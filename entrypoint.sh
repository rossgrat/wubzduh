#!/bin/sh
chown -R appuser:appuser /var/log/wubzduh
exec su-exec appuser "$@"
