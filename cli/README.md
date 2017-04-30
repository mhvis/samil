# Samil CLI

This is a very simple command-line application for Samil Power inverters.
It searches the network for inverters.
For each inverter found, model info and current generation data is printed once.
When no new inverters are found after a minute, the application exits.

You can specify the interface to bind to on the command-line, use `./samil -h`
to get usage information.

A binary for 386, amd64 or arm (Raspberry Pi) can be found
[here](https://github.com/mhvis/samil/releases/latest).

If you would like to see more useful features in this dumb application,
e.g. making it scriptable or make it possible to request history data,
please let me know.
The best way to let me know is by making a
[new issue](https://github.com/mhvis/samil/issues/new).