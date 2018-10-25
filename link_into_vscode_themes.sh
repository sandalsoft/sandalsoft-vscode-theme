#!/bin/sh

VSCODE_EXT_DIR='~/.vscode/extensions'

ln -s $(pwd) ~/.vscode/extensions/$(basename $(pwd))