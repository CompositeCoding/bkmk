## **bkmk: A Terminal Bookmark Manager**

## **Overview**

**bkmk is a command-line bookmark manager that allows you to store and open bookmarks from your terminal. It uses a simple and intuitive interface to manage your bookmarks, and supports features like tagging, searching, and importing from browser exported HTML files.**

## **Features**

* **Store and open bookmarks from your terminal**
* **Support for tagging and searching bookmarks**
* **Import bookmarks from browser exported HTML files**
* **Support for multiple profiles**

## **Usage**

## Commands

* `bkmk add <url>`: Add a new bookmark
* `bkmk open <url>`: Open a bookmark
* `bkmk delete <url>`: Delete a bookmark
* `bkmk import <path>`: Import bookmarks from a browser exported HTML file
* `bkmk profile <name>`: Create or switch to a bookmark profile

## Options

* `-a, --alias`: Set an alias for a bookmark
* `-c, --create`: Create a new profile
* `-s, --switch`: Switch to a profile

## **Installation**

**To install bkmk, simply clone this repository and run **`go build` to compile the binary. Place the binary on your OS PATH under bkmk.

## **Configuration**

**bkmk stores its configuration in a file named **`config` in the current working directory. The configuration file contains the profile name, local key, and server key. If the configuration file does not exist, bkmk will create it with default values.

## **License**

**See licence.md**
