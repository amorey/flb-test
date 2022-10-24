# Fluent Bit Test

This repo contains code to build and run a Fluent Bit go plugin that prints log records to stdout then sleeps for 5 seconds. This allows you to test the behavior of Fluent Bit when go plugin Flush() commands run slowly.

To build the plugin and deploy an instance of Fluent Bit with the plugin enabled run the following command:

```sh
docker compose up -d
```

To send a message to Fluent Bit you can exec into the running `flb` container and use the `nc` utility to send a mssage using TCP:

```sh
nc localhost 5170
msg1
msg2
```

You should see the following output in the `flb` container logs:

```
[plugin] Flush() start
[plugin] flushing: msg1
[plugin] flushing: msg2
[plugin] Flush() sleep start
[plugin] Flush() sleep stop
[plugin] Flush() ended
```
