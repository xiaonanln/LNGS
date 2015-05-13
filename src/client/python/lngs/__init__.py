# coding: utf8

import asyncore

from GameClient import GameClient
from Entity import Entity

def loop(timeout=0.001):
	asyncore.loop(timeout)

def connect(host, port):
	pass
