# External auth plugins
This repository contains two public interfaces:
- `AuthService`: is the interface implemented by all [Gloo ext auth implementations](https://gloo.solo.io/gloo_routing/virtual_services/authentication/)
- `ExtAuthPlugin`: is the interface that needs to be implemented by custom go ext auth plugins

Be sure to check out the docs at gloo.solo.io for more information!