# Moving on

Since 2019 when I first looked into the Meross devices I have been able to use Home Assistant reliably with them.
After creating this Home Assistant Meross platform, the project at https://github.com/krahabb/meross_lan has gained
traction, I have switched to it and I have no interest in maintaining my own.

Using this work to find out how the devices function, being able to run them without having internet access has
scratched my curiosity itch and I will no longer updating this repository.

I will continue to use Meross devices and I may come back and improve things should something break, or I learn
something new that is worth sharing. That being said it will likely be fed back into the meross_lan project note here.

Links:

  - https://github.com/krahabb/meross_lan
  - https://www.meross.com

# Table of Contents

* [About the Project](#about-the-project)
  * [Built With](#built-with)
* [Installation](#installation)
* [Usage](#usage)
* [Contributing](#contributing)
* [License](#license)
* [Acknowledgements](#acknowledgements)

# About The Project

I purchased a Meross mss310 mainly because it was cheap. I decided that rather than connecting it to a server and
giving control of my switches to someone else that I wanted to understand how they worked to see if I could use Home
Assistant to control these switches completly disconnected from the Meross servers.

With the help of others it turns out this not only possible, but the implementation is relatively clean provided you
are ok with running an MQTT server.

## Built With

This project can be built and compiled with Go and its standard libraries

* [Go](https://golang.org)

# Installation

This repository consists of two command line tools. To install them both you can run the following.

```bash
go get https://github.com/arandall/meross/cmd/meross-cloud # Used to get key for existing devices only.
go get https://github.com/arandall/meross/cmd/meross-device
```

<!-- USAGE EXAMPLES -->
## Usage

See the [provisioning](doc/provisioning.md) page for details.

<!-- CONTRIBUTING -->
## Contributing

If you have another meross device or find somethnig that isn't quite right raise an issue and/or PR.

<!-- LICENSE -->
## License

Distributed under the MIT License.

<!-- ACKNOWLEDGEMENTS -->
## Acknowledgements

Thanks to the following project that got me off to a good start.

* https://github.com/bapirex/meross-api for providing details for `meross-cloud` to obtain existing device keys.
* https://github.com/mrgsts/mss310-kontrol for showing the JSON API details.
