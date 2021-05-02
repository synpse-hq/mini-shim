# Mini Shim

Mini-shim allows modifying headers (and potentially other things) for your applications running in Docker containers.

## Usage

Normally mini-shim would be deployed next to the application, proxying all requests to it

```yaml
name: homeassistant
description: Home Assistant with Webhook Relay tunnel. (https://kr-homelab.webrelay.io)
scheduling:
  type: Conditional
  selectors:
    type: rpi
spec:
  containers:
    # Your main container
    - name: homeassistant
      image: docker.io/homeassistant/raspberrypi4-homeassistant:stable
      hostname: homeassistant
      ports:
        - 8123:8123
      volumes:
        - /usr/homeassistant:/config
        - /etc/localtime:/etc/localtime
    # Mini-shim container, proxying incoming requests
    # to the main container
    - name: shim
      image: quay.io/synpse/mini-shim:latest
      env:
        - name: ALLOWED_ORIGINS
          value: "*"                       # Set CORS origins
        - name: UPSTREAM_ADDR
          value: http://homeassistant:8123 # Specify where the requests should be routed 
        - name: LISTEN_ADDR
          value: :14000                    # Default port on which proxy listens for requests
    # Webhook Relay container to expose the application to the internet
    - name: relayd
      image: webhookrelay/webhookrelayd-aarch64:1
      args:
        - --mode
        - tunnel
        - -t
        - homelab
      env:
        - name: RELAY_KEY
          fromSecret: relayKey
        - name: RELAY_SECRET
          fromSecret: relaySecret
```