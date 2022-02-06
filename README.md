# hup-on-notify

This binary watches `-fileToWatch` and sends a SIGHUP to the process id (PID) found in `-pidFile`

# build for linux and docker

`make`

## notes

The included `Dockerfile` is a POC on how to send a SIGHUP to [squid-cache](http://www.squid-cache.org) when it's configuration changes.

To demonstrate:
1. build for linux and docker (above).
2. run the new container: `docker run -ti -p 3128:3128 hup-on-notify-example`
3. in a separate terminal, run the golang binary in the squid container: `docker exec -ti $(docker ps | grep hup-on-notify-example | awk '{print $1}') /usr/local/bin/hup-on-notify`
4. in a separate terminal, edit the `squid.conf`: `docker exec -ti $(docker ps | grep hup-on-notify-example | awk '{print $1}') vi /etc/squid/squid.conf` and save changes.

You should see squid reload.

# stupid docker tricks
1. [bind mount](https://docs.docker.com/storage/bind-mounts/) the `squid.conf` to `/etc/squid/squid.conf`

# stupid k8s tricks
1. run the golang binary as an additional container in the same pod as the process in question.  mount the `-fileToWatch` as a volume from a `ConfigMap`.  watch your process in question receive a SIGHUP when the `ConfigMap` object changes.
