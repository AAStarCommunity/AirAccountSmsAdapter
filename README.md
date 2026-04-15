# AirAccountSmsAdapter

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
We will need many adapters, the first is SMS adapter using SIM800/SIM900a chips as adapter.
To get SMS and parse into instructions, invoke the SDK of sim800 with instructions.

Basic instructions:
1. Check if there is any new SMS, if yes, parse it and send to the server.(use `AT+CMGL="ALL"` to get all SMS, sim800c will return a list of SMS, each SMS is a string, parse it and get the instruction)
2. Transfer into instruction and send to the Gateway.
3. Get response from the Gateway and send back to the sender.

OK !

# Auto-Update

**Only works for raspberrypi**

> See cmd/bash.sh

## License

Licensed under the [Apache License, Version 2.0](https://opensource.org/licenses/Apache-2.0). See [LICENSE](./LICENSE) for details.
