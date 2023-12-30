# AirAccountSmsAdapter

We will need many adapters, the first is SMS adapter using SIM800/SIM900a chips as adapter.
To get SMS and parse into instructions, invoke the SDK of sim800 with instructions.

Basic instructions:
1. Check if there is any new SMS, if yes, parse it and send to the server.(use `AT+CMGL="ALL"` to get all SMS, sim800c will return a list of SMS, each SMS is a string, parse it and get the instruction)
2. Transfer into instruction and send to the Gateway.
3. Get response from the Gateway and send back to the sender.

OK !

# Auto-Update

See cmd/bash.sh
