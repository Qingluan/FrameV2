#!/bin/bash
sleep 1;
which apt 2>/dev/null
if [ $? -eq 0 ] ; then
	
	sudo apt-get install -y gcc libgtk-3-dev libappindicator3-dev
	sudo apt install -y zenity
fi


ConfigDir="$HOME/.config/ProxyAnyWhere"
mkdir -p $ConfigDir


cp ProxyAnyWhere $ConfigDir

cp ico.png $ConfigDir
cp ProxyAnyWhere.desktop $HOME/Desktop/
