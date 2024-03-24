# Mozeidon native-app

The [mozeidon native-app](https://github.com/egovelox/mozeidon-native-app), a very simple ipc server written in ``go``, will allow the [mozeidon firefox-addon](https://addons.mozilla.org/en-US/firefox/addon/mozeidon) to receive commands from and send responses to the mozeidon CLI (for installing the CLI, see [mozeidon](https://github.com/egovelox/mozeidon)).


## Installation

On MacOS or Linux, you can install it using ``homebrew`` :
```bash
brew tap egovelox/homebrew-mozeidon ;

brew install egovelox/mozeidon/mozeidon-native-app ;
```

Otherwise, you may download the binary from the [release page](https://github.com/egovelox/mozeidon-native-app/releases).

If no release matches your platform, you can build the binary yourself:
```bash
git clone https://github.com/egovelox/mozeidon-native-app.git ;

cd mozeidon-native-app && go build
```

As a firefox native-app, it has to be referenced into your Firefox configuration.

### Referencing the native-app into your Firefox configuration

On ``MacOS``, first locate the ``~/Library/Application Support/Mozilla/NativeMessagingHosts`` directory (or create it if missing).

Then create a ``mozeidon.json`` file, and copy into it the following ``json``.

Note: depending on your installation, you may need to replace the value in ``"path"`` with the absolute path of the mozeidon-native-app.

```json
{
  "name": "mozeidon",
  "description": "Native messaging add-on to interact with your browser",
  "path": "/opt/homebrew/bin/mozeidon-native-app",
  "type": "stdio",
  "allowed_extensions": [
    "mozeidon-addon@egovelox.com"
  ]
}
```

Now the Mozeidon firefox-addon will be able to interact with the Mozeidon native-app.

Note : 
For other OS than ``MacOS``, please check the [Mozilla documentation](https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/Native_manifests#manifest_location) to find the correct location of the Firefox ``NativeMessagingHosts`` directory.

At last, you should be able to use the Mozeidon CLI or the Raycast extension.


## Releases

Various releases of the Mozeidon native-app can be found on the [releases page](https://github.com/egovelox/mozeidon-native-app/releases).

Releases are managed with github-actions and [goreleaser](https://github.com/goreleaser/goreleaser).

A release will be auto-published when a new git tag is pushed,
e.g :

```bash
git clone https://github.com/egovelox/mozeidon-native-app.git && cd mozeidon-native-app;

git tag -a v2.0.0 -m "A new mozeidon-native-app release"

git push origin v2.0.0
```

